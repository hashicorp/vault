# 1.14.4 (October 18, 2024)

* Update pgx to v4.18.3

# 1.14.3 (March 12, 2024)

* Update pgx to v4.18.2

# 1.14.2 (February 5, 2024)

* Fix numeric to float64 conversion (Joe Schafer)

# 1.14.1 (January 12, 2024)

* Backport fix numeric to string conversion for small negative values
* Fix EncodeValueText (horpto)
* Fix JSON.UnmarshalJSON to make copy of byte slice (horpto)

# 1.14.0 (February 11, 2023)

* Fix: BC timestamp text format support (jozeflami)
* Add Scanner and Valuer interfaces to CIDR (Yurii Popivniak)
* Fix crash when nilifying pointer to sql.Scanner

# 1.13.0 (December 1, 2022)

* Fix: Reset jsonb before unmarshal (Tomas Odinas)
* Fix: return correct zero value when UUID conversion fails (ndrpnt)
* Fix: EncodeText for Lseg includes [ and ]
* Support sql Value and Scan for custom date type (Hubert Krauze)
* Support Ltree binary encoding (AmineChikhaoui)
* Fix: dates with "BC" (jozeflami)

# 1.12.0 (August 6, 2022)

* Add JSONArray (Jakob Ackermann)
* Support Inet from fmt.Stringer and encoding.TextMarshaler (Ville Skyttä)
* Support UUID from fmt.Stringer interface (Lasse Hyldahl Jensen)
* Fix: shopspring-numeric extension does not panic on NaN
* Numeric can be assigned to string
* Fix: Do not send IPv4 networks as IPv4-mapped IPv6 (William Storey)
* Fix: PlanScan for interface{}(nil) (James Hartig)
* Fix: *sql.Scanner for NULL handling (James Hartig)
* Timestamp[tz].Set() supports string (Harmen)
* Fix: Hstore AssignTo with map of *string (Diego Becciolini)

# 1.11.0 (April 21, 2022)

* Add multirange for numeric, int4, and int8 (Vu)
* JSONBArray now supports json.RawMessage (Jens Emil Schulz Østergaard)
* Add RecordArray (WGH)
* Add UnmarshalJSON to pgtype.Int2
* Hstore.Set accepts map[string]Text

# 1.10.0 (February 7, 2022)

* Normalize UTC timestamps to comply with stdlib (Torkel Rogstad)
* Assign Numeric to *big.Rat (Oleg Lomaka)
* Fix typo in float8 error message (Pinank Solanki)
* Scan type aliases for floating point types (Collin Forsyth)

# 1.9.1 (November 28, 2021)

* Fix: binary timestamp is assumed to be in UTC (restored behavior changed in v1.9.0)

# 1.9.0 (November 20, 2021)

* Fix binary hstore null decoding
* Add shopspring/decimal.NullDecimal support to integration (Eli Treuherz)
* Inet.Set supports bare IP address (Carl Dunham)
* Add zeronull.Float8
* Fix NULL being lost when scanning unknown OID into sql.Scanner
* Fix BPChar.AssignTo **rune
* Add support for fmt.Stringer and driver.Valuer in String fields encoding (Jan Dubsky)
* Fix really big timestamp(tz)s binary format parsing (e.g. year 294276) (Jim Tsao)
* Support `map[string]*string` as hstore (Adrian Sieger)
* Fix parsing text array with negative bounds
* Add infinity support for numeric (Jim Tsao)

# 1.8.1 (July 24, 2021)

* Cleaned up Go module dependency chain

# 1.8.0 (July 10, 2021)

* Maintain host bits for inet types (Cameron Daniel)
* Support pointers of wrapping structs (Ivan Daunis)
* Register JSONBArray at NewConnInfo() (Rueian)
* CompositeTextScanner handles backslash escapes

# 1.7.0 (March 25, 2021)

* Fix scanning int into **sql.Scanner implementor
* Add tsrange array type (Vasilii Novikov)
* Fix: escaped strings when they start or end with a newline char (Stephane Martin)
* Accept nil *time.Time in Time.Set
* Fix numeric NaN support
* Use Go 1.13 errors instead of xerrors

# 1.6.2 (December 3, 2020)

* Fix panic on assigning empty array to non-slice or array
* Fix text array parsing disambiguates NULL and "NULL"
* Fix Timestamptz.DecodeText with too short text

# 1.6.1 (October 31, 2020)

* Fix simple protocol empty array support

# 1.6.0 (October 24, 2020)

* Fix AssignTo pointer to pointer to slice and named types.
* Fix zero length array assignment (Simo Haasanen)
* Add float64, float32 convert to int2, int4, int8 (lqu3j)
* Support setting infinite timestamps (Erik Agsjö)
* Polygon improvements (duohedron)
* Fix Inet.Set with nil (Tomas Volf)

# 1.5.0 (September 26, 2020)

* Add slice of slice mapping to multi-dimensional arrays (Simo Haasanen)
* Fix JSONBArray
* Fix selecting empty array
* Text formatted values except bytea can be directly scanned to []byte
* Add JSON marshalling for UUID (bakmataliev)
* Improve point type conversions (bakmataliev)

# 1.4.2 (July 22, 2020)

* Fix encoding of a large composite data type (Yaz Saito)

# 1.4.1 (July 14, 2020)

* Fix ArrayType DecodeBinary empty array breaks future reads

# 1.4.0 (June 27, 2020)

* Add JSON support to ext/gofrs-uuid
* Performance improvements in Scan path
* Improved ext/shopspring-numeric binary decoding performance
* Add composite type support (Maxim Ivanov and Jack Christensen)
* Add better generic enum type support
* Add generic array type support
* Clarify and normalize Value semantics
* Fix hstore with empty string values
* Numeric supports NaN values (leighhopcroft)
* Add slice of pointer support to array types (megaturbo)
* Add jsonb array type (tserakhau)
* Allow converting intervals with months and days to duration

# 1.3.0 (March 30, 2020)

* Get implemented on T instead of *T
* Set will call Get on src if possible
* Range types Set method supports its own type, string, and nil
* Date.Set parses string
* Fix correct format verb for unknown type error (Robert Welin)
* Truncate nanoseconds in EncodeText for Timestamptz and Timestamp

# 1.2.0 (February 5, 2020)

* Add zeronull package for easier NULL <-> zero conversion
* Add JSON marshalling for shopspring-numeric extension
* Add JSON marshalling for Bool, Date, JSON/B, Timestamptz (Jeffrey Stiles)
* Fix null status in UnmarshalJSON for some types (Jeffrey Stiles)

# 1.1.0 (January 11, 2020)

* Add PostgreSQL time type support
* Add more automatic conversions of integer arrays of different types (Jean-Philippe Quéméner)

# 1.0.3 (November 16, 2019)

* Support initializing Array types from a slice of the value (Alex Gaynor)

# 1.0.2 (October 22, 2019)

* Fix scan into null into pointer to pointer implementing Decode* interface. (Jeremy Altavilla)

# 1.0.1 (September 19, 2019)

* Fix daterange OID
