package main

import (
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/fvbommel/sortorder"
	"github.com/go-pdf/fpdf"
	"github.com/schollz/progressbar/v3"
)

func addImagePage(pdf *fpdf.Fpdf, image string) error {
	var opt fpdf.ImageOptions

	format := "jpg"

	file, err := os.Open(image)
	if err != nil {
		return err
	}
	config, err := jpeg.DecodeConfig(file)
	if err != nil {
		format = "png"
		file.Seek(0, 0)
		config, err = png.DecodeConfig(file)
		if err != nil {
			return err
		}
	}
	aspect_ratio := float64(config.Height) / float64(config.Width)
	file.Seek(0, 0)
	wd, ht := pdf.GetPageSize()
	defaultFormat := fpdf.SizeType{Wd: wd, Ht: ht}
	otherFormat := fpdf.SizeType{Wd: defaultFormat.Ht * (1 / float64(aspect_ratio)), Ht: defaultFormat.Ht}
	pdf.AddPageFormat("P", otherFormat)
	opt.ImageType = format
	pdf.RegisterImageOptionsReader(image, opt, file)
	pdf.ImageOptions(image, 0, 0, otherFormat.Wd, otherFormat.Ht, false, opt, 0, "")

	return nil
}

func generatePdf(basePath string, dirPath string) {

	pdf := fpdf.New("P", "mm", "A4", "")
	fullPath := basePath + "/" + dirPath
	files, _ := ioutil.ReadDir(fullPath)
	var links []string
	for _, file := range files {
		links = append(links, file.Name())
	}

	sort.Slice(links, func(i, j int) bool {
		return sortorder.NaturalLess(links[i], links[j])
	})

	for _, l := range links {
		link := fullPath + "/" + l
		err := addImagePage(pdf, link)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err := pdf.OutputFileAndClose("chapter_" + dirPath + ".pdf")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	dir := os.Args[1]

	dirs, _ := ioutil.ReadDir(dir)

	bar := progressbar.Default(int64(len(dirs)))

	for _, d := range dirs {
		generatePdf(dir, d.Name())
		bar.Add(1)
	}

}
