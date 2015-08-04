// Package config maintains the global configuration
// required to run LookingGlass. Configurations can be
// set on the command line or in a config file.
package config

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Contains all config options for LookingGlass
type JCheckerConfig struct {
	TempCheckInterval      time.Duration
	ShowChassisFan         bool
	ShowChassisEnvironment bool
	ShowChassisZones       bool
}

// Long/short options
const (
	CHECK_INTERVAL_LOPT = "check-interval"
	CHECK_INTERVAL_OPT  = "I"

	SHOW_CHASSIS_FAN_LOPT = "show-chassis-fan"
	SHOW_CHASSIS_FAN_OPT  = "f"

	SHOW_CHASSIS_ENV_LOPT = "show-chassis-environment"
	SHOW_CHASSIS_ENV_OPT  = "e"

	SHOW_CHASSIS_ZONES_LOPT = "show-chassis-zones"
	SHOW_CHASSIS_ZONES_OPT  = "z"

	CONFIG_FILE_LOPT = "config"
	CONFIG_FILE_OPT  = "c"

	LOG_LEVEL_LOPT = "log-level"
	LOG_LEVEL_OPT  = "l"
)

// Default options
const (
	DEFAULT_NETCONF_USERNAME   = "admin"
	DEFAULT_NETCONF_PASSWORD   = "abc123"
	DEFAULT_LOG_LEVEL          = log.InfoLevel
	DEFAULT_CHECK_INTERVAL     = time.Duration(15 * time.Minute)
	DEFAULT_SHOW_CHASSIS_ENV   = true
	DEFAULT_SHOW_CHASSIS_FAN   = false
	DEFAULT_SHOW_CHASSIS_ZONES = false
	DEFAULT_CONFIG_FILE        = "jchecker_ip_list.csv"
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

	jChecker = &cobra.Command{
		Use:   "JChecker",
		Short: "JChecker executes requested checks on Juniper devices at user defined time intervals",
		Long: "JChecker executes requested checks at user defined time intervals.\n" +
			"It currently supports:\n" +
			"\t- show chassis environment\n" +
			"\t- show chassis zones\n",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			for _, str := range args {
				if str == "help" {
					cmd.Help()
					os.Exit(0)
				}
			}
		},
	}

	jCheckerConfig = new(JCheckerConfig)

	cntxLog.Debugln("Setting persistent command flags")

	// Set available flags, their defaults,
	// and their short options
	jChecker.PersistentFlags().DurationVarP(
		&jCheckerConfig.TempCheckInterval,
		CHECK_INTERVAL_LOPT, CHECK_INTERVAL_OPT,
		DEFAULT_CHECK_INTERVAL, `Default time interval for checks. Valid time units are ms, s, m, h`)

	jChecker.PersistentFlags().BoolVarP(
		&jCheckerConfig.ShowChassisEnvironment,
		SHOW_CHASSIS_ENV_LOPT, SHOW_CHASSIS_ENV_OPT,
		DEFAULT_SHOW_CHASSIS_ENV, `Run "show chassis environment" via NETCONF?`)

	jChecker.PersistentFlags().BoolVarP(
		&jCheckerConfig.ShowChassisFan,
		SHOW_CHASSIS_FAN_LOPT, SHOW_CHASSIS_FAN_OPT,
		DEFAULT_SHOW_CHASSIS_FAN, `Run "show chassis fan" via NETCONF?`)

	jChecker.PersistentFlags().BoolVarP(
		&jCheckerConfig.ShowChassisZones,
		SHOW_CHASSIS_ZONES_LOPT, SHOW_CHASSIS_ZONES_OPT,
		DEFAULT_SHOW_CHASSIS_ZONES, `Run "show chassis zones" via NETCONF?`)

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

	if jChecker.Flags().Lookup("help").Changed {
		os.Exit(0)
	}
}

func GetConfig() JCheckerConfig {
	return *jCheckerConfig
}
