package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/kolo/xmlrpc"
)

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
}

type Torrent struct {
	Name              string `json:"name"`
	BytesDone         string `json:"bytes_done"`
	ConnectionCurrent string `json:"connection_current"`
	CreationDate      int64  `json:"creation_date"`
	GetDownRate       string `json:"get_down_rate"`
	GetDownTotal      string `json:"get_down_total"`
	SizeBytes         string `json:"size_bytes"`
	SizeFiles         int64  `json:"size_files"`
	State             int64  `json:"state"`
	LoadDate          int64  `json:"load_date"`
	Ratio             int64  `json:"ratio"`
	GetUpRate         string `json:"get_up_rate"`
	GetUpTotal        string `json:"get_up_total"`
	Hash              string `json:"hash"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleIndex")

	torrents := getTorrents()

	output, err := json.MarshalIndent(&torrents, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

var client *xmlrpc.Client

func init() {
	defer trackTime(time.Now(), "init")

	var err error
	client, err = xmlrpc.NewClient("http://192.168.2.7:5001/RPC2", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func getTorrents() []Torrent {
	defer trackTime(time.Now(), "getTorrents")

	var output [][]interface{}
	if err := client.Call("d.multicall", []interface{}{"main", "d.name=", "d.bytes_done=", "d.connection_current=", "d.creation_date=", "d.get_down_rate=", "d.get_down_total=", "d.size_bytes=", "d.size_files=", "d.state=", "d.load_date=", "d.ratio=", "d.get_up_rate=", "d.get_up_total=", "d.hash="}, &output); err != nil {
		fmt.Println("d.multicall call error: ", err)
	}
	fmt.Printf("Found %d torrents\n", len(output))

	torrents := make([]Torrent, len(output))
	for i := 0; i < len(output); i++ {
		data := output[i]

		torrent := Torrent{
			Name:              data[0].(string),
			BytesDone:         humanize.Bytes(uint64(data[1].(int64))),
			ConnectionCurrent: data[2].(string),
			CreationDate:      data[3].(int64),
			GetDownRate:       humanize.Bytes(uint64(data[4].(int64))),
			GetDownTotal:      humanize.Bytes(uint64(data[5].(int64))),
			SizeBytes:         humanize.Bytes(uint64(data[6].(int64))),
			SizeFiles:         data[7].(int64),
			State:             data[8].(int64),
			LoadDate:          data[9].(int64),
			Ratio:             data[10].(int64),
			GetUpRate:         humanize.Bytes(uint64(data[11].(int64))),
			GetUpTotal:        humanize.Bytes(uint64(data[12].(int64))),
			Hash:              data[13].(string),
		}
		torrents[i] = torrent
	}

	fmt.Printf("%#v\n", torrents[0])
	return torrents
}

func main() {
	defer trackTime(time.Now(), "xmlRpc")

	http.HandleFunc("/", handleIndex)
	fmt.Println("Will start listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
