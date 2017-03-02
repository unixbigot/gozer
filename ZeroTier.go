package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sort"
)

const apiCmdStatus = "/status"
const apiCmdNetwork = "/network"
const apiCmdNetworkDetailFormat = "/network/%s"
const apiCmdNetworkMembersFormat = "/network/%s/member"
const apiCmdNetworkMemberDetailFormat = "/network/%s/member/%s"

const apiUrl = "https://my.zerotier.com/api"

var apiToken string
var apiTokenFile string

// ZeroTierClient extends http client with methods for ZeroTier API
type ZeroTierClient struct{ http.Client }

/*
ListNetworks fetches a list of all configured networks
owned by the authenticated user, returning it as a pointer to a container struct


Not all the detail returned by the API is decoded, only what's been needed so far.
*/
func (client *ZeroTierClient) ListNetworks(print bool) (*ZeroTierNetworkList, error) {
	var networks = new(ZeroTierNetworkList)
	logger.Debugln("ListNetworks")

	// Make an API call to discover all networks that we own
	err := client.getJSON(apiCmdNetwork, &networks.networks)
	if err != nil {
		return networks, logger.Error(err, "Network request failed")
	}
	logger.Debugln(fmt.Sprintf("Parsed network response: %+v", networks.networks))

	// Augment the flat list of networks with two maps on Name and ID
	// Print the list (if requested) as a side-effect
	networks.name_index = make(map[string]*ZeroTierNetwork)
	networks.id_index = make(map[string]*ZeroTierNetwork)
	for _, network := range networks.networks {
		networks.name_index[network.Config.Name] = &network
		networks.id_index[network.ID] = &network
		if print {
			fmt.Println(network.SummaryString())
		}
	}

	// Return the slice of networks
	return networks, nil
}

// GetNetworkMembers fetches the members of a network as a map of ID to ???
// The documentation doesn't say what the integer field in the map means.
func (client *ZeroTierClient) GetNetworkMembers(networkID string) (members map[string]int, err error) {
	members = make(map[string]int)
	logger.Debugln("GetNetworkMembers", networkID)

	err = client.getJSON(
		fmt.Sprintf(apiCmdNetworkMembersFormat, networkID),
		&members)
	if err != nil {
		return
	}
	logger.Debugln(fmt.Sprintf("Parsed network response: %+v", members))
	return
}

// GetNetworkDetails fetches the information for a single network
// in the same format as returned by ListNetworks
func (client *ZeroTierClient) GetNetworkDetails(networkID string) (network *ZeroTierNetwork, err error) {
	network = new(ZeroTierNetwork)

	logger.Debugln("GetNetworkDetails", networkID)

	err = client.getJSON(
		fmt.Sprintf(apiCmdNetworkDetailFormat, networkID),
		network)
	if err != nil {
		return network, err
	}
	logger.Debugln(fmt.Sprintf("Parsed network response: %+v", network))
	return network, nil
}

func (client *ZeroTierClient) GetNetworkMemberDetails(network *ZeroTierNetwork, active bool) []ZeroTierNetworkMember {
	logger.Debugln("GetNetworkMemberDetails", network.ID)
	members := make([]ZeroTierNetworkMember, 0)

	var memberIDs []string
	if active {
		for id := range network.ActiveMembers {
			memberIDs = append(memberIDs, id)
		}
	} else {
		allMembers, err := client.GetNetworkMembers(network.ID)
		if err != nil {
			logger.Error(err, "Can't list members")
			return members
		}
		for id := range allMembers {
			memberIDs = append(memberIDs, id)
		}
	}

	for _, id := range memberIDs {

		member, err := client.GetMemberDetail(network.ID, id)
		if err != nil {
			logger.Error(err)
			continue
		}
		_, member.online = network.ActiveMembers[id]
		pretty, err := json.MarshalIndent(member, "", "  ")
		if err != nil {
			logger.Error(err, "JSON pretty print failed")
			continue
		}
		logger.Debugln(id, string(pretty))

		members = append(members, *member)
	}

	sort.Sort(byName(members))

	return members
}

