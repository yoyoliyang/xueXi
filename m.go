package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var jsonRpcPath = "http://192.168.1.52:7912/jsonrpc/0"

// var options = make(map[string]interface{})

func main() {
	// raw := `{"jsonrpc": "2.0", "id": "1", "method": "waitForExists", "params": [{"childOrSibling": [], "childOrSibling": [], "className": "android.view.textview", "mask":16}, 0]}`
	raw := `{"jsonrpc": "2.0", "id": "1", "method": "waitForExists", "params": [{"childOrSibling": [], "childOrSibling": [], "className": "android.view.textview", "mask":16}, 0]}`
	bs := []byte(raw)
	rd := bytes.NewBuffer(bs)

	resp, err := http.Post(jsonRpcPath, "application/json", rd)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))

}
