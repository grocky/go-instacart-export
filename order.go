package main

import "time"

type Item struct {
	Id        string
	ProductId string
	Quantity  int
	Name      string
}
type Delivery struct {
	Retailer    string
	DeliveredAt time.Time
	Items       []Item
}
type Order struct {
	Id         string
	Status     string
	Total      string
	CreatedAt  time.Time
	Deliveries []Delivery
}

type SortOrder []Order

func (o SortOrder) Len() int      { return len(o) }
func (o SortOrder) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

type SortOrderByDate struct {
	SortOrder
}

func (o SortOrderByDate) Less(i, j int) bool {
	return o.SortOrder[i].CreatedAt.Before(o.SortOrder[j].CreatedAt)
}
