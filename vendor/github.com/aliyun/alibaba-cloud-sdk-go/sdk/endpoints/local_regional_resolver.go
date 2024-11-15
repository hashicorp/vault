package endpoints

import (
	"fmt"
	"strings"

	"github.com/jmespath/go-jmespath"
)

type LocalRegionalResolver struct {
}

func (resolver *LocalRegionalResolver) GetName() (name string) {
	name = "local regional resolver"
	return
}

func (resolver *LocalRegionalResolver) TryResolve(param *ResolveParam) (endpoint string, support bool, err error) {
	// get the regional endpoints configs
	regionalExpression := fmt.Sprintf("products[?code=='%s'].regional_endpoints", strings.ToLower(param.Product))
	regionalData, err := jmespath.Search(regionalExpression, getEndpointConfigData())
	if err == nil && regionalData != nil && len(regionalData.([]interface{})) > 0 {
		endpointExpression := fmt.Sprintf("[0][?region=='%s'].endpoint", strings.ToLower(param.RegionId))
		var endpointData interface{}
		endpointData, err = jmespath.Search(endpointExpression, regionalData)
		if err == nil && endpointData != nil && len(endpointData.([]interface{})) > 0 {
			endpoint = endpointData.([]interface{})[0].(string)
			support = len(endpoint) > 0
			return
		}
	}
	support = false
	return
}
