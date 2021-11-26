package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	name        = ".css-1bjwylw"
	price       = ".css-o5uqvq"
	store       = ".css-1kr22w3"
)

func main() {
	crawl()
}

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func crawl() {
	fileName := "top100products.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Could not create file, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	collector := colly.NewCollector(
		colly.AllowedDomains("tokopedia.com", "www.tokopedia.com"),
	)

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	collector.OnHTML(".css-16vw0vn", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText(name),
			e.ChildText(price),
			e.ChildText(store),
		})
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	for i := 1; i <= 10; i++ {
		page := strconv.Itoa(i)
		collector.Visit("https://www.tokopedia.com/p/handphone-tablet/handphone?ob=5&page=" + page)
		time.Sleep(5)
	}
}
