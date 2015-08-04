// Package config maintains the global configuration
// required to run LookingGlass. Configurations can be
// set on the command line or in a config file.
package config

import (
	"os"
	"time"

	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Contains all config options for LookingGlass
type JCheckerConfig struct {
	// show chassis environment stuff
	ChassisEnvIntervals    []time.Duration
	ChassisEnvIPs          []net.IP
	ChassisEnvConfigFile   string

	// show chassis zones stuff
	ChassisZonesIntervals  []time.Duration
	ChassisZonesIPs        []net.IP
	ChassisZonesConfigFile string

	// NETCONF stuff
	NetconfUsername        string
	NetconfPassowrd        string
}

// Default options
const (
	DEFAULT_LOG_LEVEL = log.InfoLevel

// NETCONF stuff
	DEFAULT_NETCONF_USERNAME = "admin"
	DEFAULT_NETCONF_PASSWORD = "abc123"
	DEFAULT_NETCONF_USERNAME_LOPT = "netconf-user"
	DEFAULT_NETCONF_USERNAME_OPT = "u"
	DEFAULT_NETCONF_PASSWORD_LOPT = "netconf-password"
	DEFAULT_NETCONF_PASSWORD_OPT = "p"

// show chassis environment stuff
	DEFAULT_CHASSIS_ENV_CONFIG_FILE = "chassis_env.csv"
	DEFAULT_CHASSIS_ENV_CONFIG_LOPT = "chassis-env-config"
	DEFAULT_CHASSIS_ENV_CONFIG_OPT = "c"

// show chassis zones stuff
	DEFAULT_CHASSIS_ZONES_CONFIG_FILE = "chassis_zones.csv"
	DEFAULT_CHASSIS_ZONES_CONFIG_LOPT = "chassis-zones-config"
	DEFAULT_CHASSIS_ZONES_CONFIG_OPT = "c"
)

// Package scope variables to handle
// config file access and overall config.
var (
	jCheckerConfig *JCheckerConfig
	jChecker       *cobra.Command
)

func init() {

	cntxLog := log.WithFields(log.Fields{
		"func": "config.init()",
	})

	// Check JCHECKER_CONFIG_LOGLEVEL and set logger accordingly
	switch os.Getenv("JCHECKER_CONFIG_LOGLEVEL") {
	case "DEBUG", "Debug", "debug", "5":
		log.SetLevel(log.DebugLevel)
	case "INFO", "Info", "info", "4":
		log.SetLevel(log.InfoLevel)
	case "WARN", "Warn", "warn", "3":
		log.SetLevel(log.WarnLevel)
	case "ERROR", "Error", "error", "2":
		log.SetLevel(log.ErrorLevel)
	case "FATAL", "Fatal", "fatal", "1":
		log.SetLevel(log.FatalLevel)
	case "PANIC", "Panic", "panic", "0":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	jCheckerConfig = new(JCheckerConfig)

	jChecker = &cobra.Command{
		Use:   "JChecker",
		Short: "JChecker executes requested checks on Juniper devices at user defined time intervals",
		Long: "JChecker executes requested checks at user defined time intervals.\n\n" +
		"It currently supports:\n" +
		"    - show chassis environment\n" +
		"    - show chassis zones\n",
	}

	jChecker.LocalFlags().StringVarP(
		&jCheckerConfig.NetconfUsername,
		DEFAULT_NETCONF_USERNAME_LOPT, DEFAULT_NETCONF_USERNAME_OPT,
		DEFAULT_NETCONF_USERNAME, "The username used to connect via NETCONF.",
	)

	jChecker.LocalFlags().StringVarP(
		&jCheckerConfig.NetconfPassowrd,
		DEFAULT_NETCONF_PASSWORD_LOPT, DEFAULT_NETCONF_PASSWORD_OPT,
		DEFAULT_NETCONF_PASSWORD, "The password used to connect via NETCONF.",
	)

	helpCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(0)
		},
	}

	envRequestCmd := &cobra.Command{
		Use:   "environment",
		Short: "Get the chassis environment information.\n",
		Long:  "Gets the chassis environment from the given list of IP addresses.",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Lookup(DEFAULT_CHASSIS_ENV_CONFIG_LOPT).Changed {
				cntxLog.Info("Default chassis environment config file changed")
			}
		},
	}

	envRequestCmd.Flags().StringVarP(
		&jCheckerConfig.ChassisEnvConfigFile,
		DEFAULT_CHASSIS_ENV_CONFIG_LOPT, DEFAULT_CHASSIS_ENV_CONFIG_OPT,
		DEFAULT_CHASSIS_ENV_CONFIG_FILE, "A csv file containing IP addresses with time intervals.",
	)

	zonesCmd := &cobra.Command{
		Use:   "zones",
		Short: "Get the chassis zones information.\n",
		Long:  "Gets the chassis zones information from the given list of IP addresses and timeouts.",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Lookup(DEFAULT_CHASSIS_ZONES_CONFIG_LOPT).Changed {
				cntxLog.Info("Default chassis zones config file changed")
			}
		},
	}

	zonesCmd.Flags().StringVarP(
		&jCheckerConfig.ChassisZonesConfigFile,
		DEFAULT_CHASSIS_ZONES_CONFIG_LOPT, DEFAULT_CHASSIS_ZONES_CONFIG_OPT,
		DEFAULT_CHASSIS_ZONES_CONFIG_FILE, "A csv file containing IP addresses with time intervals.",
	)

	jChecker.AddCommand(helpCmd)
	jChecker.AddCommand(envRequestCmd)
	jChecker.AddCommand(zonesCmd)

	cntxLog.Debugln("Setting persistent command flags")

	// This message will be shown to Windows users if Looking Glass is opened from explorer.exe
	cobra.MousetrapHelpText = `

	JChecker is a command line application.

	You need to open cmd.exe and run it from there.
	`
}

// Get the config options
func Execute() {

	cntxLog := log.WithFields(log.Fields{
		"func": "config.Execute()",
	})

	cntxLog.Debugln("Executing CLI config parser")

	jChecker.Execute()

}

func GetConfig() JCheckerConfig {
	return *jCheckerConfig
}
