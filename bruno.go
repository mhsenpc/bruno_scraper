package main

import "fmt"
import "github.com/gocolly/colly/v2"
import (
    "encoding/csv"
    "log"
    "os"
)

func main(){
    fmt.Println("Collecting links in sitemap");
    //loop through links in https://www.bruno.com/sitemap

    records := [][]string{
            {"Name", "Tel", "Location"},
        }

   	// Array containing all the known URLs in a sitemap
   	knownUrls := []string{}

   	// Create a Collector specifically for Shopify
   	c := colly.NewCollector(colly.AllowedDomains("www.bruno.com"))

   	// Create a callback on the XPath query searching for the URLs
   	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
   	    knownUrls = append(knownUrls, e.Text)
   		c.Visit(e.Text)
   	})

    c.OnHTML("body", func(e *colly.HTMLElement) {
        Name := e.DOM.Find(".DealerDetailsPageContent h1").Text();
        Tel,_ := e.DOM.Find("#CtaPhone").Attr("href");
        Location,_ := e.DOM.Find("#CtaMap").Attr("href");
        if(Name != ""){
            new_row := []string{Name,Tel,Location}
            records = append(records,new_row);
            fmt.Println( "Collected data: " + Name );
        }
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

   	// Start the collector
   	c.Visit("https://www.bruno.com/sitemap")


    f, err := os.Create("list.csv")
    defer f.Close()

    if err != nil {

        log.Fatalln("failed to open file", err)
    }

    w := csv.NewWriter(f)
    err = w.WriteAll(records) // calls Flush internally

    if err != nil {
        log.Fatal(err)
    }

   	fmt.Println("Collected", len(knownUrls), "URLs")
}