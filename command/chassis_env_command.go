package command

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

const (
	CHASSIS_ENV_CMD_TMPL = "<get-environment-information></get-environment-information>"
)

type ChassisEnvCmd struct {
	// NETCONF credentials
	UserName string
	Password string

	outFile string

	// Where to execute the command
	LoopBackIP net.IP

	// The interval to execute the command
	Interval time.Duration
	// The deadline to stop executing the command
	Deadline time.Duration

	// To cancel a command prematurely
	cancelled chan struct{}
}

func ChassisEnvCmdFromRecord(record []string) (*ChassisEnvCmd, error) {

	const thisFunc = "command.ChassisEnvCmdFromRecord()"
	cntxLog := log.WithFields(log.Fields{
		"func": thisFunc,
	})

	if len(record) != NUM_FIELDS {
		cntxLog = cntxLog.WithFields(log.Fields{
			"num_fields": len(record),
		})
		cntxLog.Debugf("number of fields should be %d", NUM_FIELDS)

		return nil, &Error{
			Method:    thisFunc,
			ClientErr: fmt.Sprintf("number of fields in record is %d when it should be %d", len(record), NUM_FIELDS),
			Err:       fmt.Errorf("number of fields in record is %d when it should be %d", len(record), NUM_FIELDS),
		}
	} else if ip := net.ParseIP(record[IP_FIELD]); ip == nil {
		cntxLog = cntxLog.WithFields(log.Fields{
			"field":      record[IP_FIELD],
			"field_num":  IP_FIELD,
			"field_name": FieldNameMap[IP_FIELD],
		})
		cntxLog.Debugln("net.ParseIP() returned nil")

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: IP_FIELD,
			ClientErr:   fmt.Sprintf("could not parse IP [field %d] from record", IP_FIELD),
			Err:         fmt.Errorf("could not parse IP [field %d] from record", IP_FIELD),
		}
	} else if interval, err := time.ParseDuration(record[INTERVAL_FIELD]); err != nil {
		cntxLog = cntxLog.WithFields(log.Fields{
			"interval":   record[INTERVAL_FIELD],
			"field_num":  INTERVAL_FIELD,
			"field_name": FieldNameMap[INTERVAL_FIELD],
		})
		cntxLog.Debugln("could not parse interval duration")
		cntxLog.Debugln(err)

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: INTERVAL_FIELD,
			ClientErr:   fmt.Sprintf("could not parse interval [field %d] duration from record", INTERVAL_FIELD),
			Err:         fmt.Errorf("could not parse interval [field %d] duration from record", INTERVAL_FIELD),
		}
	} else if deadline, err := time.ParseDuration(record[DEADLINE_FIELD]); err != nil {
		cntxLog = cntxLog.WithFields(log.Fields{
			"deadline":   record[DEADLINE_FIELD],
			"field_num":  DEADLINE_FIELD,
			"field_name": FieldNameMap[DEADLINE_FIELD],
		})
		cntxLog.Debugln("could not parse deadline duration")
		cntxLog.Debugln(err)

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: DEADLINE_FIELD,
			ClientErr:   fmt.Sprintf("could not parse deadline [field %d] duration from record", DEADLINE_FIELD),
			Err:         err,
		}
	} else if err := outFileMap.addFile(record[OUT_FILE_FIELD]); err != nil {
		cntxLog = cntxLog.WithFields(log.Fields{
			"file":       record[OUT_FILE_FIELD],
			"field_num":  OUT_FILE_FIELD,
			"field_name": FieldNameMap[OUT_FILE_FIELD],
		})
		cntxLog.Debugln("error creating output file")
		cntxLog.Debugln(err)

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: OUT_FILE_FIELD,
			ClientErr:   fmt.Sprintf("could not create output file %s [field %d]", record[OUT_FILE_FIELD], OUT_FILE_FIELD),
			Err:         err,
		}
	} else if record[USER_NAME_FIELD] == "" {
		cntxLog.WithFields(log.Fields{
			"field_num":  USER_NAME_FIELD,
			"field_name": FieldNameMap[USER_NAME_FIELD],
		}).Debugln("user name field empty")

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: USER_NAME_FIELD,
			ClientErr:   fmt.Sprintf("user name [field %d] is empty", USER_NAME_FIELD),
			Err:         fmt.Errorf("user name [field %d] is empty", USER_NAME_FIELD),
		}
	} else if record[PASSWORD_FIELD] == "" {
		cntxLog = cntxLog.WithFields(log.Fields{
			"field_num":  PASSWORD_FIELD,
			"field_name": FieldNameMap[PASSWORD_FIELD],
		})
		cntxLog.Debugln("password field empty")

		return nil, &ParseError{
			Method:      thisFunc,
			Record:      record,
			RecordField: PASSWORD_FIELD,
			ClientErr:   fmt.Sprintf("password [field %d] is empty", PASSWORD_FIELD),
			Err:         fmt.Errorf("password [field %d] is empty", PASSWORD_FIELD),
		}
	} else {
		cntxLog = cntxLog.WithFields(log.Fields{
			"command":  record[COMMAND_FIELD],
			"ip":       ip,
			"interval": interval,
			"deadline": deadline,
			"out_file": record[OUT_FILE_FIELD],
			"username": record[USER_NAME_FIELD],
		})
		cntxLog.Debugln("creating ChassisEnvCmd object")

		return &ChassisEnvCmd{
			UserName:   record[USER_NAME_FIELD],
			Password:   record[PASSWORD_FIELD],
			outFile:    record[OUT_FILE_FIELD],
			LoopBackIP: ip,
			Interval:   interval,
			Deadline:   deadline,
			cancelled:  make(chan struct{}, 1),
		}, nil
	}
}

