erb range_type=Int4range element_type=Int4 typed_range.go.erb > int4range.go
erb range_type=Int8range element_type=Int8 typed_range.go.erb > int8range.go
erb range_type=Tsrange element_type=Timestamp typed_range.go.erb > tsrange.go
erb range_type=Tstzrange element_type=Timestamptz typed_range.go.erb > tstzrange.go
erb range_type=Daterange element_type=Date typed_range.go.erb > daterange.go
erb range_type=Numrange element_type=Numeric typed_range.go.erb > numrange.go
goimports -w *range.go
