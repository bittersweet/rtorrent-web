package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/bittersweet/rtorrent-web/util"
	"github.com/gorilla/mux"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleIndex")

	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatal("ParseFiles", err)
	}

	host := util.CurrentIP()
	hostName := fmt.Sprintf("http://%s", host)
	data := map[string]interface{}{
		"hostName": hostName,
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal("t.Execute", err)
	}
}

func handleTorrents(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleTorrents")

	torrents := getTorrents()
	for i := 0; i < len(torrents); i++ {
		t := &torrents[i]
		t.setTracker()
		t.Ratio = t.FormatRatio()
		t.PercentageDone = t.CalculateCompletedPercentage()
	}

	output, err := json.MarshalIndent(&torrents, "", "  ")
	if err != nil {
		log.Fatal("MarshalIndent", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(output)
}

func handleTorrentFiles(w http.ResponseWriter, r *http.Request) {
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

func handleStatic(w http.ResponseWriter, r *http.Request) {
	defer util.TrackTime(time.Now(), "handleStatic")

	http.ServeFile(w, r, r.URL.Path[1:])
}
