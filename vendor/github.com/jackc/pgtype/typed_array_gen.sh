erb pgtype_array_type=Int2Array pgtype_element_type=Int2 go_array_types=[]int16,[]*int16,[]uint16,[]*uint16,[]int32,[]*int32,[]uint32,[]*uint32,[]int64,[]*int64,[]uint64,[]*uint64,[]int,[]*int,[]uint,[]*uint element_type_name=int2 typed_array.go.erb > int2_array.go
erb pgtype_array_type=Int4Array pgtype_element_type=Int4 go_array_types=[]int16,[]*int16,[]uint16,[]*uint16,[]int32,[]*int32,[]uint32,[]*uint32,[]int64,[]*int64,[]uint64,[]*uint64,[]int,[]*int,[]uint,[]*uint element_type_name=int4 typed_array.go.erb > int4_array.go
erb pgtype_array_type=Int8Array pgtype_element_type=Int8 go_array_types=[]int16,[]*int16,[]uint16,[]*uint16,[]int32,[]*int32,[]uint32,[]*uint32,[]int64,[]*int64,[]uint64,[]*uint64,[]int,[]*int,[]uint,[]*uint element_type_name=int8 typed_array.go.erb > int8_array.go
erb pgtype_array_type=BoolArray pgtype_element_type=Bool go_array_types=[]bool,[]*bool element_type_name=bool typed_array.go.erb > bool_array.go
erb pgtype_array_type=DateArray pgtype_element_type=Date go_array_types=[]time.Time,[]*time.Time element_type_name=date typed_array.go.erb > date_array.go
erb pgtype_array_type=TimestamptzArray pgtype_element_type=Timestamptz go_array_types=[]time.Time,[]*time.Time element_type_name=timestamptz typed_array.go.erb > timestamptz_array.go
erb pgtype_array_type=TstzrangeArray pgtype_element_type=Tstzrange go_array_types=[]Tstzrange element_type_name=tstzrange typed_array.go.erb > tstzrange_array.go
erb pgtype_array_type=TsrangeArray pgtype_element_type=Tsrange go_array_types=[]Tsrange element_type_name=tsrange typed_array.go.erb > tsrange_array.go
erb pgtype_array_type=TimestampArray pgtype_element_type=Timestamp go_array_types=[]time.Time,[]*time.Time element_type_name=timestamp typed_array.go.erb > timestamp_array.go
erb pgtype_array_type=Float4Array pgtype_element_type=Float4 go_array_types=[]float32,[]*float32 element_type_name=float4 typed_array.go.erb > float4_array.go
erb pgtype_array_type=Float8Array pgtype_element_type=Float8 go_array_types=[]float64,[]*float64 element_type_name=float8 typed_array.go.erb > float8_array.go
erb pgtype_array_type=InetArray pgtype_element_type=Inet go_array_types=[]*net.IPNet,[]net.IP,[]*net.IP element_type_name=inet typed_array.go.erb > inet_array.go
erb pgtype_array_type=MacaddrArray pgtype_element_type=Macaddr go_array_types=[]net.HardwareAddr,[]*net.HardwareAddr element_type_name=macaddr typed_array.go.erb > macaddr_array.go
erb pgtype_array_type=CIDRArray pgtype_element_type=CIDR go_array_types=[]*net.IPNet,[]net.IP,[]*net.IP element_type_name=cidr typed_array.go.erb > cidr_array.go
erb pgtype_array_type=TextArray pgtype_element_type=Text go_array_types=[]string,[]*string element_type_name=text typed_array.go.erb > text_array.go
erb pgtype_array_type=VarcharArray pgtype_element_type=Varchar go_array_types=[]string,[]*string element_type_name=varchar typed_array.go.erb > varchar_array.go
erb pgtype_array_type=BPCharArray pgtype_element_type=BPChar go_array_types=[]string,[]*string element_type_name=bpchar typed_array.go.erb > bpchar_array.go
erb pgtype_array_type=ByteaArray pgtype_element_type=Bytea go_array_types=[][]byte element_type_name=bytea typed_array.go.erb > bytea_array.go
erb pgtype_array_type=ACLItemArray pgtype_element_type=ACLItem go_array_types=[]string,[]*string element_type_name=aclitem binary_format=false typed_array.go.erb > aclitem_array.go
erb pgtype_array_type=HstoreArray pgtype_element_type=Hstore go_array_types=[]map[string]string element_type_name=hstore typed_array.go.erb > hstore_array.go
erb pgtype_array_type=NumericArray pgtype_element_type=Numeric go_array_types=[]float32,[]*float32,[]float64,[]*float64,[]int64,[]*int64,[]uint64,[]*uint64 element_type_name=numeric typed_array.go.erb > numeric_array.go
erb pgtype_array_type=UUIDArray pgtype_element_type=UUID go_array_types=[][16]byte,[][]byte,[]string,[]*string element_type_name=uuid typed_array.go.erb > uuid_array.go
erb pgtype_array_type=JSONArray pgtype_element_type=JSON go_array_types=[]string,[][]byte,[]json.RawMessage element_type_name=json typed_array.go.erb > json_array.go
erb pgtype_array_type=JSONBArray pgtype_element_type=JSONB go_array_types=[]string,[][]byte,[]json.RawMessage element_type_name=jsonb typed_array.go.erb > jsonb_array.go

# While the binary format is theoretically possible it is only practical to use the text format.
erb pgtype_array_type=EnumArray pgtype_element_type=GenericText go_array_types=[]string,[]*string binary_format=false typed_array.go.erb > enum_array.go

erb pgtype_array_type=RecordArray pgtype_element_type=Record go_array_types=[][]Value element_type_name=record text_null=NULL encode_binary=false text_format=false typed_array.go.erb > record_array.go

goimports -w *_array.go
