package main

import (
	"fmt"
	"log"

	"github.com/kolo/xmlrpc"
)

// type Args struct {
// 	A int
// }

func main() {
	client, err := xmlrpc.NewClient("http://192.168.2.7:5001/RPC2", nil)
	// client, err := xmlrpc.NewClient("http://127.0.0.1:5001/RPC2", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
	var version string
	err = client.Call("system.api_version", nil, &version)
	// err = client.Call("service.time", nil, &version)
	if err != nil {
		log.Fatal("error: ", err)
	}
	fmt.Printf("Version: %s\n", version)

	type Result struct {
		Name string `xmlrpc:"string"`
	}

	var output [][]string
	// 	if err := client.Call("d.multicall", []interface{}{"main", "d.get_name=", "d.bytes_done="}, &output); err != nil {
	if err := client.Call("d.multicall", []interface{}{"main", "d.get_name="}, &output); err != nil {
		fmt.Println("service.sum call error: ", err)
	}
	fmt.Println(output)

	// var torrents []interface{}
	// args := make([]interface{}, 3)
	// args[0] = "main"
	// args[1] = "d.bytes_done="
	// args[2] = "d.size_bytes="
	// fmt.Println(args)
	// fmt.Printf("%#v\\n", args)
	// err = client.Call("d.multicall", args, &torrents)
	// err = client.Call("sum", []interface{}{2, 3}, &torrents)
	// if err != nil {
	// 	log.Fatal("multicall: ", err)
	// }
	// fmt.Printf("torrents: %#v\\\\\\\\n", torrents)

	// var sum int
	// if err := client.Call("service.sum", []interface{}{2, 3}, &sum); err != nil {
	// 	fmt.Printf("service.sum call error: %v", err)
	// }
	// fmt.Println(sum)

	// works locally:
	// type Result struct {
	// 	Version string `xmlrpc:"version"`
	// }

	// var result Result
	// client.Call("service.version", nil, &result)
	// fmt.Println(result)
	// fmt.Printf("Version: %s\\n", result) // Version: 4.2.7+

	// type Result struct {
	// 	Version string `xmlrpc:"bytes_done"`
	// 	Mark    string `xmlrpc:"bytes_done="`
	// }
	// var output [][]interface{}
	// var output []string
	// var output = make([]interface{}, 3)
	// if err := client.Call("d.multicall", []interface{}{"main", "d.bytes_done="}, &output); err != nil {
	// 	fmt.Printf("service.sum call error: %v\\n", err)
	// }
	// fmt.Println(output)

	// var output [][]interface{}
	// if err := client.Call("service.array", []interface{}{"main", "d.bytes_done="}, &output); err != nil {
	// 	fmt.Printf("service.sum call error: %v\\n", err)
	// }
	// fmt.Println("output: ", output)

	// for _, torrent := range torrents {
	// 	fmt.Println(torrent)
	// }

	// fmt.Printf("torrents: %#v\\n", len(torrents))
}
