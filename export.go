package instacart

import (
	"log"
	"time"
)

func extractOrders(orderResp OrdersResponse) []Order {
	var orders []Order
	for _, o := range orderResp.Orders {
		order := Order{}
		order.ID = o.ID
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
		orders = append(orders, order)
	}

	return orders
}
