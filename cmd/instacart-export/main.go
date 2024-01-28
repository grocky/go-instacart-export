package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/getopt/v2"

	instacart "github.com/grocky/go-instacart-export"
)

var (
	help         = false
	startPage    = 1
	endPage      = 10
	sessionToken = ""
)

func init() {
	getopt.FlagLong(&startPage, "start", 0, "The first page of order results to request")
	getopt.FlagLong(&endPage, "end", 0, "The last page of order results to request")
	getopt.FlagLong(&help, "help", 'h', "Help!")
}

func main() {

	getopt.Parse()
	if help {
		getopt.Usage()
		return
	}

	sessionToken = os.Getenv("INSTACART_SESSION_TOKEN")
	if sessionToken == "" {
		log.Println("Session token missing. Please provide the INSTACART_SESSION_TOKEN environment variable")
		getopt.Usage()
		return
	}

	client := instacart.NewClient(sessionToken)

	log.Print("Fetching orders...")
	orders := client.FetchOrders(startPage, endPage)
	data := extractOrdersData(orders)
	writeToCSV(data)

	log.Print("Done!")
}

func extractOrdersData(orders []*instacart.Order) [][]string {
	log.Print("Processing orders")
	data := [][]string{{
		"id",
		"satus",
		"total",
		"createdAt",
		"retailers",
		"numItems",
	}}
	for _, o := range orders {

		var retailers []string
		numItems := 0

		for _, d := range o.Deliveries {
			retailers = append(retailers, d.Retailer)
			numItems += len(d.Items)
		}

		order := []string{
			o.ID,
			o.Status,
			o.Total,
			o.CreatedAt.Format("2006-01-02"),
			strings.Join(retailers, "|"),
			strconv.Itoa(numItems),
		}
		data = append(data, order)
	}

	return data
}

func writeToCSV(data [][]string) {
	log.Print("Writing orders to a CSV")

	now := time.Now()
	file, err := os.Create("data/instacart_orders_" + now.Format("01-02-2006_03-04-05") + ".csv")
	if err != nil {
		log.Fatal("Unable to create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Fatal("Error writing data", err)
		}
	}
}
