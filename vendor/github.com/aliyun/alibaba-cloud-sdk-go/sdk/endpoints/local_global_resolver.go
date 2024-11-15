package endpoints

import (
	"fmt"
	"strings"

	"github.com/jmespath/go-jmespath"
)

type LocalGlobalResolver struct {
}

func (resolver *LocalGlobalResolver) GetName() (name string) {
	name = "local global resolver"
	return
}

func (resolver *LocalGlobalResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
	// get the global endpoints configs
	endpointExpression := fmt.Sprintf("products[?code=='%s'].global_endpoint", strings.ToLower(param.Product))
	endpointData, err := jmespath.Search(endpointExpression, getEndpointConfigData())
	if err == nil && endpointData != nil && len(endpointData.([]interface{})) > 0 {
		endpoint = endpointData.([]interface{})[0].(string)
		support = len(endpoint) > 0
		return
	}
	support = false
	return
}
