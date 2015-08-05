package main

import (
	"github.com/JReyLBC/JChecker/config"
	"github.com/JReyLBC/JChecker/request"
	log "github.com/Sirupsen/logrus"
	"github.com/Juniper/go-netconf/netconf"
	"os"
)

func init() {
	config.Execute()

}

func main() {
	cfg := config.GetConfig()
	requests := make([]request.Request, 0, 32)

	if cfg.ChassisZonesIPs != nil {
		for _, ip := range cfg.ChassisZonesIPs {
			for _, dur := range cfg.ChassisEnvIntervals {
				requests = append(requests, request.NewChassisZonesRequest(
					ip, dur, cfg.NetconfUsername, cfg.NetconfPassowrd))
			}
		}

		var file *os.File
		var err error
		if file, err = os.Open(cfg.ChassisZonesResultsFile); err != nil {
			log.Errorln("Could not open results file.")
			log.Fatalln(err)
		} else {
			defer file.Close()
		}

		replyChan := make(chan *netconf.RPCReply)
		for _, request := range requests {
			request.Run(24, replyChan)
		}

/*		for ncReply := range replyChan {
			resp := response.N

		}*/
	}
/*
	if cfg.ChassisEnvIPs != nil {
		for _, ip := range cfg.ChassisEnvIPs {
			for _, dur := range cfg.ChassisEnvIPs {
			}
		}
	}*/
}
