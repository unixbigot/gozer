package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var activeOnly bool
var allNetworks bool

func init() {
	flag.BoolVar(&activeOnly, "active", false, "Show only active network members")
	flag.BoolVar(&allNetworks, "all", false, "Show status for all networks")
	flag.BoolVar(&allNetworks, "a", false, "Show status for all networks")
}

func main() {
	flag.Parse()
	if len(apiToken) == 0 {
		var token_buf []byte
		token_buf, err := ioutil.ReadFile(os.ExpandEnv(apiTokenFile))
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatal("Cannot read", err)
			}
			log.Println("No token file " + apiTokenFile)
		} else {
			apiToken = string(token_buf)
		}
	}
	if len(apiToken) == 0 {
		log.Fatal("No API token supplied")
	}

	client := ZeroTierClient{}

	// var payload interface{}
	//err = client.getJSON(apiUrl+apiCmdStatus, &payload)
	//if err != nil {
	//	logger.Fatal("Status request failed", err)
	//}
	if flag.NArg() == 0 && !allNetworks {
		_, err := client.ListNetworks(true)
		if err != nil {
			logger.Fatal(err)
		}
		os.Exit(0)
	}

	network_names := make([]string, 0, 100)

	if allNetworks {
		networks, err := client.ListNetworks(false)
		if err != nil {
			logger.Fatal(err)
		}
		for name := range networks.id_index {
			network_names = append(network_names, name)
		}
	} else {
		for _, name := range flag.Args() {
			network_names = append(network_names, name)
		}
	}
	logger.Verboseln(fmt.Sprintf("Showing detail for: %v", network_names))

	for _, arg := range network_names {

		// Get the information for a named network
		network, err := client.GetNetworkDetails(arg)
		if err != nil {
			fmt.Println(arg, " NOT FOUND ", err)
			continue
		}
		fmt.Println(network.SummaryString())

		// Get a list of members for the network, and iterate
		members := client.GetNetworkMemberDetails(network, activeOnly)
		for _, member := range members {
			fmt.Println("    ", member.SummaryString())
		}
		fmt.Println()
	}

}
