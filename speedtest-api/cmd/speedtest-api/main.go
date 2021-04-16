package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"os/exec"

	"log"

	"github.com/julienschmidt/httprouter"
)

const (
	layoutISO8601 = "2006-01-02T15:04:05Z"
)

type SpeedtestResult struct {
	TimeStamp time.Time `json:"timestamp"`
	Ping      Latency   `json:"ping"`
	Download  Stats     `json:"download"`
	Upload    Stats     `json:"upload"`
}

type Latency struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
}

type Stats struct {
	Bandwidth int `json:"bandwidth"`
}

type outputData struct {
	TimeStamp     string  `json:"timestamp"`
	Jitter        float64 `json:"jitter"`
	Latency       float64 `json:"latency"`
	DownBandwidth int     `json:"down_bandwidth"`
	UpBandwidth   int     `json:"up_bandwidth"`
}

func main() {

	router := httprouter.New()
	router.GET("/metrics", speedtestExport())

	log.Fatal(http.ListenAndServe(":9096", router))

}

func speedtestExport() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		result := performSpeedtest()
		list := printData(result)
		json.NewEncoder(w).Encode(list)
	}
}

func performSpeedtest() (result SpeedtestResult) {
	cmd := exec.Command("/bin/sh", "/exec/run")
	path := "/resources/report.json"

	if fileExists(path) {
		result = parseJson(path)
		loc, _ := time.LoadLocation("UTC")
		t := time.Now().In(loc)

		if t.Sub(result.TimeStamp).Seconds() > 600 {
			cmd.Run()
			return result
		}

		if t.Sub(result.TimeStamp).Seconds() > 60 {
			cmd.Start()
			return result
		}

		return result

	}
	cmd.Run()
	result = parseJson(path)
	return result
}

func parseJson(path string) (result SpeedtestResult) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("File missing")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &result)

	return result
}

func printData(result SpeedtestResult) outputData {
	list := outputData{
		TimeStamp:     changeFormat(result.TimeStamp),
		Jitter:        result.Ping.Jitter,
		Latency:       result.Ping.Latency,
		DownBandwidth: result.Download.Bandwidth,
		UpBandwidth:   result.Upload.Bandwidth,
	}
	return list
}

func changeFormat(t time.Time) string {
	output := t.Format(layoutISO8601)
	return output
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
