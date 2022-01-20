package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type ErrorReport struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type ChapterData struct {
	Hash string   `json:"hash"`
	Data []string `json:"data"`
}

type ChapterResponse struct {
	Result  string        `json:"result"`
	Errors  []ErrorReport `json:"errors"`
	BaseUrl string        `json:"baseUrl"`
	Chapter ChapterData   `json:"chapter"`
}

func GetChapter(id string) ChapterResponse {
	ret := ChapterResponse{}
	resp, err := http.Get("https://api.mangadex.org/at-home/server/" + id)

	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	json.Unmarshal(body, &ret)

	if ret.Result != "ok" {
		log.Fatalf("Api response is %v %v\n", ret.Errors, ret.Result)
	}

	return ret
}
