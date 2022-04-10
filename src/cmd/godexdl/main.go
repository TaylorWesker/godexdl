package main

import (
	"flag"
	"fmt"
	"godexdl/api"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

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

func getAllDownloadLinksManga(id string, min, max float64, outdir string) {
	c := api.GetAllChapter(id)

	sort.SliceStable(c.Data, func(i, j int) bool {
		a, _ := strconv.ParseFloat(c.Data[i].Attributes.Chapter, 64)
		b, _ := strconv.ParseFloat(c.Data[j].Attributes.Chapter, 64)
		return a < b
	})

	last_chapter := c.Data[len(c.Data)-1].Attributes.Chapter

	err := os.Mkdir(outdir, 0755)
	if err != nil && err.Error() != "mkdir download: file exists" {
		log.Println(os.ErrExist)
		log.Fatalf("%v\n", err.Error())
	}
	os.Chdir(outdir)

	m := api.GetManga(id)

	folder := m.Title
	// folder := "test"
	err = os.Mkdir(folder, 0755)
	if err != nil && err.Error() != "mkdir "+folder+": file exists" {
		log.Fatalln(err)
	}
	os.Chdir(folder)

	seen := make(map[string]int)

	for _, c := range c.Data {
		cn, _ := strconv.ParseFloat(c.Attributes.Chapter, 64)
		if c.Attributes.ExternalUrl != "" || cn < min || cn > max {
			continue
		}

		chapkey := c.Attributes.Chapter
		v, ok := seen[c.Attributes.Chapter]
		if !ok {
			seen[c.Attributes.Chapter] = 1
		} else {
			chapkey = fmt.Sprintf("%s-%d", c.Attributes.Chapter, v)
			seen[c.Attributes.Chapter] += 1
		}

		fmt.Printf("%#v\n", c)

		err = os.Mkdir(chapkey, 0755)
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

func main() {
	var id string
	min := 0.
	max := float64(1<<(32-1) - 1)
	var outdir string
	flag.StringVar(&id, "id", "", "id to the mangadex manga")
	flag.Float64Var(&min, "min", min, "minimun chapter number")
	flag.Float64Var(&max, "max", max, "maximum chapter number")
	flag.StringVar(&outdir, "outdir", "download", "folder where the scans will be downloaded")
	flag.Parse()
	getAllDownloadLinksManga(id, min, max, outdir)
}
