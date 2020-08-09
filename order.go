package instacart

import "time"

// Item represents a purchased item in teh delivery.
type Item struct {
	ID        string
	ProductID string
	Quantity  int
	Name      string
}

// Delivery is a purhcase at a particular retailer.
type Delivery struct {
	Retailer    string
	DeliveredAt time.Time
	Items       []*Item
}

// Order is the complete transaction.
type Order struct {
	ID         string
	Status     string
	Total      string
	CreatedAt  time.Time
	Deliveries []*Delivery
}

type sortOrders []*Order

func (o sortOrders) Len() int      { return len(o) }
func (o sortOrders) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

type sortOrderByDate struct{ sortOrders }

func (o sortOrderByDate) Less(i, j int) bool {
	return o.sortOrders[i].CreatedAt.Before(o.sortOrders[j].CreatedAt)
}
