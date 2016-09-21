/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package gabs implements a simplified wrapper around creating and parsing JSON.
package gabs

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

/*---------------------------------------------------------------------------------------------------
 */

var (
	// ErrOutOfBounds - Index out of bounds.
	ErrOutOfBounds = errors.New("out of bounds")

	// ErrNotObjOrArray - The target is not an object or array type.
	ErrNotObjOrArray = errors.New("not an object or array")

	// ErrNotObj - The target is not an object type.
	ErrNotObj = errors.New("not an object")

	// ErrNotArray - The target is not an array type.
	ErrNotArray = errors.New("not an array")

	// ErrPathCollision - Creating a path failed because an element collided with an existing value.
	ErrPathCollision = errors.New("encountered value collision whilst building path")

	// ErrInvalidInputObj - The input value was not a map[string]interface{}.
	ErrInvalidInputObj = errors.New("invalid input object")

	// ErrInvalidInputText - The input data could not be parsed.
	ErrInvalidInputText = errors.New("input text could not be parsed")

	// ErrInvalidPath - The filepath was not valid.
	ErrInvalidPath = errors.New("invalid file path")

	// ErrInvalidBuffer - The input buffer contained an invalid JSON string
	ErrInvalidBuffer = errors.New("input buffer contained invalid JSON")
)

/*---------------------------------------------------------------------------------------------------
 */

/*
Container - an internal structure that holds a reference to the core interface map of the parsed
json. Use this container to move context.
*/
type Container struct {
	object interface{}
}

