package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Chapter struct {
	Id      string   `json:"id"`
	OtherId []string `json:"others"`
}

type Volume struct {
	Chapters map[string]Chapter `json:"chapters"`
}

/*really dirty way of geting back the manga title*/

type Manga4 struct {
	Title map[string]string `json:"title"`
}

type Manga3 struct {
	Attributes Manga4 `json:"attributes"`
}

type Manga2 struct {
	Data Manga3 `json:"data"`
}

/*******************************************/

type Manga struct {
	Result  string        `json:"result"`
	Errors  []ErrorReport `json:"errors"`
	Title   string
	Volumes map[string]Volume `json:"volumes"`
}

func GetManga(id string) Manga {
	ret := Manga{}
	resp, err := http.Get("https://api.mangadex.org/manga/" + id + "/aggregate?translatedLanguage[]=en")
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

	extra := Manga2{}

	resp, err = http.Get("https://api.mangadex.org/manga/" + id)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	json.Unmarshal(body, &extra)

	if ret.Result != "ok" {
		log.Fatalf("Api response is %v %v\n", ret.Errors, ret.Result)
	}

	ret.Title = extra.Data.Attributes.Title["en"]

	return ret
}
