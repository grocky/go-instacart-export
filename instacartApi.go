package instacart

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sort"
	"strconv"
	"time"
)

// Client is the HTTP client for the Instacart orders API.
type Client struct {
	SessionToken string
	httpClient   *http.Client
}

// NewClient constructs a Client.
func NewClient(sessionToken string) *Client {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: jar,
	}
	return &Client{
		SessionToken: sessionToken,
		httpClient:   httpClient,
	}
}

const timeFormat = "Jan 2, 2006,  3:04 PM"

// FetchOrders retrieves all orders sorted by date created, descending.
func (c *Client) FetchOrders(start, end int) []*Order {
	var orders []*Order
	var resp OrdersResponse
	var nextPage = new(int)
	*nextPage = start

	for nextPage != nil && *nextPage <= end {
		log.Printf("Fetching page: %d", *nextPage)
		resp = c.getPage(*nextPage)
		orders = append(orders, extractOrders(resp)...)
		nextPage = resp.Meta.Pagination.NextPage
	}

	sort.Sort(sortOrderByDate{orders})
	return orders
}

func (c *Client) getPage(page int) OrdersResponse {

	url := "https://www.instacart.com/v3/orders?page=" + strconv.Itoa(page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Host", "www.instacart.com")
	req.Header.Set("User-Agent", "Instacart Orders To CSV Client")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Identifier", "web")

	sessionCookie := &http.Cookie{
		Name:  "_instacart_session_id",
		Value: c.SessionToken,
	}
	req.AddCookie(sessionCookie)

	// sets the ":authority" pseudo-header field for HTTP/2
	//req.Host = "www.instacart.com"

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("There was a problem with the request to instacart...")
		log.Println(resp.Status)
		log.Fatalf("%+v\n", resp)
	}

	defer resp.Body.Close()

	var ordersResp OrdersResponse

	if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
		log.Fatal(err)
	}

	return ordersResp
}

