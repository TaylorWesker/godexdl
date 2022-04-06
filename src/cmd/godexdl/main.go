package main

import (
	"fmt"
	"godexdl/api"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

func downloadFile(url string, i int) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	filename := strconv.Itoa(i) + ".jpg"

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

func getAllDownloadLinksMangaV2(id string, min, max int) {
	c := api.GetAllChapter(id)

	sort.SliceStable(c.Data, func(i, j int) bool {
		a, _ := strconv.Atoi(c.Data[i].Attributes.Chapter)
		b, _ := strconv.Atoi(c.Data[j].Attributes.Chapter)
		return a < b
	})

	last_chapter := c.Data[len(c.Data)-1].Attributes.Chapter

	err := os.Mkdir("download", 0755)
	if err != nil && err.Error() != "mkdir download: file exists" {
		log.Println(os.ErrExist)
		log.Fatalf("%v\n", err.Error())
	}
	os.Chdir("download")

	m := api.GetManga(id)

	folder := m.Title
	// folder := "test"
	err = os.Mkdir(folder, 755)
	if err != nil && err.Error() != "mkdir "+folder+": file exists" {
		log.Fatalln(err)
	}
	os.Chdir(folder)

	for _, c := range c.Data {
		cn, _ := strconv.Atoi(c.Attributes.Chapter)
		if c.Attributes.ExternalUrl != "" || cn < min || cn > max {
			continue
		}
		fmt.Printf("%#v\n", c)
		chapkey := c.Attributes.Chapter

		err = os.Mkdir(chapkey, 755)
		if err != nil && err.Error() != "mkdir "+chapkey+": file exists" {
			log.Fatalln(err)
		}
		os.Chdir(chapkey)

		c := api.GetChapter(c.Id)

		page_bar := progressbar.Default(int64(len(c.Chapter.Data)), fmt.Sprintf("downloading chapter number %s/%s", chapkey, last_chapter))
		for i, file := range c.Chapter.Data {
			path := "/data/" + c.Chapter.Hash + "/" + file
			downloadFile(c.BaseUrl+path, i+1)
			page_bar.Add(1)
		}
		os.Chdir("..")
	}
}

func getAllDownloadLinksManga(id string) {

	err := os.Mkdir("download", 0755)
	if err != nil && err.Error() != "mkdir download: file exists" {
		log.Println(os.ErrExist)
		log.Fatalf("%v\n", err.Error())
	}
	os.Chdir("download")

	var last_time time.Time

	nb_download := 0

	m := api.GetManga(id)

	folder := m.Title
	err = os.Mkdir(folder, 755)
	if err != nil && err.Error() != "mkdir "+folder+": file exists" {
		log.Fatalln(err)
	}
	os.Chdir(folder)
	for volkey, vol := range m.Volumes {
		// if volkey == "none" {
		// 	continue
		// }

		chapter_bar := progressbar.Default(int64(len(vol.Chapters)), fmt.Sprint("downloading volume number ", volkey))
		for chapkey, chap := range vol.Chapters {
			err = os.Mkdir(chapkey, 755)
			if err != nil && err.Error() != "mkdir "+chapkey+": file exists" {
				log.Fatalln(err)
			}
			os.Chdir(chapkey)
			if last_time.IsZero() || time.Now().Sub(last_time) > time.Second*60 {
				last_time = time.Now()
				nb_download = 0
			}
			if nb_download >= 40 {
				time.Sleep(last_time.Add(time.Second*60).Sub(time.Now()) + time.Microsecond*1000)
			}
			nb_download += 1
			c := api.GetChapter(chap.Id)

			for i, file := range c.Chapter.Data {
				path := "/data/" + c.Chapter.Hash + "/" + file
				downloadFile(c.BaseUrl+path, i+1)
			}
			os.Chdir("..")
			chapter_bar.Add(1)
		}
	}
}

func main() {
	// id := "80422e14-b9ad-4fda-970f-de370d5fa4e5" // Made in Abyss
	// id := "f8e294c0-7c11-4c66-bdd7-4e25df52bf69" // Blue Period
	// id := "881c4da6-87c5-437f-8660-70f3fa4b476a" // Wandering Emanon
	// https://mangadex.org/title/a77742b1-befd-49a4-bff5-1ad4e6b0ef7b/chainsaw-man
	id := "162146eb-672a-4a05-b3b2-0c6303f9614e"
	id = "a77742b1-befd-49a4-bff5-1ad4e6b0ef7b"
	id = "6fcfaa0e-6023-403e-97f9-5301dd3c258c"
	id = "f5e3baad-3cd4-427c-a2ec-ad7d776b370d"
	getAllDownloadLinksMangaV2(id, 0, 1<<(32-1)-1)
}
