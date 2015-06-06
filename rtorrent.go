package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/bittersweet/rtorrent-web/util"

	"github.com/bittersweet/xmlrpc"
	"github.com/dustin/go-humanize"
	"github.com/gorilla/mux"
)

var client *xmlrpc.Client
var trackers = map[string]string{}

type Torrent struct {
	Name              string `json:"name"`
	BytesDone         string `json:"bytes_done"`
	BytesDoneRaw      int64
	ConnectionCurrent string `json:"connection_current"`
	CreationDate      int64  `json:"creation_date"`
	GetDownRate       string `json:"get_down_rate"`
	GetDownRateRaw    int64  `json:"get_down_rate_raw"`
	GetDownTotal      string `json:"get_down_total"`
	SizeBytes         string `json:"size_bytes"`
	SizeBytesRaw      int64
	SizeFiles         int64 `json:"size_files"`
	State             int64 `json:"state"`
	LoadDate          int64 `json:"load_date"`
	RatioRaw          int64
	Ratio             string `json:"ratio"`
	GetUpRate         string `json:"get_up_rate"`
	GetUpRateRaw      int64  `json:"get_up_rate_raw"`
	GetUpTotal        int64  `json:"get_up_total"`
	Hash              string `json:"hash"`
	PeersConnected    int64  `json:"peers_connected"`
	Tracker           string `json:"tracker"`
	PercentageDone    string `json:"percentage_done"`
}

func (t *Torrent) FormatRatio() string {
	result := float64(t.RatioRaw) / 1000
	return fmt.Sprintf("%.2f", result)
}

func (t *Torrent) CalculateCompletedPercentage() string {
	result := float64(t.BytesDoneRaw) / float64(t.SizeBytesRaw) * 100
	return fmt.Sprintf("%.0f%%", result)
}

type File struct {
	Name            string `json:"name"`
	SizeBytes       string `json:"size_bytes"`
	Priority        int64  `json:"priority"`
	IsOpen          int64  `json:"is_open"`
	CompletedChunks int64  `json:"completed_chunks"`
	FrozenPath      string `json:"get_frozen_path"`
}

type Tracker struct {
	hash string
	url  string
}

func (t *Torrent) getTracker() string {
	defer util.TrackTime(time.Now(), "getTracker")

	var tracker string
	if err := client.Call("t.get_url", []interface{}{t.Hash, 0}, &tracker); err != nil {
		log.Fatal("error: ", err)
	}

	url, err := url.Parse(tracker)
	if err != nil {
		panic(err)
	}

	host, port, _ := net.SplitHostPort(url.Host)
	if port == "" {
		host = url.Host
	}
	return host
	// Store an inmemory list of torrents, update that, and let the getTorrents
	// update fields that have changed
	// or for first version just store a map[hash] => tracker
	// Take a look at this concurrency model: https://golang.org/doc/codewalk/sharemem/
}

func (t *Torrent) setTracker() {
	// defer util.TrackTime(time.Now(), "setTracker")

	url := trackers[t.Hash]
	// fmt.Printf("Setting url to %s\\n", url)
	// if key is not in map, do a request
	t.Tracker = url
}

func getFiles(hash string) []File {
	var output [][]interface{}
	args := []interface{}{hash, 0, "f.get_path=", "f.size_bytes=", "f.get_priority=", "f.is_open=", "f.get_completed_chunks=", "f.get_frozen_path="}
	if err := client.Call("f.multicall", args, &output); err != nil {
		fmt.Println("d.multicall call error: ", err)
	}

	files := make([]File, len(output))
	for i := 0; i < len(output); i++ {
		data := output[i]

		fmt.Printf("%#v\n", data)

		file := File{
			Name:            data[0].(string),
			SizeBytes:       humanize.Bytes(uint64(data[1].(int64))),
			Priority:        data[2].(int64),
			IsOpen:          data[3].(int64),
			CompletedChunks: data[4].(int64),
			FrozenPath:      data[5].(string),
		}
		files[i] = file
	}

	return files
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleIndex")

	http.ServeFile(w, r, "static/index.html")
}

func rtorrentMethodFromString(status string) string {
	switch status {
	case "stop":
		return "d.stop"
	case "start":
		return "d.start"
	case "remove":
		return "d.erase"
	}

	return "d.start"
}

func changeStatus(hash string, status string) {
	method := rtorrentMethodFromString(status)
	var statusCode int
	if err := client.Call(method, []interface{}{hash, 0}, &statusCode); err != nil {
		log.Fatal("error: ", err)
	}

	fmt.Printf("%#v\n", statusCode)
}

func copyTorrent(hash string) {
	files := getFiles(hash)

	destinationFolder := os.Getenv("DESTINATIONFOLDER")

	for _, file := range files {
		source := file.FrozenPath
		target := fmt.Sprintf("%s/%s", destinationFolder, file.Name)

		copyFile(source, target)
	}
}

func copyFile(source string, target string) {
	mvCmd := exec.Command("mv", source, target)
	err := mvCmd.Run()

	fmt.Printf("Copying %s to %s\n", source, target)
	if err != nil {
		log.Fatal("Copying file failure: ", err)
	}
}

