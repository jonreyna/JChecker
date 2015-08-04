package main

import (
	"github.com/JReyLBC/JChecker/config"
	"fmt"
)

func main() {
	config.Execute()

	cfg := config.GetConfig()

	fmt.Printf("cfg.ShowChassisZones = %t\n", cfg.ShowChassisZones)
	fmt.Printf("cfg.ShowChassisFan = %t\n", cfg.ShowChassisFan)
	fmt.Printf("cfg.ShowChassisEnvironment = %t\n", cfg.ShowChassisEnvironment)
	fmt.Printf("cfg.TempCheckInterval = %s\n", cfg.TempCheckInterval.String())
}