/*
Data - Return the contained data as an interface{}.
*/
func (g *Container) Data() interface{} {
	return g.object
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Path - Search for a value using dot notation.
*/
func (g *Container) Path(path string) *Container {
	return g.Search(strings.Split(path, ".")...)
}

/*
Search - Attempt to find and return an object within the JSON structure by specifying the hierarchy
of field names to locate the target. If the search encounters an array and has not reached the end
target then it will iterate each object of the array for the target and return all of the results in
a JSON array.
*/
func (g *Container) Search(hierarchy ...string) *Container {
	var object interface{}

	object = g.object
	for target := 0; target < len(hierarchy); target++ {
		if mmap, ok := object.(map[string]interface{}); ok {
			object = mmap[hierarchy[target]]
		} else if marray, ok := object.([]interface{}); ok {
			tmpArray := []interface{}{}
			for _, val := range marray {
				tmpGabs := &Container{val}
				res := tmpGabs.Search(hierarchy[target:]...).Data()
				if res != nil {
					tmpArray = append(tmpArray, res)
				}
			}
			if len(tmpArray) == 0 {
				return &Container{nil}
			}
			return &Container{tmpArray}
		} else {
			return &Container{nil}
		}
	}
	return &Container{object}
}

/*
S - Shorthand method, does the same thing as Search.
*/
func (g *Container) S(hierarchy ...string) *Container {
	return g.Search(hierarchy...)
}

/*
Exists - Checks whether a path exists.
*/
func (g *Container) Exists(hierarchy ...string) bool {
	return g.Search(hierarchy...).Data() != nil
}

/*
ExistsP - Checks whether a dot notation path exists.
*/
func (g *Container) ExistsP(path string) bool {
	return g.Exists(strings.Split(path, ".")...)
}

/*
Index - Attempt to find and return an object with a JSON array by specifying the index of the
target.
*/
func (g *Container) Index(index int) *Container {
	if array, ok := g.Data().([]interface{}); ok {
		if index >= len(array) {
			return &Container{nil}
		}
		return &Container{array[index]}
	}
	return &Container{nil}
}

/*
Children - Return a slice of all the children of the array. This also works for objects, however,
the children returned for an object will NOT be in order and you lose the names of the returned
objects this way.
*/
func (g *Container) Children() ([]*Container, error) {
	if array, ok := g.Data().([]interface{}); ok {
		children := make([]*Container, len(array))
		for i := 0; i < len(array); i++ {
			children[i] = &Container{array[i]}
		}
		return children, nil
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		children := []*Container{}
		for _, obj := range mmap {
			children = append(children, &Container{obj})
		}
		return children, nil
	}
	return nil, ErrNotObjOrArray
}

/*
ChildrenMap - Return a map of all the children of an object.
*/
func (g *Container) ChildrenMap() (map[string]*Container, error) {
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		children := map[string]*Container{}
		for name, obj := range mmap {
			children[name] = &Container{obj}
		}
		return children, nil
	}
	return nil, ErrNotObj
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Set - Set the value of a field at a JSON path, any parts of the path that do not exist will be
constructed, and if a collision occurs with a non object type whilst iterating the path an error is
returned.
*/
func (g *Container) Set(value interface{}, path ...string) (*Container, error) {
	if len(path) == 0 {
		g.object = value
		return g, nil
	}
	var object interface{}
	if g.object == nil {
		g.object = map[string]interface{}{}
	}
	object = g.object
	for target := 0; target < len(path); target++ {
		if mmap, ok := object.(map[string]interface{}); ok {
			if target == len(path)-1 {
				mmap[path[target]] = value
			} else if mmap[path[target]] == nil {
				mmap[path[target]] = map[string]interface{}{}
			}
			object = mmap[path[target]]
		} else {
			return &Container{nil}, ErrPathCollision
		}
	}
	return &Container{object}, nil
}

/*
SetP - Does the same as Set, but using a dot notation JSON path.
*/
func (g *Container) SetP(value interface{}, path string) (*Container, error) {
	return g.Set(value, strings.Split(path, ".")...)
}

/*
SetIndex - Set a value of an array element based on the index.
*/
func (g *Container) SetIndex(value interface{}, index int) (*Container, error) {
	if array, ok := g.Data().([]interface{}); ok {
		if index >= len(array) {
			return &Container{nil}, ErrOutOfBounds
		}
		array[index] = value
		return &Container{array[index]}, nil
	}
	return &Container{nil}, ErrNotArray
}

/*
Object - Create a new JSON object at a path. Returns an error if the path contains a collision with
a non object type.
*/
func (g *Container) Object(path ...string) (*Container, error) {
	return g.Set(map[string]interface{}{}, path...)
}

/*
ObjectP - Does the same as Object, but using a dot notation JSON path.
*/
func (g *Container) ObjectP(path string) (*Container, error) {
	return g.Object(strings.Split(path, ".")...)
}

/*
ObjectI - Create a new JSON object at an array index. Returns an error if the object is not an array
or the index is out of bounds.
*/
func (g *Container) ObjectI(index int) (*Container, error) {
	return g.SetIndex(map[string]interface{}{}, index)
}

/*
Array - Create a new JSON array at a path. Returns an error if the path contains a collision with
a non object type.
*/
func (g *Container) Array(path ...string) (*Container, error) {
	return g.Set([]interface{}{}, path...)
}

/*
ArrayP - Does the same as Array, but using a dot notation JSON path.
*/
func (g *Container) ArrayP(path string) (*Container, error) {
	return g.Array(strings.Split(path, ".")...)
}

/*
ArrayI - Create a new JSON array at an array index. Returns an error if the object is not an array
or the index is out of bounds.
*/
func (g *Container) ArrayI(index int) (*Container, error) {
	return g.SetIndex([]interface{}{}, index)
}

/*
ArrayOfSize - Create a new JSON array of a particular size at a path. Returns an error if the path
contains a collision with a non object type.
*/
func (g *Container) ArrayOfSize(size int, path ...string) (*Container, error) {
	a := make([]interface{}, size)
	return g.Set(a, path...)
}

/*
ArrayOfSizeP - Does the same as ArrayOfSize, but using a dot notation JSON path.
*/
func (g *Container) ArrayOfSizeP(size int, path string) (*Container, error) {
	return g.ArrayOfSize(size, strings.Split(path, ".")...)
}

/*
ArrayOfSizeI - Create a new JSON array of a particular size at an array index. Returns an error if
the object is not an array or the index is out of bounds.
*/
func (g *Container) ArrayOfSizeI(size, index int) (*Container, error) {
	a := make([]interface{}, size)
	return g.SetIndex(a, index)
}

/*
Delete - Delete an element at a JSON path, an error is returned if the element does not exist.
*/
func (g *Container) Delete(path ...string) error {
	var object interface{}

	if g.object == nil {
		return ErrNotObj
	}
	object = g.object
	for target := 0; target < len(path); target++ {
		if mmap, ok := object.(map[string]interface{}); ok {
			if target == len(path)-1 {
				delete(mmap, path[target])
			} else if mmap[path[target]] == nil {
				return ErrNotObj
			}
			object = mmap[path[target]]
		} else {
			return ErrNotObj
		}
	}
	return nil
}

/*
DeleteP - Does the same as Delete, but using a dot notation JSON path.
*/
func (g *Container) DeleteP(path string) error {
	return g.Delete(strings.Split(path, ".")...)
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Array modification/search - Keeping these options simple right now, no need for anything more
complicated since you can just cast to []interface{}, modify and then reassign with Set.
*/

/*
ArrayAppend - Append a value onto a JSON array.
*/
func (g *Container) ArrayAppend(value interface{}, path ...string) error {
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return ErrNotArray
	}
	array = append(array, value)
	_, err := g.Set(array, path...)
	return err
}

/*
ArrayAppendP - Append a value onto a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayAppendP(value interface{}, path string) error {
	return g.ArrayAppend(value, strings.Split(path, ".")...)
}

/*
ArrayRemove - Remove an element from a JSON array.
*/
func (g *Container) ArrayRemove(index int, path ...string) error {
	if index < 0 {
		return ErrOutOfBounds
	}
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return ErrNotArray
	}
	if index < len(array) {
		array = append(array[:index], array[index+1:]...)
	} else {
		return ErrOutOfBounds
	}
	_, err := g.Set(array, path...)
	return err
}

