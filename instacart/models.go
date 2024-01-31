package instacart

// OrdersResponse is the response from the orders API
type OrdersResponse struct {
	Orders []Order    `json:"orders"`
	Meta   OrdersMeta `json:"meta"`
}

type Order struct {
	ID              string            `json:"id"`
	LegacyID        string            `json:"legacy_id"`
	Status          string            `json:"status"`
	Rating          float32           `json:"rating"`
	Total           string            `json:"total"`
	CreatedAt       string            `json:"created_at"`
	Actions         map[string]Action `json:"actions"`
	OrderDeliveries []OrderDelivery   `json:"order_deliveries"`
}

type OrdersMeta struct {
	Pagination struct {
		Total    int  `json:"total"`
		PerPage  int  `json:"per_page"`
		Page     int  `json:"page"`
		NextPage *int `json:"next_page"`
	} `json:"pagination"`
}

type Action struct {
	Label           string `json:"label"`
	InProgressLabel string `json:"in_progress_label"`
	OrderUUID       string `json:"order_uuid"`
	SourceType      string `json:"source_type"`
}

type OrderDelivery struct {
	ID          string      `json:"id"`
	OrderID     string      `json:"order_id"`
	Description string      `json:"description"`
	Base62ID    string      `json:"base62_id"`
	Status      string      `json:"status"`
	DeliveredAt string      `json:"delivered_at"`
	Retailer    Retailer    `json:"retailer"`
	OrderItems  []OrderItem `json:"order_items"`
}

type Retailer struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Slug            string       `json:"slug"`
	Logo            RetailerLogo `json:"logo"`
	BackgroundColor string       `json:"background_color"`
}
type RetailerLogo struct {
	URL        string `json:"url"`
	Alt        string `json:"alt"`
	Responsive struct {
		Template string `json:"template"`
		Defaults struct {
			Width int `json:"width"`
		} `json:"defaults"`
	} `json:"responsive"`
	Sizes []interface{} `json:"sizes"`
}

type OrderItem struct {
	Qty  float32 `json:"qty"`
	Item Item    `json:"item"`
}

type Item struct {
	ID                          string            `json:"id"`
	LegacyID                    int               `json:"legacy_id"`
	ProductID                   string            `json:"product_id"`
	Name                        string            `json:"name"`
	Attributes                  []string          `json:"attributes"`
	PriceAffix                  interface{}       `json:"price_affix"`
	PriceAffixAria              interface{}       `json:"price_affix_aria"`
	SecondaryPriceAffix         string            `json:"secondary_price_affix"`
	SecondaryPriceAffixAria     string            `json:"secondary_price_affix_aria"`
	Size                        string            `json:"size"`
	SizeAria                    string            `json:"size_aria"`
	ImageList                   []ItemImage       `json:"image_list"`
	Image                       ItemImage         `json:"image"`
	VariableAttributesMap       interface{}       `json:"variable_attributes_map"`
	ClickAction                 ItemClickAction   `json:"click_action"`
	WineRatingBadge             interface{}       `json:"wine_rating_badge"`
	Weekly                      interface{}       `json:"weekly"`
	WeeklyOrderID               interface{}       `json:"weekly_order_id"`
	QtyAttributes               ItemQtyAttributes `json:"qty_attributes"`
	QtyAttributesPerUnit        interface{}       `json:"qty_attributes_per_unit"`
	DeliveryPromotionAttributes interface{}       `json:"delivery_promotion_attributes"`
}

type ItemImage struct {
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
}

type ItemClickAction struct {
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
}

type ItemQtyAttributes struct {
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
}
