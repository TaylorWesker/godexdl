package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

type ChapterAttribute struct {
	Hash string `json:"hash"`
	Data []string `json:"data"`
	DataSaver []string `json:"dataSaver"`
}

type ChapterData struct {
	Id string `json:"id"`
	Attributes ChapterAttribute `json:"attributes"`
}

type ChapterResponse struct {
	Data ChapterData `json:"data"`
}

func GetChapter(id string) ChapterData {
	ret := ChapterResponse{}
	resp, err := http.Get("https://api.mangadex.org/chapter/"+id)
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

	return ret.Data
}