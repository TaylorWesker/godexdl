package main

import (
	"fmt"
	"godexdl/api"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func downloadFile(url string, i int) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	filename := strconv.Itoa(i) + ".jpg"

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
		if volkey == "none" {
			continue
		}
		// err = os.Mkdir(volkey, 755)
		// if err != nil && err.Error() != "mkdir "+volkey+": file exists" {
		// 	log.Fatalln(err)
		// }
		// os.Chdir(volkey)
		// fmt.Printf("%v\n", volkey)
		for chapkey, chap := range vol.Chapters {
			err = os.Mkdir(chapkey, 755)
			if err != nil && err.Error() != "mkdir "+chapkey+": file exists" {
				log.Fatalln(err)
			}
			os.Chdir(chapkey)
			fmt.Printf("    %v  %v\n", chapkey, chap.Id)
			if last_time.IsZero() || time.Now().Sub(last_time) > time.Second*60 {
				log.Printf("time reset %v\n", nb_download)
				last_time = time.Now()
				nb_download = 0
			}
			if nb_download >= 40 {
				log.Printf("Waiting\n")
				time.Sleep(last_time.Add(time.Second*60).Sub(time.Now()) + time.Microsecond*1000)
			}
			nb_download += 1
			c := api.GetChapter(chap.Id)

			for i, file := range c.Chapter.Data {
				path := "/data/" + c.Chapter.Hash + "/" + file
				downloadFile(c.BaseUrl+path, i+1)
			}
			os.Chdir("..")
		}
		// os.Chdir("..")
	}
}

func main() {
	// id := "80422e14-b9ad-4fda-970f-de370d5fa4e5" // Made in Abyss
	// id := "f8e294c0-7c11-4c66-bdd7-4e25df52bf69" // Blue Period
	id := "881c4da6-87c5-437f-8660-70f3fa4b476a" // Wandering Emanon
	getAllDownloadLinksManga(id)
}
