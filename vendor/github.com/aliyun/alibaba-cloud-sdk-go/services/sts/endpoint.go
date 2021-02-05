package sts

// EndpointMap Endpoint Data
var EndpointMap map[string]string

// EndpointType regional or central
var EndpointType = "regional"

// GetEndpointMap Get Endpoint Data Map
func GetEndpointMap() map[string]string {
	if EndpointMap == nil {
		EndpointMap = map[string]string{
			"cn-shanghai-internal-test-1": "sts.aliyuncs.com",
			"cn-beijing-gov-1":            "sts.aliyuncs.com",
			"cn-shenzhen-su18-b01":        "sts.aliyuncs.com",
			"cn-shanghai-inner":           "sts.aliyuncs.com",
			"cn-shenzhen-st4-d01":         "sts.aliyuncs.com",
			"cn-haidian-cm12-c01":         "sts.aliyuncs.com",
			"cn-hangzhou-internal-prod-1": "sts.aliyuncs.com",
			"cn-north-2-gov-1":            "sts.aliyuncs.com",
			"cn-yushanfang":               "sts.aliyuncs.com",
			"cn-hongkong-finance-pop":     "sts.aliyuncs.com",
			"cn-qingdao-nebula":           "sts.aliyuncs.com",
			"cn-shanghai-finance-1":       "sts.aliyuncs.com",
			"cn-beijing-finance-pop":      "sts.aliyuncs.com",
			"cn-wuhan":                    "sts.aliyuncs.com",
			"cn-zhengzhou-nebula-1":       "sts.aliyuncs.com",
			"rus-west-1-pop":              "sts.ap-northeast-1.aliyuncs.com",
			"cn-shanghai-et15-b01":        "sts.aliyuncs.com",
			"cn-hangzhou-bj-b01":          "sts.aliyuncs.com",
			"cn-hangzhou-internal-test-1": "sts.aliyuncs.com",
			"eu-west-1-oxs":               "sts.ap-northeast-1.aliyuncs.com",
			"cn-zhangbei-na61-b01":        "sts.aliyuncs.com",
			"cn-beijing-finance-1":        "sts.aliyuncs.com",
			"cn-hangzhou-internal-test-3": "sts.aliyuncs.com",
			"cn-hangzhou-internal-test-2": "sts.aliyuncs.com",
			"cn-shenzhen-finance-1":       "sts.aliyuncs.com",
			"cn-hangzhou-test-306":        "sts.aliyuncs.com",
			"cn-shanghai-et2-b01":         "sts.aliyuncs.com",
			"cn-hangzhou-finance":         "sts.aliyuncs.com",
			"cn-beijing-nu16-b01":         "sts.aliyuncs.com",
			"cn-edge-1":                   "sts.aliyuncs.com",
			"cn-fujian":                   "sts.aliyuncs.com",
			"ap-northeast-2-pop":          "sts.ap-northeast-1.aliyuncs.com",
			"cn-shenzhen-inner":           "sts.aliyuncs.com",
			"cn-zhangjiakou-na62-a01":     "sts.aliyuncs.com",
		}
	}
	return EndpointMap
}

// GetEndpointType Get Endpoint Type Value
func GetEndpointType() string {
	return EndpointType
}