// GetMemberDetail returns the information for one member of a network
// including its name and IP (which annoyingly is not included
// in the network details)
func (client *ZeroTierClient) GetMemberDetail(networkId string, clientId string) (*ZeroTierNetworkMember, error) {
	logger.Debugln("GetMemberDetail")

	var detail ZeroTierNetworkMember

	// Make an API call to discover all networks that we own
	err := client.getJSON(fmt.Sprintf(apiCmdNetworkMemberDetailFormat, networkId, clientId), &detail)
	if err != nil {
		return &detail, logger.Error(err, "Network request failed")
	}
	logger.Debugln(fmt.Sprintf("Parsed detail response: %+v", detail))

	return &detail, nil
}

// getJSON makes an API call to the ZeroTier rest API
// and decodes the result into the supplied Go structure
func (client *ZeroTierClient) getJSON(cmd string, payload interface{}) error {

	logger.Debugln("GET ", cmd)
	// Build a request
	req, err := http.NewRequest("GET", apiUrl+cmd, nil)
	if err != nil {
		return logger.Error(err, "Can't create request object")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiToken)

	// Send the request
	dump, err := httputil.DumpRequest(req, true)
	logger.Debugln("Sending request\n" + string(dump))
	resp, err := client.Do(req)
	if err != nil {
		return logger.Error(err, "HTTP request failed")
	}

	// Process the response
	defer resp.Body.Close()
	dump, err = httputil.DumpResponse(resp, true)
	logger.Debugln("Received response\n" + string(dump))

	dec := json.NewDecoder(resp.Body)
	for {
		if err := dec.Decode(payload); err == io.EOF {
			break
		} else if err != nil {
			logger.Error(err, "Cannot parse response", string(dump))
		}
	}

	// Log the prettified the response
	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return logger.Error(err, "JSON pretty print failed")
	}
	logger.Debugln("Decoded response", string(pretty))

	// Convert server error to go error
	if resp.StatusCode != 200 {
		return logger.Error(errors.New(resp.Status), "Request returned error ")
	}

	return nil
}

// ZeroTierNetwork represents the most useful information about a single network
type ZeroTierNetwork struct {
	ID            string
	Description   string
	Config        ZeroTierNetworkConfig
	ActiveMembers map[string]interface{}
}

// ZeroTierNetworkConfig holds the configuration information for a network
type ZeroTierNetworkConfig struct {
	Name              string
	Private           bool
	IPAssignmentPools []struct {
		IPRangeStart string
		IPRangeEnd   string
	}
}

// ZeroTierNetworkList encapsulates a list of networks returned by /network
type ZeroTierNetworkList struct {
	networks   []ZeroTierNetwork
	name_index map[string]*ZeroTierNetwork
	id_index   map[string]*ZeroTierNetwork
}

func (network *ZeroTierNetwork) SummaryString() string {
	return fmt.Sprintf("%-20s %s %s--%s %s",
		network.Config.Name,
		network.ID,
		network.Config.IPAssignmentPools[0].IPRangeStart,
		network.Config.IPAssignmentPools[0].IPRangeEnd,
		network.Description)
}

func (networks *ZeroTierNetworkList) FindIDorName(s string) (network *ZeroTierNetwork) {
	var ok bool
	if network, ok = networks.id_index[s]; ok {
		return
	}
	if network, ok = networks.id_index[s]; ok {
		return
	}
	return nil
}

type ZeroTierNetworkMember struct {
	NetworkID   string
	NodeID      string
	Name        string
	Description string
	Config      struct {
		Authorized    bool
		ActiveBridge  bool
		Hidden        bool
		IPAssignments []string
	}
	online bool
}

type byName []ZeroTierNetworkMember

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (member *ZeroTierNetworkMember) isOnline() bool {
	return member.online
}

func (member *ZeroTierNetworkMember) SummaryString() (summary string) {
	summary = fmt.Sprintf("%-15s %s %-30v %s",
		member.Name,
		member.NodeID,
		member.Config.IPAssignments,
		member.Description)
	if !member.Config.Authorized {
		summary += " Unauthorized"
	}
	if member.Config.ActiveBridge {
		summary += " Bridged"
	}
	if member.Config.Hidden {
		summary += " Hidden"
	}
	if !member.isOnline() {
		summary += " Offline"
	}
	return

}

func init() {
	flag.StringVar(&apiToken, "api-token", "", "ZeroTier API token")
	flag.StringVar(&apiTokenFile, "api-token-file", "${HOME}/.gozer-token", "File containing ZeroTier API token")
}
