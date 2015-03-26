package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/kolo/xmlrpc"
)

var client *xmlrpc.Client
var trackers = map[string]string{}

// var updates = make(chan Tracker)

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
	GetUpRateRaw      int64
	GetUpTotal        string `json:"get_up_total"`
	Hash              string `json:"hash"`
	Tracker           string `json:"tracker"`
}

func (t *Torrent) getTracker() string {
	defer trackTime(time.Now(), "getTracker")

	var tracker string
	if err := client.Call("t.get_url", []interface{}{t.Hash, 0}, &tracker); err != nil {
		log.Fatal("error: ", err)
	}

	url, err := url.Parse(tracker)
	if err != nil {
		panic(err)
	}

	fmt.Println(url)
	host, _, _ := net.SplitHostPort(url.Host)
	return host
	// Store an inmemory list of torrents, update that, and let the getTorrents
	// update fields that have changed
	// or for first version just store a map[hash] => tracker
	// Take a look at this concurrency model: https://golang.org/doc/codewalk/sharemem/
}

func (t *Torrent) setTracker() {
	// defer trackTime(time.Now(), "setTracker")

	url := trackers[t.Hash]
	// fmt.Printf("Setting url to %s\\n", url)
	// if key is not in map, do a request
	t.Tracker = url
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleIndex")

	http.ServeFile(w, r, "static/index.html")
}

func handleTorrents(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleIndex")

	torrents := getTorrents()
	for i := 0; i < len(torrents); i++ {
		torrent := torrents[i]
		torrent.setTracker()
		torrents[i] = torrent
	}

	output, err := json.MarshalIndent(&torrents, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTrackers(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "handleTrackers")

	output, err := json.MarshalIndent(&trackers, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

type Tracker struct {
	hash string
	url  string
}

func stateMonitor() chan<- *Tracker {
	updates := make(chan *Tracker)
	go func() {
		for {
			select {
			case t := <-updates:
				fmt.Println("update")
				trackers[t.hash] = t.url
			}
		}
	}()

	fmt.Println("Statemonitor launched")
	return updates
}

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
			GetUpRateRaw:      data[11].(int64),
			GetUpTotal:        humanize.Bytes(uint64(data[12].(int64))),
			Hash:              data[13].(string),
		}
		torrents[i] = torrent
	}

	fmt.Printf("%#v\n", torrents[0])
	return torrents
}

func Poller(process <-chan *Torrent, updates chan<- *Tracker) {
	for torrent := range process {
		tracker := torrent.getTracker()
		fmt.Println("updating")
		updates <- &Tracker{torrent.Hash, tracker}
	}
}

func main() {
	defer trackTime(time.Now(), "xmlRpc")

	updates := stateMonitor()
	incoming := make(chan *Torrent)

	torrents := getTorrents()
	fmt.Println("have torrents", len(torrents))

	// Launch poller goroutines
	for i := 0; i < 4; i++ {
		go Poller(incoming, updates)
	}

	// Send all torrents to the queue to be processed
	for i := 0; i < len(torrents); i++ {
		torrent := torrents[i]
		fmt.Println("updating")
		incoming <- &torrent
	}
	fmt.Printf("%#v\n", trackers)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/torrents", handleTorrents)
	http.HandleFunc("/trackers", handleTrackers)
	fmt.Println("Will start listening on port 8000")
	http.ListenAndServe(":8000", nil)
}
