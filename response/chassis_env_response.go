package response

import (
	"encoding/csv"
	"encoding/xml"
	"io"
	"time"
)

/*var (
	recordChan chan *csvWriteJob
)*/

//func init() {
	//recordChan = make(chan *csvWriteJob)
	//go writeCSV()
//}

type EnvItem struct {
	Name   string `xml:"name,omitempty"   json:"name,omitempty"`
	Class  string `xml:"class,omitempty"  json:"class,omitempty"`
	Status string `xml:"status,omitempty" json:"status,omitempty"`
}

type ChassisEnvResponse struct {
	XMLName   xml.Name  `xml:"environment-information,omitempty" json:"-"`
	EnvItem   []EnvItem `xml:"environment-item,omitempty"        json:"environment-item,omitempty"`
	csvWriter *csv.Writer
}

func NewChassisEnvResponse(ncReply *string) (*ChassisEnvResponse, error) {

	var (
		replyNoNewlines    = newLine.Replace(*ncReply)
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

func (cer *ChassisEnvResponse) WriteCSV(w io.Writer) {

	record := []string{}
	records := [][]string{}
	tStamp := time.Now()

	for _, envItem := range cer.EnvItem {
		record = append(record, tStamp.Format(time.Stamp))
		record = append(record, envItem.Name)
		record = append(record, envItem.Class)
		record = append(record, envItem.Status)
		records = append(records, record)
	}

	csvWriter := csv.NewWriter(w)
	csvWriter.WriteAll(records)
	//recordChan <- newCSVWriteJob(csvWriter, records)
}

/*func writeCSV() {
	for rec := range recordChan {
		rec.WriteAll()
	}
}*/
