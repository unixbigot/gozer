# gozer - Golang + Zerotier
Gozer is an API client for the Zerotier One software-defined-network written in Go

Zerotier is a great (freemium) service that allows your Mac, Linux, Windows, iOS and
Android systems to form a "virtual LAN" no matter where they are in the world.

With Zerotier you can manage your servers even if they are behind a distant
firewall, or mobile, or dynamic IP.distant

I wrote this client becuase I wanted a way (besides using ZeroTier's web console)
to find the status and IP of all the hosts on my zerotier SDN.

By default gozer shows your connected network(s)

```bash
chris@Auric: () You are in ~/src/go/me/gozer (master)
So what are you going to DO about it? $ gozer
christopher_biggs    e5cdREDACTEDf921 192.168.191.1--192.168.191.254 devices belonging to Christopher Biggs
```

With flags, it can show all the hosts:

```bash
chris@Auric: () You are in ~/src/go/me/gozer (master)
So what are you going to DO about it? $ gozer -h
Usage of gozer:
  -a	Show status for all networks
  -active
    	Show only active network members
  -all
    	Show status for all networks
  -api-token string
    	ZeroTier API token
  -api-token-file string
    	File containing ZeroTier API token (default "${HOME}/.gozer-token")
  -debug
    	Turn on debug trace
  -verbose
    	Turn on verbose log

chris@Auric: () You are in ~/src/go/me/gozer (master)
So what are you going to DO about it? $ gozer -a
christopher_biggs    e5cdREDACTEDf921 192.168.191.1--192.168.191.254 devices belonging to Christopher Biggs
                     45SEKRITe3 []  Unauthorized Offline
     auric           0fSEKRITa5 [192.168.191.231               ]
     enceladus       38SEKRITbb [192.168.191.83                ]  Offline
     gabriel         6aSEKRITc4 [192.168.191.232               ]
     palantir        f1SEKRIT10 [192.168.191.54                ]  Offline
     pimento         03SEKRIT15 [192.168.191.180               ]
     pinacolada      47SEKRIT9f [192.168.191.76                ]  Offline
     pintsize        b5SEKRIT96 [192.168.191.105               ]  Offline
     serendipity     dfSEKRIT93 [192.168.191.164               ]
     tindog          97SEKRIT8e [192.168.191.4                 ]
     trippy          b8SEKRIT70 [192.168.191.170               ]  Offline
     tweety          53SEKRITd5 [192.168.191.206               ]
     umbriel         83SEKRITcb [192.168.191.124               ]
     vulpine         bfSEKRIT5a []
```

# FAQ

* **What if I do not own the network**

  I think you need to be able to authenticate as the network owner to 
  use the relevant parts of the API, sorry.

* **What if I am joined to several networks with different owners?**

  You need to authenticate as the network owner for each.

  You will need to install an API key for each owner, and use the command
  line options to select which owner to authenticate as.
  
  I welcome a pull request to query multiple accounts at once.

* **Gozer?**

   Gozer the Traveler. He will come in one of the pre-chosen forms. During the rectification of the Vuldrini, the traveler came  as a large and moving Torg! Then, during the third reconciliation of the last of the McKetrick supplicants, they chose a new form for him: that of a giant Slor! Many Shuvs and Zuuls knew what it was to be roasted in the depths of the Slor that day, I can tell you!
   
   
