package packngo

const spotMarketBasePath = "/market/spot/prices"

// SpotMarketService expooses Spot Market methods
type SpotMarketService interface {
	Prices() (PriceMap, *Response, error)
}

// SpotMarketServiceOp implements SpotMarketService
type SpotMarketServiceOp struct {
	client *Client
}

// PriceMap is a map of [facility][plan]-> float Price
type PriceMap map[string]map[string]float64

// Prices gets current PriceMap from the API
func (s *SpotMarketServiceOp) Prices() (PriceMap, *Response, error) {
	root := new(struct {
		SMPs map[string]map[string]struct {
			Price float64 `json:"price"`
		} `json:"spot_market_prices"`
	})

	resp, err := s.client.DoRequest("GET", spotMarketBasePath, nil, root)
	if err != nil {
		return nil, resp, err
	}

	prices := make(PriceMap)
	for facility, planMap := range root.SMPs {
		prices[facility] = map[string]float64{}
		for plan, v := range planMap {
			prices[facility][plan] = v.Price
		}
	}
	return prices, resp, err
}
