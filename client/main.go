package main

import "net/http"
import "net/url"
import "fmt"
import "io/ioutil"

func main() {
	//client := &http.Client{}
	//req, err := http.NewRequest("POST", "http://127.0.0.1:8081/api/get-token/", nil)

	usernamePwd := url.Values{}
	usernamePwd.Set("username", "suraj")
	usernamePwd.Set("password", "suraj")

	//	if err != nil {
	//		fmt.Println("Unable to form a POST")
	//	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//resp, err := client.Do(req)
	resp, err := http.PostForm("http://127.0.0.1:8081/api/get-token/", usernamePwd)
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
