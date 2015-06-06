package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/bittersweet/rtorrent-web/util"
	"github.com/gorilla/mux"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleIndex")

	http.ServeFile(w, r, "static/index.html")
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
