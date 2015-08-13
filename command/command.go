package command

import (
	"fmt"
	"os"
	"sync"

	"github.com/Juniper/go-netconf/netconf"
	log "github.com/Sirupsen/logrus"
)

const (
	// Command fields
	COMMAND_FIELD = iota
	IP_FIELD
	INTERVAL_FIELD
	DEADLINE_FIELD
	OUT_FILE_FIELD
	USER_NAME_FIELD
	PASSWORD_FIELD
	NUM_FIELDS
)

const (
	// Available commands
	CHASSIS_ENV_CMD = "show chassis environment"
)

type outputFileMap struct {
	m       sync.RWMutex
	fileMap map[string]*os.File
}

var (
	outFileMap   *outputFileMap
	FieldNameMap map[int]string
)

func init() {
	outFileMap = new(outputFileMap)
	outFileMap.fileMap = make(map[string]*os.File)

	FieldNameMap = make(map[int]string)
	FieldNameMap[COMMAND_FIELD] = "command"
	FieldNameMap[IP_FIELD] = "ip"
	FieldNameMap[INTERVAL_FIELD] = "interval"
	FieldNameMap[DEADLINE_FIELD] = "deadline"
	FieldNameMap[OUT_FILE_FIELD] = "out_file"
	FieldNameMap[USER_NAME_FIELD] = "username"
	FieldNameMap[PASSWORD_FIELD] = "password"
	FieldNameMap[NUM_FIELDS] = "num_fields"
}

func (ofm *outputFileMap) addFile(fileName string) error {
	ofm.m.Lock()
	defer ofm.m.Unlock()
	const thisMethod = "*outputFileMap.addFile()"

	cntxLog := log.WithFields(log.Fields{
		"addr":      fmt.Sprintf("%p", ofm),
		"method":    thisMethod,
		"file_name": fileName,
	})
	cntxLog.Debugln("attempting to add file")

	if _, ok := ofm.fileMap[fileName]; !ok {
		cntxLog.Debugln("file not in map")
		if file, err := os.Create(fileName); err != nil {
			return err
		} else {
			cntxLog.Debugln("adding file to map")
			ofm.fileMap[fileName] = file
			return nil
		}
	} else {
		cntxLog.Debugln("file is already in map")
		return nil
	}
}

func (ofm *outputFileMap) getFile(fileName string) (*os.File, error) {
	ofm.m.RLock()
	defer ofm.m.RUnlock()
	const thisMethod = "*outputFileMap.getFile()"

	cntxLog := log.WithFields(log.Fields{
		"addr":      fmt.Sprintf("%p", ofm),
		"method":    thisMethod,
		"file_name": fileName,
	})
	cntxLog.Debugln("attempting to get file")

	if file, ok := ofm.fileMap[fileName]; !ok {
		cntxLog.Debugln("file is not in map")
		return nil, fmt.Errorf("Error getting file from map")
	} else {
		cntxLog.Debugln("returning file")
		return file, nil
	}
}

type Commander interface {
	Run() (<-chan *netconf.RPCReply, <-chan error)
	OutFile() (*os.File, error)
	Cancel()
}
