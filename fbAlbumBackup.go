package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	fb "github.com/huandu/facebook"
)

func getFBParams(token string) fb.Params {
	return fb.Params{
		"access_token": token}
}

func getPhotoName(photoSource string, height int64, width int64) (photoName string, err error) {
	url, err := url.Parse(photoSource)
	photoName = ""
	if err == nil {
		path := url.Path
		splitted := strings.Split(path, "/")
		photoName = splitted[len(splitted)-1]
		photoName = fmt.Sprintf("%dx%d-%s", width, height, photoName)
	}
	return
}

func getPhoto(source string, photoName string) (err error) {
	resp, err := http.Get(source)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		ioutil.WriteFile(photoName, body, 0644)
	}
	return
}

func worker(id int, jobs <-chan fb.Result, result chan<- error) {
	for photo := range jobs {
		source := photo["source"].(string)
		width := photo["width"].(json.Number)
		height := photo["height"].(json.Number)
		photoWidth, _ := width.Int64()
		photoHeight, _ := height.Int64()
		photoName, err := getPhotoName(source, photoWidth, photoHeight)
		if err != nil {
			result <- err
		} else {
			result <- getPhoto(source, photoName)
		}
	}
}

func main() {
	token := flag.String("token", "", "FB access token")
	workers := flag.Int("worker", 5, "number of workers")
	var items []fb.Result
	var albums []fb.Result
	jobs := make(chan fb.Result)
	jobErr := make(chan error)
	flag.Parse()
	if *token == "" {
		log.Fatalln("-token flag is required")
	}
	fbParams := getFBParams(*token)
	res, err := fb.Get("/me/albums", fbParams)
	for w := 1; w <= *workers; w++ {
		go worker(w, jobs, jobErr)
	}
	if err == nil {
		res.DecodeField("data", &items)
		for _, item := range items {
			var albumURL = fmt.Sprintf("/%s/photos", item["id"])
			res, err := fb.Get(albumURL, fbParams)
			res.DecodeField("data", &albums)
			if err == nil {
				for _, photo := range albums {
					photoURL := fmt.Sprintf("/%s?fields=images", photo["id"])
					res, _ := fb.Get(photoURL, fbParams)
					res.DecodeField("images", &albums)
					for _, photo := range albums {
						jobs <- photo
						err := <-jobErr
						if err != nil {
							log.Println(err)
						}
					}
				}
			} else {
				log.Fatalln(err)
			}
		}
	} else {
		log.Fatalln(err)
	}
}
