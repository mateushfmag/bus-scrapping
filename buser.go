package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func parseNumber(text string ) float32{
	regex := regexp.MustCompile(`R\$\s+`)
	value := regex.ReplaceAllString(text, "")
	regex = regexp.MustCompile(`,`)
	value = regex.ReplaceAllString(value,".")
	result, _ := strconv.ParseFloat(value, 32)
	return float32(result)
}

func trim(text string) string{
	return strings.Trim(text,"\n ")
}

func saveData(buserOptions []BusSummary){
	// sort.SliceStable(buserOptions, func(i,j int) bool {
	// 	return buserOptions[i].Price < buserOptions[j].Price
	// })
	
	file, _ := json.MarshalIndent(buserOptions,"","    ")
	os.Mkdir("./results", 0700)
	os.WriteFile("./results/buser.json",file, 0700)
}

func getBusOptions(c *colly.Collector){
	var buserOptions []BusSummary
	saveData(buserOptions)
	c.OnHTML(".g-header", func(h *colly.HTMLElement) {
		summary := BusSummary{}
		_, hasDisabled := h.DOM.Find(".g-reservar button").Attr("disabled")

		summary.Available = hasDisabled == false
		destination := h.DOM.Find(".ir-data.is-destino")
		summary.Destination.Date = trim(destination.Find(".ird-dia").Text()) 
		summary.Destination.Hour = trim(destination.Find(".ird-hora").Text())


		regex := regexp.MustCompile(`.+:`)
		summary.Destination.Location = trim(regex.ReplaceAllString(h.DOM.Find(".ir-endereco.is-destino").Text(),""))

		origin := h.DOM.Find(".ir-data.is-origem")
		summary.Origin.Date = trim(origin.Find(".ird-dia").Text())
		summary.Origin.Hour = trim(origin.Find(".ird-hora").Text())
		summary.Origin.Location = trim(regex.ReplaceAllString(h.DOM.Find(".ir-endereco.is-origem").Text(),""))

		summary.Price = parseNumber(trim(h.DOM.Find(".p-preco").First().Text()))
		summary.Seat = trim(h.DOM.Find(".p-assento").Text())

		buserOptions = append(buserOptions, summary)

		saveData(buserOptions)
	})
}

func search(c *colly.Collector, origin string, destination string, departure string, arrival string){
	query := fmt.Sprintf("%s/%s?ida=%s&volta=%s", origin,destination, departure, arrival)
	c.Visit(fmt.Sprintf("https://www.buser.com.br/onibus/%s", query))
}

func mapper(origin string, destination string) (string , string) {
	if origin == "MG-BH" {
		origin = "belo-horizonte-mg"
	}

	if destination == "RJ-RJ" {
		destination = "rio-de-janeiro-rj"
	}

	return origin, destination
}

func GetBuser(origin string, destination string, departure string, arrival string){
	c := colly.NewCollector()
	getBusOptions(c)
	origin, destination = mapper(origin,destination)
	search(c, origin, destination, departure, arrival)
}