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

func main() {
	defer trackTime(time.Now(), "xmlRpc")
	client, err := xmlrpc.NewClient("http://192.168.2.7:5001/RPC2", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
	var version string
	err = client.Call("system.api_version", nil, &version)
	if err != nil {
		log.Fatal("error: ", err)
	}
	fmt.Printf("Version: %s\n", version)

	var output [][]interface{}
	if err := client.Call("d.multicall", []interface{}{"main", "d.bytes_done=", "d.connection_current=", "d.creation_date=", "d.get_down_rate=", "d.get_down_total=", "d.size_bytes=", "d.size_files=", "d.state=", "d.load_date=", "d.name=", "d.ratio=", "d.get_up_rate=", "d.get_up_total=", "d.hash="}, &output); err != nil {

		fmt.Println("service.sum call error: ", err)
	}
	fmt.Println(output)
}
