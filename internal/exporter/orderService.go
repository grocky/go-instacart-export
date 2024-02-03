package exporter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/grocky/go-instacart-export/instacart"
)

const timeFormat = "Jan 2, 2006,  3:04 PM"

// OrderService implements business logic for exporting orders from Instacart.
type OrderService struct {
	instacartClient *instacart.Client
	numWorkers      int
}

// NewOrderService creates a new OrderService.
func NewOrderService(client *instacart.Client) *OrderService {
	return &OrderService{
		instacartClient: client,
		numWorkers:      10,
	}
}

// GetOrderPages retrieves pages of orders starting with start and ending with end, inclusive.
// Results are returned in reverse chronological order.
func (o *OrderService) GetOrderPages(start, end int) ([]*Order, error) {

	if start > end {
		return nil, errors.New("start must be less than or equal to end")
	}

	var wg sync.WaitGroup

	numTasks := end - start + 1
	tasks := make(chan task, numTasks)
	results := make(chan task, numTasks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create the workers
	for i := 0; i < o.numWorkers; i++ {
		wg.Add(1)
		go worker(workerContext{
			ctx:     ctx,
			client:  o.instacartClient,
			tasks:   tasks,
			results: results,
			cancel:  cancel,
			wg:      &wg,
		})
	}

	// generate the tasks
	for i := start; i <= end; i++ {
		tasks <- task{page: i}
	}

	close(tasks)

	wg.Wait()

	close(results)

	// collect the results
	var orders []*Order
	for r := range results {
		if r.orders != nil {
			orders = append(orders, r.orders...)
		}
	}

	sort.Sort(ByDate(orders))

	return orders, nil
}

type workerContext struct {
	client  *instacart.Client
	ctx     context.Context
	cancel  context.CancelFunc
	tasks   <-chan task
	results chan<- task
	wg      *sync.WaitGroup
}

func worker(wctx workerContext) {
	defer wctx.wg.Done()

	var err error
	var resp *instacart.OrdersResponse
	var orders []*Order

	for t := range wctx.tasks {
		// detect cancellation
		select {
		case <-wctx.ctx.Done():
			return
		default:
		}

		log.Printf("fetching page: %d", t.page)
		resp, err = wctx.client.FetchPage(t.page)
		if err != nil {
			log.Printf("failed to fetch page %d: %s", t.page, err)
			wctx.cancel()
			return
		}

		if len(resp.Orders) == 0 {
			log.Printf("no items to process on page: %d", t.page)
			return
		}

		if resp.Meta.Pagination.NextPage == nil {
			log.Printf("no more pages left to process: current page: %d", t.page)
			wctx.cancel()
		}

		orders, err = extractOrdersFromResponse(resp)
		if err != nil {
			log.Printf("unable to convert order: %s", err)
		}

		t.orders = orders

		wctx.results <- t
		log.Printf("completed processing page: %d", t.page)
	}
}

type task struct {
	page   int
	orders []*Order
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
