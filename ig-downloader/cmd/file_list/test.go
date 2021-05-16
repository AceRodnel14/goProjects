package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	req, err := http.NewRequest("GET", "https://www.instagram.com/p/CMcLaT8FezH/", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		//log.Info(bodyString)
		regexMatchString := " {\"@context\"(.*)"

		regex, err := regexp.Compile(regexMatchString)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(regex.FindString(bodyString))
	}
}
