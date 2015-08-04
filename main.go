package main

import (
	"fmt"
	"github.com/JReyLBC/JChecker/config"
)

func main() {
	config.Execute()

	cfg := config.GetConfig()
	fmt.Printf("cfg.ChassisEnvConfigFile: %s\n", cfg.ChassisEnvConfigFile)
	fmt.Printf("cfg.ChassisEnvIntervals: %v\n", cfg.ChassisEnvIntervals)
	fmt.Printf("cfg.ChassisEnvIPs: %v\n\n", cfg.ChassisEnvIPs)

	fmt.Printf("cfg.ChassisZonesConfigFile: %s\n", cfg.ChassisZonesConfigFile)
	fmt.Printf("cfg.ChassisZonesIntervals: %v\n", cfg.ChassisZonesIntervals)
	fmt.Printf("cfg.ChassisZonesIPs: %v\n\n", cfg.ChassisZonesIPs)

	fmt.Printf("cfg.NetconfUsername: %s\n", cfg.NetconfUsername)
	fmt.Printf("cfg.NetconfPassword: %s\n", cfg.NetconfPassowrd)
}
