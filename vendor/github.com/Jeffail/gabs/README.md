![Gabs](gabs_logo.png "Gabs")

Gabs is a small utility for dealing with dynamic or unknown JSON structures in golang. It's pretty much just a helpful wrapper around the golang json.Marshal/json.Unmarshal behaviour and map[string]interface{} objects. It does nothing spectacular except for being fabulous.

https://godoc.org/github.com/Jeffail/gabs

## How to install:

```bash
go get github.com/Jeffail/gabs
```

## How to use

### Parsing and searching JSON

```go
...

import "github.com/Jeffail/gabs"

jsonParsed, err := gabs.ParseJSON([]byte(`{
	"outter":{
		"inner":{
			"value1":10,
			"value2":22
		},
		"alsoInner":{
			"value1":20
		}
	}
}`))

var value float64
var ok bool

value, ok = jsonParsed.Path("outter.inner.value1").Data().(float64)
// value == 10.0, ok == true

value, ok = jsonParsed.Search("outter", "inner", "value1").Data().(float64)
// value == 10.0, ok == true

value, ok = jsonParsed.Path("does.not.exist").Data().(float64)
// value == 0.0, ok == false

exists := jsonParsed.Exists("outter", "inner", "value1")
// exists == true

exists := jsonParsed.Exists("does", "not", "exist")
// exists == false

exists := jsonParsed.ExistsP("does.not.exist")
// exists == false

...
```

### Iterating objects

```go
...

jsonParsed, _ := gabs.ParseJSON([]byte(`{"object":{ "first": 1, "second": 2, "third": 3 }}`))

// S is shorthand for Search
children, _ := jsonParsed.S("object").ChildrenMap()
for key, child := range children {
	fmt.Printf("key: %v, value: %v\n", key, child.Data().(string))
}

...
```

### Iterating arrays

```go
...

jsonParsed, _ := gabs.ParseJSON([]byte(`{"array":[ "first", "second", "third" ]}`))

// S is shorthand for Search
children, _ := jsonParsed.S("array").Children()
for _, child := range children {
	fmt.Println(child.Data().(string))
}

...
```

Will print:

```
first
second
third
```

Children() will return all children of an array in order. This also works on objects, however, the children will be returned in a random order.

### Searching through arrays

If your JSON structure contains arrays you can still search the fields of the objects within the array, this returns a JSON array containing the results for each element.

```go
...

jsonParsed, _ := gabs.ParseJSON([]byte(`{"array":[ {"value":1}, {"value":2}, {"value":3} ]}`))
fmt.Println(jsonParsed.Path("array.value").String())

...
```

Will print:

```
[1,2,3]
```

### Generating JSON

```go
...

jsonObj := gabs.New()
// or gabs.Consume(jsonObject) to work on an existing map[string]interface{}

jsonObj.Set(10, "outter", "inner", "value")
jsonObj.SetP(20, "outter.inner.value2")
jsonObj.Set(30, "outter", "inner2", "value3")

fmt.Println(jsonObj.String())

...
```

Will print:

```
{"outter":{"inner":{"value":10,"value2":20},"inner2":{"value3":30}}}
```

To pretty-print:

```go
...

fmt.Println(jsonObj.StringIndent("", "  "))

...
```

Will print:

```
{
  "outter": {
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
...

jsonObj := gabs.New()

jsonObj.Array("foo", "array")
// Or .ArrayP("foo.array")

jsonObj.ArrayAppend(10, "foo", "array")
jsonObj.ArrayAppend(20, "foo", "array")
jsonObj.ArrayAppend(30, "foo", "array")

fmt.Println(jsonObj.String())

...
```

Will print:

```
{"foo":{"array":[10,20,30]}}
```

Working with arrays by index:

```go
...

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

...
```

Will print:

```
{"foo":["test1","test2",[1,2,3]]}
```

### Converting back to JSON

This is the easiest part:

```go
...

jsonParsedObj, _ := gabs.ParseJSON([]byte(`{
	"outter":{
		"values":{
			"first":10,
			"second":11
		}
	},
	"outter2":"hello world"
}`))

jsonOutput := jsonParsedObj.String()
// Becomes `{"outter":{"values":{"first":10,"second":11}},"outter2":"hello world"}`

...
```

And to serialize a specific segment is as simple as:

```go
...

jsonParsedObj := gabs.ParseJSON([]byte(`{
	"outter":{
		"values":{
			"first":10,
			"second":11
		}
	},
	"outter2":"hello world"
}`))

jsonOutput := jsonParsedObj.Search("outter").String()
// Becomes `{"values":{"first":10,"second":11}}`

...
```

### Parsing Numbers

Gabs uses the `json` package under the bonnet, which by default will parse all number values into `float64`. If you need to parse `Int` values then you should use a `json.Decoder` (https://golang.org/pkg/encoding/json/#Decoder):

```go
sample := []byte(`{"test":{"int":10, "float":6.66}}`)
dec := json.NewDecoder(bytes.NewReader(sample))
dec.UseNumber()

val, err := gabs.ParseJSONDecoder(dec)
if err != nil {
    t.Errorf("Failed to parse: %v", err)
    return
}

intValue, err := val.Path("test.int").Data().(json.Number).Int64()
```
