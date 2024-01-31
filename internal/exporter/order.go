package exporter

import "time"

// Item represents a purchased item in the delivery.
type Item struct {
	ID        string
	ProductID string
	Quantity  int
	Name      string
}

// Delivery is a purchase at a particular retailer.
type Delivery struct {
	Retailer    string
	DeliveredAt time.Time
	Items       []Item
}

// Order is the complete transaction.
type Order struct {
	ID         string
	Status     string
	Total      string
	CreatedAt  time.Time
	Deliveries []Delivery
}

// ByDate implements sort.Interface for []Order based on the CreatedAt field.
type ByDate []*Order

func (o ByDate) Len() int           { return len(o) }
func (o ByDate) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ByDate) Less(i, j int) bool { return o[i].CreatedAt.After(o[j].CreatedAt) }
