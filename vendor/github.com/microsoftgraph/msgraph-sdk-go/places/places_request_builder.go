package places

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
)

// PlacesRequestBuilder builds and executes requests for operations under \places
type PlacesRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// ByPlaceId provides operations to manage the collection of place entities.
// returns a *PlaceItemRequestBuilder when successful
func (m *PlacesRequestBuilder) ByPlaceId(placeId string)(*PlaceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.BaseRequestBuilder.PathParameters {
        urlTplParams[idx] = item
    }
    if placeId != "" {
        urlTplParams["place%2Did"] = placeId
    }
    return NewPlaceItemRequestBuilderInternal(urlTplParams, m.BaseRequestBuilder.RequestAdapter)
}
// NewPlacesRequestBuilderInternal instantiates a new PlacesRequestBuilder and sets the default values.
func NewPlacesRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PlacesRequestBuilder) {
    m := &PlacesRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/places", pathParameters),
    }
    return m
}
// NewPlacesRequestBuilder instantiates a new PlacesRequestBuilder and sets the default values.
func NewPlacesRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*PlacesRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewPlacesRequestBuilderInternal(urlParams, requestAdapter)
}
// Count provides operations to count the resources in the collection.
// returns a *CountRequestBuilder when successful
func (m *PlacesRequestBuilder) Count()(*CountRequestBuilder) {
    return NewCountRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphRoom casts the previous resource to room.
// returns a *GraphRoomRequestBuilder when successful
func (m *PlacesRequestBuilder) GraphRoom()(*GraphRoomRequestBuilder) {
    return NewGraphRoomRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// GraphRoomList casts the previous resource to roomList.
// returns a *GraphRoomListRequestBuilder when successful
func (m *PlacesRequestBuilder) GraphRoomList()(*GraphRoomListRequestBuilder) {
    return NewGraphRoomListRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
