package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
)

type Province struct {
	ID   int
	Name string
	URL  string
}

type Kota struct {
	ID   int
	Name string
	URL  string
}

type Kecamatan struct {
	ID   int
	Name string
	URL  string
}

type School struct {
	ID      int
	Name    string
	URL     string
	Email   string
	Website string
}

var tipe string = ""

func GetProvince(url string, types string) {

	if types == "paud" {
		tipe = "paud"
	} else if types == "dikdas" {
		tipe = "dikdas"
	} else if types == "dikmen" {
		tipe = "dikmen"
	}

	start := time.Now()
	defer TimeTrack(start, "Download")

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]Province, 0)

	doc.Find("tbody").Children().Each(func(i int, sel *goquery.Selection) {
		row := new(Province)
		row.ID = i + 1
		row.Name = sel.Find("tr td a").Text()
		row.URL, _ = sel.Find("tr td a").Attr("href")
		rows = append(rows, *row)
		minimalize := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(row.Name, " ", ""), ".", ""))
		os.MkdirAll(tipe+"/data/"+minimalize, 0755)
		GetKota(minimalize, row.Name, row.URL)
	})

	_, err = json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

}

func GetKota(parent string, name string, province string) {

	start := time.Now()
	defer TimeTrack(start, "Download "+name)
	minimalize := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", ""), ".", ""))
	os.MkdirAll(tipe+"/data/"+parent, 0755)
	f, err := os.Create(tipe + "/data/" + parent + "/" + minimalize + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := http.Get(province)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]Kota, 0)

	doc.Find("tbody").Children().Each(func(i int, sel *goquery.Selection) {
		row := new(Kota)
		row.ID = i + 1
		row.Name = sel.Find("tr td a").Text()
		row.URL, _ = sel.Find("tr td a").Attr("href")
		rows = append(rows, *row)

		GetKecamatan(parent, minimalize, row.Name, row.URL)
	})

	bts, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(bts)
	f.Close()

}

func GetKecamatan(bigparent string, parent string, name string, province string) {

	start := time.Now()
	defer TimeTrack(start, "Download "+name)

	minimalize := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", ""), ".", ""))

	os.MkdirAll(tipe+"/data/"+bigparent+"/"+minimalize, 0755)
	f, err := os.Create(tipe + "/data/" + parent + "/" + minimalize + "/" + minimalize + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := http.Get(province)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]Kecamatan, 0)

	doc.Find("tbody").Children().Each(func(i int, sel *goquery.Selection) {
		row := new(Kecamatan)
		row.ID = i + 1
		row.Name = sel.Find("tr td a").Text()
		row.URL, _ = sel.Find("tr td a").Attr("href")
		rows = append(rows, *row)

		GetSchool(bigparent, minimalize, row.Name, row.URL)
	})

	bts, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	f.Write(bts)
	f.Close()

}

func GetSchool(bigparent string, parent string, name string, region string) {

	minimalize := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(name, " ", ""), ".", ""))
	os.MkdirAll(tipe+"/data/"+bigparent+"/"+parent+"/"+minimalize, 0755)
	f, err := os.Create(tipe + "/data/" + bigparent + "/" + parent + "/" + minimalize + "/" + minimalize + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := http.Get(region)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]School, 0)

	doc.Find("tbody").Children().Each(func(i int, sel *goquery.Selection) {

		row := new(School)
		row.ID = i + 1
		row.URL, _ = sel.Find("tr td a").Attr("href")
		sel.Find("tr td").Each(func(j int, s2 *goquery.Selection) {

			if j == 2 {
				email, website := GetAttribute(bigparent, parent, name, row.URL)
				row.Name = s2.Text()
				row.Email = email
				row.Website = website
			}
		})
		rows = append(rows, *row)

	})

	bts, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Create excel file
	xlsx := excelize.NewFile()
	sheet1Name := "Sheet One"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	// set filter cell
	xlsx.SetCellValue(sheet1Name, "A1", "Name")
	xlsx.SetCellValue(sheet1Name, "B1", "Email")
	xlsx.SetCellValue(sheet1Name, "C1", "Website")
	err = xlsx.AutoFilter(sheet1Name, "A1", "C1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	// parse json
	var schools []School
	err = json.Unmarshal(bts, &schools)
	if err != nil {
		log.Fatal(err)
	}

	// insert to cell
	for i, each := range schools {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), each.Name)
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), each.Email)   // Assuming Gender is Email
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), each.Website) // Assuming Age is Website
	}

	// Save file
	err = xlsx.SaveAs(tipe + "/data/" + bigparent + "/" + parent + "/" + minimalize + "/" + minimalize + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}

	f.Write(bts)
	f.Close()

}

func GetAttribute(bigparent string, parent string, name string, region string) (string, string) {

	res, err := http.Get(region)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var email string
	var website string

	doc.Find(".tabs").Children().Each(func(i int, sel *goquery.Selection) {

		if i == 3 {

			sel.Find("tbody tr td").Each(func(j int, s2 *goquery.Selection) {

				// Email
				if j == 11 {
					email = s2.Text()
				}

				// website
				if j == 15 {
					website = s2.Text()
				}

			})

		}

	})

	return email, website

}
