package main

import (
    "fmt"
    "os"
    "strconv"
    "io/ioutil"
    "io"
	"log"
	"net/http"
	"encoding/json"
    "godexdl/api"
)

type DownloadInfo struct {
    BaseUrl string `json:"baseUrl"`
}

func downloadFile(url string, i int) {
    res, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    filename := strconv.Itoa(i)+".jpg"

    fmt.Printf("        %v\n", filename)
    file, err := os.Create(filename)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    //Write the bytes to the fiel
    _, err = io.Copy(file, res.Body)
    if err != nil {
        panic(err)
    }
}

func getDownloadBaseUrlChapter(id string) string {
    ret := DownloadInfo{}

    resp, err := http.Get("https://api.mangadex.org/at-home/server/"+id)
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

    return ret.BaseUrl
}

func downloadManga(id string) {
    os.Mkdir("download", 0755)
    os.Chdir("download")
    m := api.GetManga(id)
    os.Mkdir("MangaDir", 755)
    
    os.Chdir("MangaDir")
    for volkey, vol := range m.Volumes {
        if volkey == "none" {
            continue
        }
        fmt.Printf("%v\n", volkey)
        os.Mkdir(volkey, 755)
        os.Chdir(volkey)
        for chapkey, chap := range vol.Chapters {
            os.Mkdir(chapkey, 755)
            os.Chdir(chapkey)
            base := getDownloadBaseUrlChapter(chap.Id)
            fmt.Printf("    %v \"%v\" %v\n", chapkey, base, chap.Id)

            c := api.GetChapter(chap.Id)

            for i, file := range c.Attributes.Data {
                path := "/data/"+ c.Attributes.Hash + "/" + file
                downloadFile(base + path, i)
            }
            os.Chdir("..")
        }
        os.Chdir("..")
    }
}

func main() {
	id := "80422e14-b9ad-4fda-970f-de370d5fa4e5"
    downloadManga(id)
}