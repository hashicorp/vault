erb range_type=Numrange multirange_type=Nummultirange typed_multirange.go.erb > num_multirange.go
erb range_type=Int4range multirange_type=Int4multirange typed_multirange.go.erb > int4_multirange.go
erb range_type=Int8range multirange_type=Int8multirange typed_multirange.go.erb > int8_multirange.go
# TODO
# erb range_type=Tsrange multirange_type=Tsmultirange typed_multirange.go.erb > ts_multirange.go
# erb range_type=Tstzrange multirange_type=Tstzmultirange typed_multirange.go.erb > tstz_multirange.go
# erb range_type=Daterange multirange_type=Datemultirange typed_multirange.go.erb > date_multirange.go
goimports -w *multirange.go