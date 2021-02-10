package common

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
)

type InstanceChargeType string

const (
	PrePaid  = InstanceChargeType("PrePaid")
	PostPaid = InstanceChargeType("PostPaid")
)

var SpecialDeployedProducts = map[string]map[Region]interface{}{
	"vpc": {
		Hangzhou:     Hangzhou,
		Shenzhen:     Shenzhen,
		APSouthEast1: APSouthEast1,
		USWest1:      USWest1,
		USEast1:      USEast1,
		Chengdu:      Chengdu,
		Zhangjiakou:  Zhangjiakou,
		Huhehaote:    Huhehaote,
		APSouthEast3: APSouthEast3,
		EUCentral1:   EUCentral1,
		EUWest1:      EUWest1,
		APSouth1:     APSouth1,
		APNorthEast1: APNorthEast1,
		APSouthEast5: APSouthEast5,
		APSouthEast2: APSouthEast2,
		MEEast1:      MEEast1,
		CNNorth2Gov1: CNNorth2Gov1,
	},
}

var CentralDomainServices = map[string]string{
	"pvtz": "pvtz.vpc-proxy.aliyuncs.com",
}

var RegionalDomainServices = []string{
	"ecs",
	"vpc",
	"slb",
}

// Unit-Domain of central product
var UnitRegions = map[Region]interface{}{
	Hangzhou:     Hangzhou,
	Shenzhen:     Shenzhen,
	APSouthEast1: APSouthEast1,
	USWest1:      USWest1,
	USEast1:      USEast1,
	Chengdu:      Chengdu,
	Zhangjiakou:  Zhangjiakou,
	Huhehaote:    Huhehaote,
	APSouthEast3: APSouthEast3,
	EUCentral1:   EUCentral1,
	EUWest1:      EUWest1,
	APSouth1:     APSouth1,
	APNorthEast1: APNorthEast1,
	APSouthEast5: APSouthEast5,
	APSouthEast2: APSouthEast2,
	CNNorth2Gov1: CNNorth2Gov1,
	//MEEast1:      MEEast1,
	//RUSWest1:        RUSWest1,
	//Beijing:         Beijing,
	//Shanghai:        Shanghai,
	//Hongkong:        Hongkong,
	//ShanghaiFinance: ShanghaiFinance,
	//ShenZhenFinance: ShenZhenFinance,
	HangZhouFinance: Hangzhou,
}

type DescribeEndpointArgs struct {
	Id          Region
	ServiceCode string
	Type        string
}

type EndpointItem struct {
	Protocols struct {
		Protocols []string
	}
	Type        string
	Namespace   string
	Id          Region
	SerivceCode string
	Endpoint    string
}

type DescribeEndpointResponse struct {
	Response
	EndpointItem
}

type DescribeEndpointsArgs struct {
	Id          Region
	ServiceCode string
	Type        string
}

type DescribeEndpointsResponse struct {
	Response
	Endpoints APIEndpoints
	RequestId string
	Success   bool
}

type APIEndpoints struct {
	Endpoint []EndpointItem
}

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type NetworkType string

const (
	Classic = NetworkType("Classic")
	VPC     = NetworkType("VPC")
)

type BusinessInfo struct {
	Pack       string `json:"pack,omitempty"`
	ActivityId string `json:"activityId,omitempty"`
}

//xml
type Endpoints struct {
	Endpoint []Endpoint `xml:"Endpoint"`
}

type Endpoint struct {
	Name      string    `xml:"name,attr"`
	RegionIds RegionIds `xml:"RegionIds"`
	Products  Products  `xml:"Products"`
}

type RegionIds struct {
	RegionId string `xml:"RegionId"`
}

type Products struct {
	Product []Product `xml:"Product"`
}

type Product struct {
	ProductName string `xml:"ProductName"`
	DomainName  string `xml:"DomainName"`
}