func handleTorrent(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrent")

	http.ServeFile(w, r, "static/show.html")
}

func handleTorrentChangeStatus(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrentChangeStatus")

	vars := mux.Vars(r)
	hash := vars["hash"]

	u, _ := url.Parse(r.URL.String())
	queryParams := u.Query()
	requestedStatus := queryParams["status"][0]
	changeStatus(hash, requestedStatus)

	response := map[string]string{
		"status": "ok",
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTorrentCopy(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrentCopy")

	vars := mux.Vars(r)
	hash := vars["hash"]

	copyTorrent(hash)

	response := map[string]string{
		"status": "ok",
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTorrents(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrents")

	torrents := getTorrents()
	for i := 0; i < len(torrents); i++ {
		torrent := torrents[i]
		torrent.setTracker()
		// Work around to pass by value or pointer type thing
		// updating in setTracker didn't work
		torrents[i] = torrent
		torrents[i].Ratio = torrents[i].FormatRatio()
		torrents[i].PercentageDone = torrents[i].CalculateCompletedPercentage()
	}

	output, err := json.MarshalIndent(&torrents, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTorrentJson(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrentJson")

	vars := mux.Vars(r)
	files := getFiles(vars["hash"])

	output, err := json.MarshalIndent(&files, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent files", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTrackers(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTrackers")

	output, err := json.MarshalIndent(&trackers, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleStatic")

	http.ServeFile(w, r, r.URL.Path[1:])
}

func stateMonitor() chan<- *Tracker {
	updates := make(chan *Tracker)
	go func() {
		for {
			select {
			case t := <-updates:
				trackers[t.hash] = t.url
			}
		}
	}()

	fmt.Println("Statemonitor launched")
	return updates
}

func init() {
	defer util.TrackTime(time.Now(), "init")

	var err error
	client, err = xmlrpc.NewClient("http://192.168.2.7:5001/RPC2", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func getTorrents() []Torrent {
	defer util.TrackTime(time.Now(), "getTorrents")

	var output [][]interface{}
	if err := client.Call("d.multicall", []interface{}{"main", "d.name=", "d.bytes_done=", "d.connection_current=", "d.creation_date=", "d.get_down_rate=", "d.get_down_total=", "d.size_bytes=", "d.size_files=", "d.state=", "d.load_date=", "d.ratio=", "d.get_up_rate=", "d.get_up_total=", "d.hash=", "d.peers_connected="}, &output); err != nil {
		fmt.Println("d.multicall call error: ", err)
	}

	torrents := make([]Torrent, len(output))
	for i := 0; i < len(output); i++ {
		data := output[i]

		torrent := Torrent{
			Name:              data[0].(string),
			BytesDone:         humanize.Bytes(uint64(data[1].(int64))),
			BytesDoneRaw:      data[1].(int64),
			ConnectionCurrent: data[2].(string),
			CreationDate:      data[3].(int64),
			GetDownRate:       humanize.Bytes(uint64(data[4].(int64))),
			GetDownRateRaw:    data[4].(int64),
			GetDownTotal:      humanize.Bytes(uint64(data[5].(int64))),
			SizeBytes:         humanize.Bytes(uint64(data[6].(int64))),
			SizeBytesRaw:      data[6].(int64),
			SizeFiles:         data[7].(int64),
			State:             data[8].(int64),
			LoadDate:          data[9].(int64),
			RatioRaw:          data[10].(int64),
			GetUpRate:         humanize.Bytes(uint64(data[11].(int64))),
			GetUpRateRaw:      data[11].(int64),
			GetUpTotal:        data[12].(int64),
			Hash:              data[13].(string),
			PeersConnected:    data[14].(int64),
		}
		torrents[i] = torrent
	}

	return torrents
}

func Poller(process <-chan *Torrent, updates chan<- *Tracker) {
	for torrent := range process {
		tracker := torrent.getTracker()
		updates <- &Tracker{torrent.Hash, tracker}
	}
}

func main() {
	updates := stateMonitor()
	incoming := make(chan *Torrent)

	torrents := getTorrents()

	// Launch poller goroutines
	for i := 0; i < 4; i++ {
		go Poller(incoming, updates)
	}

	// Send all torrents to the queue to be processed
	for i := 0; i < len(torrents); i++ {
		torrent := torrents[i]
		incoming <- &torrent
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/torrents.json", handleTorrents)
	mux.HandleFunc("/torrents/{hash}.json", handleTorrentJson)
	mux.HandleFunc("/torrents/{hash}", handleTorrent)
	mux.HandleFunc("/torrents/{hash}/changestatus", handleTorrentChangeStatus)
	mux.HandleFunc("/torrents/{hash}/copy", handleTorrentCopy)
	mux.HandleFunc("/trackers", handleTrackers)
	mux.HandleFunc("/static/{file}", handleStatic)
	fmt.Println("Will start listening on port 8000")
	http.ListenAndServe(":8000", mux)
}
