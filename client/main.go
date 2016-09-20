package main

import "net/http"
import "fmt"
import "io/ioutil"

func main() {
	resp, err := http.Get("http://127.0.0.1:8081/api/get-token")
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body")
	}
	fmt.Println(string(body))
}
