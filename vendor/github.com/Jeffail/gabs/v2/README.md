![Gabs](gabs_logo.png "Gabs")

[![pkg.go for Jeffail/gabs][godoc-badge]][godoc-url]

Gabs is a small utility for dealing with dynamic or unknown JSON structures in Go. It's pretty much just a helpful wrapper for navigating hierarchies of `map[string]interface{}` objects provided by the `encoding/json` package. It does nothing spectacular apart from being fabulous.

If you're migrating from version 1 check out [`migration.md`][migration-doc] for guidance.

## Use

### Import

Using modules:

```go
import (
	"github.com/Jeffail/gabs/v2"
)
```

Without modules:

```go
import (
	"github.com/Jeffail/gabs"
)
```

### Parsing and searching JSON

```go
jsonParsed, err := gabs.ParseJSON([]byte(`{
	"outer":{
		"inner":{
			"value1":10,
			"value2":22
		},
		"alsoInner":{
			"value1":20,
			"array1":[
				30, 40
			]
		}
	}
}`))
if err != nil {
	panic(err)
}

var value float64
var ok bool

value, ok = jsonParsed.Path("outer.inner.value1").Data().(float64)
// value == 10.0, ok == true

value, ok = jsonParsed.Search("outer", "inner", "value1").Data().(float64)
// value == 10.0, ok == true

value, ok = jsonParsed.Search("outer", "alsoInner", "array1", "1").Data().(float64)
// value == 40.0, ok == true

gObj, err := jsonParsed.JSONPointer("/outer/alsoInner/array1/1")
if err != nil {
	panic(err)
}
value, ok = gObj.Data().(float64)
// value == 40.0, ok == true

value, ok = jsonParsed.Path("does.not.exist").Data().(float64)
// value == 0.0, ok == false

exists := jsonParsed.Exists("outer", "inner", "value1")
// exists == true

exists = jsonParsed.ExistsP("does.not.exist")
// exists == false
```

### Iterating objects

```go
jsonParsed, err := gabs.ParseJSON([]byte(`{"object":{"first":1,"second":2,"third":3}}`))
if err != nil {
	panic(err)
}

// S is shorthand for Search
for key, child := range jsonParsed.S("object").ChildrenMap() {
	fmt.Printf("key: %v, value: %v\n", key, child.Data().(float64))
}
```

### Iterating arrays

```go
jsonParsed, err := gabs.ParseJSON([]byte(`{"array":["first","second","third"]}`))
if err != nil {
	panic(err)
}

for _, child := range jsonParsed.S("array").Children() {
	fmt.Println(child.Data().(string))
}
```

Will print:

```
first
second
third
```

Children() will return all children of an array in order. This also works on objects, however, the children will be returned in a random order.

### Searching through arrays

If your structure contains arrays you must target an index in your search.

```go
jsonParsed, err := gabs.ParseJSON([]byte(`{"array":[{"value":1},{"value":2},{"value":3}]}`))
if err != nil {
	panic(err)
}
fmt.Println(jsonParsed.Path("array.1.value").String())
```

Will print `2`.

### Generating JSON

```go
jsonObj := gabs.New()
// or gabs.Wrap(jsonObject) to work on an existing map[string]interface{}

jsonObj.Set(10, "outer", "inner", "value")
jsonObj.SetP(20, "outer.inner.value2")
jsonObj.Set(30, "outer", "inner2", "value3")

fmt.Println(jsonObj.String())
```

Will print:

```
{"outer":{"inner":{"value":10,"value2":20},"inner2":{"value3":30}}}
```

To pretty-print:

```go
fmt.Println(jsonObj.StringIndent("", "  "))
```

Will print:

```
{
  "outer": {
    "inner": {
      "value": 10,
      "value2": 20
    },
    "inner2": {
      "value3": 30
    }
  }
}
```

### Generating Arrays

```go
jsonObj := gabs.New()

jsonObj.Array("foo", "array")
// Or .ArrayP("foo.array")

jsonObj.ArrayAppend(10, "foo", "array")
jsonObj.ArrayAppend(20, "foo", "array")
jsonObj.ArrayAppend(30, "foo", "array")

fmt.Println(jsonObj.String())
```

Will print:

```
{"foo":{"array":[10,20,30]}}
```

Working with arrays by index:

```go
jsonObj := gabs.New()

// Create an array with the length of 3
jsonObj.ArrayOfSize(3, "foo")

jsonObj.S("foo").SetIndex("test1", 0)
jsonObj.S("foo").SetIndex("test2", 1)

// Create an embedded array with the length of 3
jsonObj.S("foo").ArrayOfSizeI(3, 2)

jsonObj.S("foo").Index(2).SetIndex(1, 0)
jsonObj.S("foo").Index(2).SetIndex(2, 1)
jsonObj.S("foo").Index(2).SetIndex(3, 2)

fmt.Println(jsonObj.String())
```

Will print:

```
{"foo":["test1","test2",[1,2,3]]}
```

### Converting back to JSON

This is the easiest part:

```go
jsonParsedObj, _ := gabs.ParseJSON([]byte(`{
	"outer":{
		"values":{
			"first":10,
			"second":11
		}
	},
	"outer2":"hello world"
}`))

jsonOutput := jsonParsedObj.String()
// Becomes `{"outer":{"values":{"first":10,"second":11}},"outer2":"hello world"}`
```

And to serialize a specific segment is as simple as:

```go
jsonParsedObj := gabs.ParseJSON([]byte(`{
	"outer":{
		"values":{
			"first":10,
			"second":11
		}
	},
	"outer2":"hello world"
}`))

jsonOutput := jsonParsedObj.Search("outer").String()
// Becomes `{"values":{"first":10,"second":11}}`
```

### Merge two containers

You can merge a JSON structure into an existing one, where collisions will be converted into a JSON array.

```go
jsonParsed1, _ := ParseJSON([]byte(`{"outer":{"value1":"one"}}`))
jsonParsed2, _ := ParseJSON([]byte(`{"outer":{"inner":{"value3":"three"}},"outer2":{"value2":"two"}}`))

jsonParsed1.Merge(jsonParsed2)
// Becomes `{"outer":{"inner":{"value3":"three"},"value1":"one"},"outer2":{"value2":"two"}}`
```

Arrays are merged:

```go
jsonParsed1, _ := ParseJSON([]byte(`{"array":["one"]}`))
jsonParsed2, _ := ParseJSON([]byte(`{"array":["two"]}`))

jsonParsed1.Merge(jsonParsed2)
// Becomes `{"array":["one", "two"]}`
```

### Parsing Numbers

Gabs uses the `json` package under the bonnet, which by default will parse all number values into `float64`. If you need to parse `Int` values then you should use a [`json.Decoder`](https://golang.org/pkg/encoding/json/#Decoder):

```go
sample := []byte(`{"test":{"int":10,"float":6.66}}`)
dec := json.NewDecoder(bytes.NewReader(sample))
dec.UseNumber()

val, err := gabs.ParseJSONDecoder(dec)
if err != nil {
    t.Errorf("Failed to parse: %v", err)
    return
}

intValue, err := val.Path("test.int").Data().(json.Number).Int64()
```

[godoc-badge]: https://godoc.org/github.com/Jeffail/gabs?status.svg
[godoc-url]: https://pkg.go.dev/github.com/Jeffail/gabs/v2
[migration-doc]: ./migration.md
