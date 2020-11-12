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
	// 子元素获取
	// raw := `{"jsonrpc":"2.0","id":"400bfb24369e9551b529ac8f46a46efb","method":"getText","params":[{"childOrSibling":[],"childOrSiblingSelector":[],"className":"android.widget.TextView","instance":22,"mask":16777232}]}`

	// raw := `{"jsonrpc":"2.0","id":"2753badbb83fecdb214c725173e7e903","method":"count","params":[{"childOrSibling":[],"childOrSiblingSelector":[],"className":"android.widget.TextView","mask":16}]}`
	// 	data := `{"jsonrpc":"2.0","id":"2753badbb83fecdb214c725173e7e903","method":"count","params":[{"childOrSibling":[],"childOrSiblingSelector":[],"className":"android.widget.LinearLayout","mask":16}]}`
	// data := `{"jsonrpc":"2.0","id":"441f93cd4abf6cbb4e8e562ac1c0f7a3","method":"getText","params":[{"childOrSibling":[],"childOrSiblingSelector":[],"className":"android.view.View","instance":20,"mask":16777232}]}`
	data := `{"jsonrpc":"2.0","id":"441f93cd4abf6cbb4e8e562ac1c0f7a3","method":"getText","params":[{"childOrSibling":[],"childOrSiblingSelector":[],"className":"android.view.View","instance":20,"mask":16777232}]}`

	bs := []byte(data)
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
