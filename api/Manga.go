package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

type Chapter struct {
	Id string `json:"id"`
	OtherId []string `json:"others"`
}

type Volume struct {
	Chapters map[string] Chapter `json:"chapters"`
}

type Manga struct {
	Volumes map[string] Volume `json:"volumes"`
}

func GetManga(id string) Manga {
	ret := Manga{}
	resp, err := http.Get("https://api.mangadex.org/manga/"+id+"/aggregate")
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

	return ret
}