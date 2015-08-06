package main

import (
	"os"

	"encoding/csv"
	"github.com/JReyLBC/JChecker/config"
	"github.com/JReyLBC/JChecker/request"
	"github.com/JReyLBC/JChecker/response"
	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

func init() {
	config.Execute()

}

func main() {
	cfg := config.GetConfig()
	requests := make([]request.Request, 0, 32)

	log.Infof("chassisEnvIntervals: %v", cfg.ChassisEnvIntervals)
	if cfg.ChassisEnvIPs != nil {
		for _, ip := range cfg.ChassisEnvIPs {
			for _, dur := range cfg.ChassisEnvIntervals {
				requests = append(requests, request.NewChassisEnvRequest(
					ip, dur, cfg.NetconfUsername, cfg.NetconfPassowrd))
			}
		}

		var file *os.File
		var err error
		if file, err = os.OpenFile(cfg.ChassisEnvResultsFile, os.O_RDWR, os.ModeAppend); err != nil {
			log.Errorln("Could not open results file.")
			log.Fatalln(err)
		} else {
			defer file.Close()
		}
		csvWriter := csv.NewWriter(file)

		replyChan := make(chan *netconf.RPCReply)
		for _, request := range requests {
			request.Run(96, replyChan)
		}

		for ncReply := range replyChan {
			if resp, err := response.NewChassisEnvResponse(ncReply); err != nil {
				log.Errorln("Error creating response struct")
				log.Errorln(err)
			} else {
				log.Infoln("Writing to CSV file")
				resp.WriteCSV(csvWriter)
			}

		}
	} else if cfg.ChassisZonesIPs != nil {
		log.Errorln("Zones not implemented yet")
	}
	/*
		if cfg.ChassisEnvIPs != nil {
			for _, ip := range cfg.ChassisEnvIPs {
				for _, dur := range cfg.ChassisEnvIPs {
				}
			}
		}*/
}
