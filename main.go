package main

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "Jan 2, 2006,  3:04 PM"

func extractOrders(orderResp OrdersResponse) []Order {
	var orders []Order
	for _, o := range orderResp.Orders {
		order := Order{}
		order.Id = o.ID
		order.Status = o.Status
		order.Total = o.Total

		createdAt, err := time.Parse(timeFormat, o.CreatedAt)
		if err != nil {
			log.Fatalf("Unable to parse order time: %s | %v", o.CreatedAt, err)
		}
		order.CreatedAt = createdAt

		var deliveries []Delivery
		for _, d := range o.OrderDeliveries {
			delivery := Delivery{}
			delivery.Retailer = d.Retailer.Name

			if d.DeliveredAt != "" {
				deliveredAt, err := time.Parse(timeFormat, d.DeliveredAt)
				if err != nil {
					log.Fatalf("Unable to parse delivery time: %s | %v", d.DeliveredAt, err)
				}
				delivery.DeliveredAt = deliveredAt
			}

			var items []Item
			for _, i := range d.OrderItems {
				item := Item{}
				item.Id = i.Item.ID
				item.ProductId = i.Item.ProductID
				item.Quantity = int(i.Qty)
				item.Name = i.Item.Name

				items = append(items, item)
			}

			delivery.Items = items
			deliveries = append(deliveries, delivery)
		}

		order.Deliveries = deliveries
		orders = append(orders, order)
	}

	return orders
}

func fetchOrders(client Client) []Order {
	var orders []Order
	var resp OrdersResponse
	var nextPage *int
	nextPage = new(int)

	*nextPage = 1

	for nextPage != nil {
		log.Printf("Fetching page: %d", *nextPage)
		resp = client.getPage(*nextPage)
		orders = append(orders, extractOrders(resp)...)
		nextPage = resp.Meta.Pagination.NextPage
	}

	return orders
}

func main() {
	authCookie := os.Getenv("AUTH_COOKIE")
	csrf := os.Getenv("AUTH_CSRF")

	client := Client{
		userCookie: authCookie,
		csrfToken:  csrf,
	}

	log.Print("Fetching orders...")
	orders := fetchOrders(client)
	sort.Sort(SortOrderByDate{orders})

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
			o.Id,
			o.Status,
			o.Total,
			o.CreatedAt.Format("2006-01-02"),
			strings.Join(retailers, "|"),
			strconv.Itoa(numItems),
		}
		data = append(data, order)
	}

	log.Print("Writing orders to a CSV")

	now := time.Now()
	file, err := os.Create("instacart_orders_" + now.Format("01-02-2006_03:04:05") + ".csv")
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
	log.Print("Done!")
}
