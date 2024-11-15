package packngo

import "path"

const (
	spotMarketBasePath = "/market/spot/prices"

	spotMarketMetrosPath = "metros"
)

// SpotMarketService expooses Spot Market methods
type SpotMarketService interface {
	// Prices gets current spot market prices by facility.
	//
	// Deprecated: Use PricesByFacility
	Prices() (PriceMap, *Response, error)

	// PricesByFacility gets current spot market prices by facility. The map is
	// indexed by facility code and then plan name.
	PricesByFacility() (PriceMap, *Response, error)

	// PricesByMetro gets current spot market prices by metro. The map is
	// indexed by metro code and then plan name.
	PricesByMetro() (PriceMap, *Response, error)
}

// SpotMarketServiceOp implements SpotMarketService
type SpotMarketServiceOp struct {
	client *Client
}

// PriceMap is a map of [location][plan]-> float Price
type PriceMap map[string]map[string]float64

// Prices gets current spot market prices by facility.
//
// Deprecated: Use PricesByFacility which this function thinly wraps.
func (s *SpotMarketServiceOp) Prices() (PriceMap, *Response, error) {
	return s.PricesByFacility()
}

// PricesByFacility gets current spot market prices by facility. The map is
// indexed by facility code and then plan name.
//
// price := client.SpotMarket.PricesByFacility()["ny5"]["c3.medium.x86"]
func (s *SpotMarketServiceOp) PricesByFacility() (PriceMap, *Response, error) {
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

// PricesByMetro gets current spot market prices by metro. The map is
// indexed by metro code and then plan name.
//
// price := client.SpotMarket.PricesByMetro()["sv"]["c3.medium.x86"]
func (s *SpotMarketServiceOp) PricesByMetro() (PriceMap, *Response, error) {
	root := new(struct {
		SMPs map[string]map[string]struct {
			Price float64 `json:"price"`
		} `json:"spot_market_prices"`
	})

	resp, err := s.client.DoRequest("GET", path.Join(spotMarketBasePath, spotMarketMetrosPath), nil, root)
	if err != nil {
		return nil, resp, err
	}

	prices := make(PriceMap)
	for metro, planMap := range root.SMPs {
		prices[metro] = map[string]float64{}
		for plan, v := range planMap {
			prices[metro][plan] = v.Price
		}
	}
	return prices, resp, err
}