/*
ArrayRemoveP - Remove an element from a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayRemoveP(index int, path string) error {
	return g.ArrayRemove(index, strings.Split(path, ".")...)
}

/*
ArrayElement - Access an element from a JSON array.
*/
func (g *Container) ArrayElement(index int, path ...string) (*Container, error) {
	if index < 0 {
		return &Container{nil}, ErrOutOfBounds
	}
	array, ok := g.Search(path...).Data().([]interface{})
	if !ok {
		return &Container{nil}, ErrNotArray
	}
	if index < len(array) {
		return &Container{array[index]}, nil
	}
	return &Container{nil}, ErrOutOfBounds
}

/*
ArrayElementP - Access an element from a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayElementP(index int, path string) (*Container, error) {
	return g.ArrayElement(index, strings.Split(path, ".")...)
}

/*
ArrayCount - Count the number of elements in a JSON array.
*/
func (g *Container) ArrayCount(path ...string) (int, error) {
	if array, ok := g.Search(path...).Data().([]interface{}); ok {
		return len(array), nil
	}
	return 0, ErrNotArray
}

/*
ArrayCountP - Count the number of elements in a JSON array using a dot notation JSON path.
*/
func (g *Container) ArrayCountP(path string) (int, error) {
	return g.ArrayCount(strings.Split(path, ".")...)
}

/*---------------------------------------------------------------------------------------------------
 */

/*
Bytes - Converts the contained object back to a JSON []byte blob.
*/
func (g *Container) Bytes() []byte {
	if g.object != nil {
		if bytes, err := json.Marshal(g.object); err == nil {
			return bytes
		}
	}
	return []byte("{}")
}

