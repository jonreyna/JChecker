package response

import (
	"encoding/xml"
)

type Zone struct {
	Name           string `xml:"name,omitempty"             json:"name,omitempty"`
	DrivingFRUName string `xml:"driving-fru-name,omitempty" json:"driving-fru-name,omitempty"`
	Temperature    string `xml:"temperature,omitempty"      json:"temperature,omitempty"`
	ZoneStatus     string `xml:"zone-status,omitempty"      json:"zone-status,omitempty"`
	FanMissingCnt  int    `xml:"fan-missing-cnt,omitempty"  json:"fan-missing-cnt,omitempty"`
	FanFailedCnt   int    `xml:"fan-failed-cnt,omitempty"   json:"fan-failed-cnt,omitempty"`
	FanDutyCycle   int    `xml:"fan-dutycycle,omitempty"    json:"fan-dutycycle,omitempty"`
}
type ChassisZonesResponse struct {
	XMLName xml.Name `xml:"chassis-zones-information,omitempty" json:"-"`
	Zones   []Zone
}

func NewZonesResponse(ncReply *string) (*ChassisZonesResponse, error) {

	var (
		replyNoNewlines    = newLine.Replace(*ncReply)
		chassisZonesResponse = new(ChassisZonesResponse)
	)

	if err := xml.Unmarshal([]byte(replyNoNewlines), chassisZonesResponse); err != nil {
		return nil, err
	} else {
		return chassisZonesResponse, nil
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

/*func (cer *ChassisZonesResponse) WriteCSV(w io.Writer) {

	record := []string{}
	records := [][]string{}
	tStamp := time.Now()

	for _, envItem := range cer.Item {
		record = append(record, tStamp.Format(time.Stamp))
		record = append(record, envItem.Name)
		record = append(record, envItem.Class)
		record = append(record, envItem.Status)
		records = append(records, record)
	}

	csvWriter := csv.NewWriter(w)
	csvWriter.WriteAll(records)
	//recordChan <- newCSVWriteJob(csvWriter, records)
}*/

/*func writeCSV() {
	for rec := range recordChan {
		rec.WriteAll()
	}
}*/
