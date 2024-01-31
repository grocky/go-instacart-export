package exporter

import (
	"errors"
	"fmt"
	"github.com/grocky/go-instacart-export/instacart"
	"log"
	"sort"
	"time"
)

const timeFormat = "Jan 2, 2006,  3:04 PM"

// OrderService implements business logic for exporting orders from Instacart.
type OrderService struct {
	instacartClient *instacart.Client
}

// NewOrderService creates a new OrderService.
func NewOrderService(client *instacart.Client) *OrderService {
	return &OrderService{instacartClient: client}
}

// GetOrderPages retrieves pages of orders starting with start and ending with end, inclusive.
// Results are returned in reverse chronological order.
func (o *OrderService) GetOrderPages(start, end int) ([]*Order, error) {
	var orders []*Order
	var nextPage = new(int)
	*nextPage = start

	if start > end {
		return nil, errors.New("start must be less than or equal to end")
	}

	for nextPage != nil && *nextPage <= end {
		log.Printf("fetching page: %d", *nextPage)
		resp, err := o.instacartClient.FetchPage(*nextPage)
		if err != nil {
			return orders, fmt.Errorf("failed to fetch page %d: %w", *nextPage, err)
		}
		o, err := extractOrdersFromResponse(resp)
		if err != nil {
			return orders, err
		}
		orders = append(orders, o...)
		nextPage = resp.Meta.Pagination.NextPage
	}

	sort.Sort(ByDate(orders))
	return orders, nil
}

func extractOrdersFromResponse(orderResp *instacart.OrdersResponse) ([]*Order, error) {
	var orders []*Order
	for _, instacartOrder := range orderResp.Orders {
		order, err := convertInstacartApiOrder(instacartOrder)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func convertInstacartApiOrder(instacartOrder instacart.Order) (*Order, error) {
	order := &Order{}
	order.ID = instacartOrder.ID
	order.Status = instacartOrder.Status
	order.Total = instacartOrder.Total

	createdAt, err := time.Parse(timeFormat, instacartOrder.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to parse order time %s: %w", instacartOrder.CreatedAt, err)
	}
	order.CreatedAt = createdAt

	var deliveries []Delivery
	for _, d := range instacartOrder.OrderDeliveries {
		delivery := Delivery{}
		delivery.Retailer = d.Retailer.Name

		if d.DeliveredAt != "" {
			deliveredAt, err := time.Parse(timeFormat, d.DeliveredAt)
			if err != nil {
				return nil, fmt.Errorf("unable to parse delivery time %s, %w", d.DeliveredAt, err)
			}
			delivery.DeliveredAt = deliveredAt
		}

		var items []Item
		for _, i := range d.OrderItems {
			item := Item{}
			item.ID = i.Item.ID
			item.ProductID = i.Item.ProductID
			item.Quantity = int(i.Qty)
			item.Name = i.Item.Name

			items = append(items, item)
		}

		delivery.Items = items
		deliveries = append(deliveries, delivery)
	}

	order.Deliveries = deliveries

	return order, nil
}
