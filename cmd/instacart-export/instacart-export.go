package main

import (
	"encoding/csv"
	"fmt"
	"github.com/grocky/go-instacart-export/instacart"
	"github.com/grocky/go-instacart-export/internal/exporter"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/getopt/v2"
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
	order := exporter.NewOrderService(client)

	log.Print("Fetching orders...")
	orders, err := order.GetOrderPages(startPage, endPage)
	if err != nil {
		log.Printf("some page requests failed...: %s", err)
		log.Print(err)
	}

	if len(orders) == 0 {
		log.Print("nothing to write to the file...exiting")
		os.Exit(3)
	}

	log.Print("Processing orders")
	data := convertOrderToCSV(orders)

	log.Print("Writing orders to a CSV")
	err = writeToCSV(data)
	if err != nil {
		fmt.Print(err)
	}

	log.Print("Done!")
}

func convertOrderToCSV(orders []*exporter.Order) [][]string {
	data := [][]string{{
		"id",
		"status",
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

func writeToCSV(data [][]string) error {
	now := time.Now()
	file, err := os.Create("data/instacart_orders_" + now.Format("01-02-2006_03-04-05") + ".csv")
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failure writing data: %w", err)
		}
	}

	return nil
}
