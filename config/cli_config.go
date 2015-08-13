// Package config maintains the global configuration
// required to run JChecker. Configurations can be
// set on the command line or in a config file.
package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// Default options
const (
	DEFAULT_LOG_LEVEL = log.InfoLevel

	DEFAULT_COMMAND_CSV_CONFIG_FILE = "jchecker_command_config.csv"
	COMMAND_CSV_CONFIG_LOPT         = "command-config"
	COMMAND_CSV_CONFIG_OPT          = "c"
	COMMAND_CSV_CONFIG_DESCRIPTION  = "Location of the file with CSV commands."
)

// Contains all config options for JChecker
type JCheckerConfig struct {
	CommandConfigCSVFile string
}

// Package scope variables to handle
// config file access and overall config.
var (
	jCheckerConfig *JCheckerConfig
	jChecker       *cobra.Command
)

// Initialize config and commands
func init() {

	const thisFunc = "config.init()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})

	// Check JCHECKER_LOG_LEVEL and set logger accordingly
	switch os.Getenv("JCHECKER_LOG_LEVEL") {
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

	cntxLog.Debugln("setting JChecker main command usage")

	jCheckerConfig = new(JCheckerConfig)
	jChecker = &cobra.Command{
		Use:   "JChecker",
		Short: "JChecker executes requested commands on Juniper devices at user defined time intervals.",

		Long: `
JChecker executes requested commands at user defined time intervals, until
a user defined deadline is reached. Deadlines, time intervals, and device
definitions are defined in a user provided CSV file (the command CSV file).

It's format is:
command, ip, interval, deadline, results file, username, password

A comment character '#' at the beginning of the line, causes the parser to
ignore the entire line.

Commands currently supported include:
  - show chassis environment`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				log.Errorf("Unknown command(s): %v", args)
				cmd.Help()
				os.Exit(1)
			}
		},
	}

	jChecker.PersistentFlags().StringVarP(
		&jCheckerConfig.CommandConfigCSVFile,
		COMMAND_CSV_CONFIG_LOPT, COMMAND_CSV_CONFIG_OPT,
		DEFAULT_COMMAND_CSV_CONFIG_FILE, COMMAND_CSV_CONFIG_DESCRIPTION,
	)

	// This message will be shown to Windows users if JChecker is opened from explorer.exe
	cobra.MousetrapHelpText = `

	JChecker is a command line application.

	You need to open cmd.exe and run it from there.
	`
}

// Get the config options
func Execute() {

	const thisFunc = "config.Execute()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})
	cntxLog.Debugln("executing CLI config parser")

	if err := jChecker.Execute(); err != nil {
		log.Fatalln(err)
	}
	if jChecker.PersistentFlags().Lookup("help").Changed {
		os.Exit(0)
	}

	if jChecker.HasSubCommands() {
		log.Debugln("Has subcommand")
	}

}

func GetConfig() JCheckerConfig {
	return *jCheckerConfig
}
