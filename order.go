package instacart

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

// byDate implements sort.Interface for []Order based on the CreatedAt field.
type byDate []Order

func (o byDate) Len() int           { return len(o) }
func (o byDate) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o byDate) Less(i, j int) bool { return o[i].CreatedAt.Before(o[j].CreatedAt) }