/*
BytesIndent - Converts the contained object back to a JSON []byte blob formatted with prefix and indent.
*/
func (g *Container) BytesIndent(prefix string, indent string) []byte {
	if g.object != nil {
		if bytes, err := json.MarshalIndent(g.object, prefix, indent); err == nil {
			return bytes
		}
	}
	return []byte("{}")
}

/*
String - Converts the contained object back to a JSON formatted string.
*/
func (g *Container) String() string {
	return string(g.Bytes())
}

/*
StringIndent - Converts the contained object back to a JSON formatted string with prefix and indent.
*/
func (g *Container) StringIndent(prefix string, indent string) string {
	return string(g.BytesIndent(prefix, indent))
}

/*
New - Create a new gabs JSON object.
*/
func New() *Container {
	return &Container{map[string]interface{}{}}
}

/*
Consume - Gobble up an already converted JSON object, or a fresh map[string]interface{} object.
*/
func Consume(root interface{}) (*Container, error) {
	return &Container{root}, nil
}

/*
ParseJSON - Convert a string into a representation of the parsed JSON.
*/
func ParseJSON(sample []byte) (*Container, error) {
	var gabs Container

	if err := json.Unmarshal(sample, &gabs.object); err != nil {
		return nil, err
	}

	return &gabs, nil
}

/*
ParseJSONDecoder - Convert a json.Decoder into a representation of the parsed JSON.
*/
func ParseJSONDecoder(decoder *json.Decoder) (*Container, error) {
	var gabs Container

	if err := decoder.Decode(&gabs.object); err != nil {
		return nil, err
	}

	return &gabs, nil
}

/*
ParseJSONFile - Read a file and convert into a representation of the parsed JSON.
*/
func ParseJSONFile(path string) (*Container, error) {
	if len(path) > 0 {
		cBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		container, err := ParseJSON(cBytes)
		if err != nil {
			return nil, err
		}

		return container, nil
	}
	return nil, ErrInvalidPath
}

/*
ParseJSONBuffer - Read the contents of a buffer into a representation of the parsed JSON.
*/
func ParseJSONBuffer(buffer io.Reader) (*Container, error) {
	var gabs Container
	jsonDecoder := json.NewDecoder(buffer)
	if err := jsonDecoder.Decode(&gabs.object); err != nil {
		return nil, err
	}

	return &gabs, nil
}

/*---------------------------------------------------------------------------------------------------
 */

// DEPRECATED METHODS

/*
Push - DEPRECATED: Push a value onto a JSON array.
*/
func (g *Container) Push(target string, value interface{}) error {
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			mmap[target] = append(array, value)
		} else {
			return ErrNotArray
		}
	} else {
		return ErrNotObj
	}
	return nil
}

/*
RemoveElement - DEPRECATED: Remove a value from a JSON array.
*/
func (g *Container) RemoveElement(target string, index int) error {
	if index < 0 {
		return ErrOutOfBounds
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			if index < len(array) {
				mmap[target] = append(array[:index], array[index+1:]...)
			} else {
				return ErrOutOfBounds
			}
		} else {
			return ErrNotArray
		}
	} else {
		return ErrNotObj
	}
	return nil
}

/*
GetElement - DEPRECATED: Get the desired element from a JSON array
*/
func (g *Container) GetElement(target string, index int) *Container {
	if index < 0 {
		return &Container{nil}
	}
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			if index < len(array) {
				return &Container{array[index]}
			}
		}
	}
	return &Container{nil}
}

/*
CountElements - DEPRECATED: Count the elements of a JSON array, returns -1 if the target is not an
array
*/
func (g *Container) CountElements(target string) int {
	if mmap, ok := g.Data().(map[string]interface{}); ok {
		arrayTarget := mmap[target]
		if array, ok := arrayTarget.([]interface{}); ok {
			return len(array)
		}
	}
	return -1
}

/*---------------------------------------------------------------------------------------------------
 */
