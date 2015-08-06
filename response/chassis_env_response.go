package response

import (
	"encoding/csv"
	"encoding/xml"
	"time"

	"github.com/Juniper/go-netconf/netconf"
)

/*var (
	recordChan chan *csvWriteJob
)*/

//func init() {
//recordChan = make(chan *csvWriteJob)
//go writeCSV()
//}

type EnvItem struct {
	Name        string `xml:"name,omitempty"   json:"name,omitempty"`
	Class       string `xml:"class,omitempty"  json:"class,omitempty"`
	Status      string `xml:"status,omitempty" json:"status,omitempty"`
	Temperature string `xml:"temperature,omitempty" json:"temperature,omitempty"`
}

type ChassisEnvResponse struct {
	XMLName xml.Name  `xml:"environment-information,omitempty" json:"-"`
	EnvItem []EnvItem `xml:"environment-item,omitempty"        json:"environment-item,omitempty"`
}

func NewChassisEnvResponse(ncReply *netconf.RPCReply) (*ChassisEnvResponse, error) {

	var (
		replyNoNewlines    = newLine.Replace(ncReply.Data)
		chassisEnvResponse = new(ChassisEnvResponse)
	)

	if err := xml.Unmarshal([]byte(replyNoNewlines), chassisEnvResponse); err != nil {
		return nil, err
	} else {
		return chassisEnvResponse, nil
	}
}

/*type csvWriteJob struct {
	CSVWriter      *csv.Writer
	RecordsToWrite [][]string
}

func (j *csvWriteJob) WriteAll() {
	j.CSVWriter.WriteAll(j.RecordsToWrite)
}*/

/*func newCSVWriteJob(writer *csv.Writer, records [][]string) *csvWriteJob {
	return &csvWriteJob{
		CSVWriter:      writer,
		RecordsToWrite: records,
	}
}*/

func (cer *ChassisEnvResponse) WriteCSV(w *csv.Writer) error {

	records := make([][]string, 0, 32)
	tStamp := time.Now()

	for _, envItem := range cer.EnvItem {
		record := make([]string, 0, 4)
		if envItem.Class == "Power" {
			continue
		}
		record = append(record, tStamp.Format(time.Stamp))
		record = append(record, envItem.Name)
		//record = append(record, envItem.Class)
		record = append(record, envItem.Status)
		record = append(record, envItem.Temperature)
		records = append(records, record)
	}

	if err := w.WriteAll(records); err != nil {
		return err
	} else {
		return nil
	}
	//recordChan <- newCSVWriteJob(csvWriter, records)
}

/*func writeCSV() {
	for rec := range recordChan {
		rec.WriteAll()
	}
}*/
