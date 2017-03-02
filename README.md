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
ris@Auric: () You are in ~/src/go/me/gozer (master)
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
