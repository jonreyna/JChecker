package main

import (
	"github.com/JReyLBC/JChecker/command"
	"github.com/JReyLBC/JChecker/config"
	"github.com/JReyLBC/JChecker/reply"
	log "github.com/Sirupsen/logrus"

	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func init() {
	config.Execute()
}

func main() {
	const thisFunc = "main.main()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})

	cfg := config.GetConfig()

	cntxLog.WithFields(log.Fields{
		"file": cfg.CommandConfigCSVFile,
	})

	if file, err := os.Open(cfg.CommandConfigCSVFile); err != nil {
		cntxLog.Errorln("could not open command config CSV file")
		cntxLog.Fatalln(err)
	} else if records, err := readCSVRecs(file); err != nil {
		cntxLog.Errorln("could not read command records")
		cntxLog.Fatalln(err)
	} else if cmds, err := buildCommands(records); err != nil {
		cntxLog.Errorln("could not construct commands from records")
		cntxLog.Fatal(err)
	} else {
		var wg sync.WaitGroup
		wg.Add(len(cmds))
		runCommands(cmds, &wg)
		wg.Wait()
	}
}

func readCSVRecs(file io.Reader) ([][]string, error) {

	const thisFunc = "main.readCSVRecs()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})

	cntxLog.Debugln("setting up CSV reader")
	csvReader := csv.NewReader(file)
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = command.NUM_FIELDS
	csvReader.TrimLeadingSpace = true

	cntxLog.Debugln("reading all records")
	return csvReader.ReadAll()
}

func buildCommands(records [][]string) ([]command.Commander, error) {

	const thisFunc = "main.buildCommands()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})

	cmds := []command.Commander{}
	for recNum, record := range records {
		record[command.COMMAND_FIELD] = strings.ToLower(record[command.COMMAND_FIELD])

		switch record[command.COMMAND_FIELD] {
		case command.CHASSIS_ENV_CMD:
			cntxLog = cntxLog.WithFields(log.Fields{
				"command": command.CHASSIS_ENV_CMD,
			})
			if cmd, err := command.ChassisEnvCmdFromRecord(record); err != nil {
				cntxLog.Debugln("could not create ChassisEnvCmd")
				cntxLog.Debugln(err)
				return nil, err
			} else {
				cntxLog.Debugln("appending command")
				cmds = append(cmds, cmd)
			}
		default:
			cntxLog = cntxLog.WithFields(log.Fields{
				"command": record[command.COMMAND_FIELD],
			})
			cntxLog.Debugln("unknown command")
			return nil, fmt.Errorf("unknown command '%s' in record %d", record[command.COMMAND_FIELD], recNum+1)
		}
	}
	return cmds, nil
}

func runCommands(cmds []command.Commander, wg *sync.WaitGroup) {

	for _, cmd := range cmds {
		replyChan, errChan := reply.ProcessChassisEnv(cmd.Run())

		go func() {
			const thisFunc = "main.runCommands.func()"
			cntxLog := log.WithFields(log.Fields{
				"func": thisFunc,
			})

			defer wg.Done()
			done := false

			for !done {
				select {
				case r, ok := <-replyChan:
					if !ok {
						cntxLog.Debugln("replyChan closed")
						done = true
						break
					} else if outFile, err := cmd.OutFile(); err != nil {
						cntxLog.Errorln(err)
					} else if err := r.WriteCSV(outFile); err != nil {
						cntxLog.Errorln(err)
					}
				case e, ok := <-errChan:
					if !ok {
						cntxLog.Debugln("errChan closed")
						done = true
						break
					} else {
						log.Errorln(e)
					}
				}
			}
		}()
	}
}