func (*ChassisEnvCmd) Method() netconf.RawMethod {
	return netconf.RawMethod(CHASSIS_ENV_CMD_TMPL)
}

func (cer *ChassisEnvCmd) Run() (<-chan *netconf.RPCReply, <-chan error) {

	replyChan := make(chan *netconf.RPCReply)
	errChan := make(chan error)

	go func() {
		const thisFunc = "*ChassisEnvCmd.Run()"
		cntxLog := log.WithFields(log.Fields{
			"func":     thisFunc,
			"deadline": cer.Deadline,
			"interval": cer.Interval,
			"username": cer.UserName,
			"ip":       cer.LoopBackIP,
		})

		cntxLog.Debugln("Starting")

		// Flag for breaking loop
		done := false

		// Use a ticker to notify us when it's time to execute the command again.
		ticker := time.NewTicker(cer.Interval)
		// Use time.After to let us know when we've reached our deadline.
		deadline := time.After(cer.Deadline)

		for !done {
			select {
			case <-ticker.C:
				cntxLog.Debugln("dialing SSH")
				if s, err := netconf.DialSSH(cer.LoopBackIP.String(),
					netconf.SSHConfigPassword(cer.UserName, cer.Password)); err != nil {
					cntxLog.Debugln(err)
					errChan <- &RunError{
						Method:    thisFunc,
						ClientErr: fmt.Sprintf("could not dial %s via NETCONF over SSH", cer.LoopBackIP.String()),
						Err:       err,
						Cmd:       cer,
					}
				} else if ncReply, err := s.Exec(cer.Method()); err != nil {
					cntxLog.Debugln(err)
					s.Close()
					errChan <- &RunError{
						Method:    thisFunc,
						ClientErr: fmt.Sprintf("could not execute show chassis enviornment command"),
					}
				} else {
					s.Close()
					replyChan <- ncReply
				}
			case <-deadline:
				cntxLog.Debugln("deadline reached")
				close(replyChan)
				close(errChan)
				done = true
			case <-cer.cancelled:
				cntxLog.Debugln("command cancelled")
				close(replyChan)
				close(errChan)
				done = true
			}
		}

		cntxLog.Debugln("Exiting")
	}()

	return replyChan, errChan
}

func (cer *ChassisEnvCmd) Cancel() {
	const thisMethod = "*ChassisEnvCmd.Cancel()"
	cntxLog := log.WithFields(log.Fields{
		"addr":   fmt.Sprintf("%p", cer),
		"method": thisMethod,
	})
	cntxLog.Debugln("cancelling command execution")
	cer.cancelled <- struct{}{}
}

func (cer *ChassisEnvCmd) OutFile() (*os.File, error) {
	const thisMethod = "*ChassisEnvCmd.OutFile()"
	cntxLog := log.WithFields(log.Fields{
		"addr":   fmt.Sprintf("%p", cer),
		"method": thisMethod,
	})
	cntxLog.Debugln("returning output file")
	return outFileMap.getFile(cer.outFile)
}