func extractOrders(orderResp OrdersResponse) []*Order {
	var orders []*Order
	for _, o := range orderResp.Orders {
		order := &Order{}
		order.ID = o.ID
		order.Status = o.Status
		order.Total = o.Total

		createdAt, err := time.Parse(timeFormat, o.CreatedAt)
		if err != nil {
			log.Fatalf("Unable to parse order time: %s | %v", o.CreatedAt, err)
		}
		order.CreatedAt = createdAt

		var deliveries []*Delivery
		for _, d := range o.OrderDeliveries {
			delivery := &Delivery{}
			delivery.Retailer = d.Retailer.Name

			if d.DeliveredAt != "" {
				deliveredAt, err := time.Parse(timeFormat, d.DeliveredAt)
				if err != nil {
					log.Fatalf("Unable to parse delivery time: %s | %v", d.DeliveredAt, err)
				}
				delivery.DeliveredAt = deliveredAt
			}

			var items []*Item
			for _, i := range d.OrderItems {
				item := &Item{}
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

// OrdersResponse is the response from the orders API
// auto-generated from: https://mholt.github.io/json-to-go/
//   - Updated Actions to be map[string]struct
//   - Updated .orders.order_deliveries.order_items.qty to be float
//   - Updated .orders.order_deliveries.order_items.item.qty_attributes.increment  to be float
//   - Updated .orders.order_deliveries.order_items.item.qty_attributes.min  to be float
//   - Updated .orders.order_deliveries.order_items.item.qty_attributes.max  to be float
//   - Updated .orders.order_deliveries.order_items.item.qty_attributes.select.options to be float
//   - Updated .orders.rating to be float
type OrdersResponse struct {
	Orders []struct {
		ID        string  `json:"id"`
		LegacyID  string  `json:"legacy_id"`
		Status    string  `json:"status"`
		Rating    float32 `json:"rating"`
		Total     string  `json:"total"`
		CreatedAt string  `json:"created_at"`
		Actions   map[string]struct {
			Label           string `json:"label"`
			InProgressLabel string `json:"in_progress_label"`
			OrderUUID       string `json:"order_uuid"`
			SourceType      string `json:"source_type"`
		} `json:"actions"`
		OrderDeliveries []struct {
			ID          string `json:"id"`
			OrderID     string `json:"order_id"`
			Description string `json:"description"`
			Base62ID    string `json:"base62_id"`
			Status      string `json:"status"`
			DeliveredAt string `json:"delivered_at"`
			Retailer    struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
				Logo struct {
					URL        string `json:"url"`
					Alt        string `json:"alt"`
					Responsive struct {
						Template string `json:"template"`
						Defaults struct {
							Width int `json:"width"`
						} `json:"defaults"`
					} `json:"responsive"`
					Sizes []interface{} `json:"sizes"`
				} `json:"logo"`
				BackgroundColor string `json:"background_color"`
			} `json:"retailer"`
			OrderItems []struct {
				Qty  float32 `json:"qty"`
				Item struct {
					ID                      string      `json:"id"`
					LegacyID                int         `json:"legacy_id"`
					ProductID               string      `json:"product_id"`
					Name                    string      `json:"name"`
					Attributes              []string    `json:"attributes"`
					PriceAffix              interface{} `json:"price_affix"`
					PriceAffixAria          interface{} `json:"price_affix_aria"`
					SecondaryPriceAffix     string      `json:"secondary_price_affix"`
					SecondaryPriceAffixAria string      `json:"secondary_price_affix_aria"`
					Size                    string      `json:"size"`
					SizeAria                string      `json:"size_aria"`
					ImageList               []struct {
						URL        string `json:"url"`
						Alt        string `json:"alt"`
						Responsive struct {
							Template string `json:"template"`
							Defaults struct {
								Width  int    `json:"width"`
								Fill   string `json:"fill"`
								Format string `json:"format"`
							} `json:"defaults"`
						} `json:"responsive"`
						Sizes []interface{} `json:"sizes"`
					} `json:"image_list"`
					Image struct {
						URL        string `json:"url"`
						Alt        string `json:"alt"`
						Responsive struct {
							Template string `json:"template"`
							Defaults struct {
								Width  int    `json:"width"`
								Fill   string `json:"fill"`
								Format string `json:"format"`
							} `json:"defaults"`
						} `json:"responsive"`
						Sizes []interface{} `json:"sizes"`
					} `json:"image"`
					VariableAttributesMap interface{} `json:"variable_attributes_map"`
					ClickAction           struct {
						Type string `json:"type"`
						Data struct {
							Container struct {
								Title            string        `json:"title"`
								Path             string        `json:"path"`
								InitialStep      interface{}   `json:"initial_step"`
								Modules          []interface{} `json:"modules"`
								DataDependencies []interface{} `json:"data_dependencies"`
							} `json:"container"`
							TrackingParams struct {
							} `json:"tracking_params"`
							TrackingEventNames struct {
							} `json:"tracking_event_names"`
						} `json:"data"`
					} `json:"click_action"`
					WineRatingBadge interface{} `json:"wine_rating_badge"`
					Weekly          interface{} `json:"weekly"`
					WeeklyOrderID   interface{} `json:"weekly_order_id"`
					QtyAttributes   struct {
						Initial          int         `json:"initial"`
						Increment        float32     `json:"increment"`
						Min              float32     `json:"min"`
						Max              float32     `json:"max"`
						Unit             interface{} `json:"unit"`
						UnitAria         interface{} `json:"unit_aria"`
						MaxReachedLabel  string      `json:"max_reached_label"`
						MinReachedLabel  interface{} `json:"min_reached_label"`
						MinWeightExp     bool        `json:"min_weight_exp"`
						Editable         bool        `json:"editable"`
						QtyEnforcedLabel interface{} `json:"qty_enforced_label"`
						Select           struct {
							Options       []float32 `json:"options"`
							DefaultOption int       `json:"default_option"`
							CustomOption  struct {
								Label string `json:"label"`
							} `json:"custom_option"`
						} `json:"select"`
					} `json:"qty_attributes"`
					QtyAttributesPerUnit        interface{} `json:"qty_attributes_per_unit"`
					DeliveryPromotionAttributes interface{} `json:"delivery_promotion_attributes"`
				} `json:"item"`
			} `json:"order_items"`
		} `json:"order_deliveries"`
	} `json:"orders"`
	Meta struct {
		Pagination struct {
			Total    int  `json:"total"`
			PerPage  int  `json:"per_page"`
			Page     int  `json:"page"`
			NextPage *int `json:"next_page"`
		} `json:"pagination"`
	} `json:"meta"`
}
