package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kolo/xmlrpc"
)

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s\n", name, elapsed)
}

type Torrent struct {
	Name              string
	BytesDone         int64
	ConnectionCurrent string
	CreationDate      int64
	GetDownRate       int64
	GetDownTotal      int64
	SizeBytes         int64
	SizeFiles         int64
	State             int64
	LoadDate          int64
	Ratio             int64
	GetUpRate         int64
	GetUpTotal        int64
	Hash              string
}

func main() {
	defer trackTime(time.Now(), "xmlRpc")
	client, err := xmlrpc.NewClient("http://192.168.2.7:5001/RPC2", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}

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
			BytesDone:         data[1].(int64),
			ConnectionCurrent: data[2].(string),
			CreationDate:      data[3].(int64),
			GetDownRate:       data[4].(int64),
			GetDownTotal:      data[5].(int64),
			SizeBytes:         data[6].(int64),
			SizeFiles:         data[7].(int64),
			State:             data[8].(int64),
			LoadDate:          data[9].(int64),
			Ratio:             data[10].(int64),
			GetUpRate:         data[11].(int64),
			GetUpTotal:        data[12].(int64),
			Hash:              data[13].(string),
		}
		// fmt.Println(torrent)
		// fmt.Println(i)
		// torrents = append(torrents, torrent)
		torrents[i] = torrent
	}

	fmt.Printf("%#v\n", torrents[0])
}
