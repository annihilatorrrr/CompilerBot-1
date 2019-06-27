package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ArrayContains(arr []string, str string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}
	return false
}

func SendHTTPPostJson(url string, json string) []byte {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(json)),
	)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	res.Body.Close()
	return response
}

func SendHTTPGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	response, err := ioutil.ReadAll(resp.Body)
	return response
}
