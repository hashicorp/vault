package gocb

import (
	"errors"
	"fmt"
)

// CouchbaseList represents a list document.
type CouchbaseList struct {
	collection *Collection
	id         string
}

// List returns a new CouchbaseList for the document specified by id.
func (c *Collection) List(id string) *CouchbaseList {
	return &CouchbaseList{
		collection: c,
		id:         id,
	}
}

// Iterator returns an iterable for all items in the list.
func (cl *CouchbaseList) Iterator() ([]interface{}, error) {
	span := cl.collection.startKvOpTrace("list_iterator", nil, false)
	defer span.End()

	return dsListIterator(span, cl.collection, cl.id)
}

func dsListIterator(span RequestSpan, collection *Collection, id string) ([]interface{}, error) {
	content, err := collection.Get(id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return nil, err
	}

	var listContents []interface{}
	err = content.Content(&listContents)
	if err != nil {
		return nil, err
	}

	return listContents, nil
}

// At retrieves the value specified at the given index from the list.
func (cl *CouchbaseList) At(index int, valuePtr interface{}) error {
	span := cl.collection.startKvOpTrace("list_at", nil, false)
	defer span.End()
	ops := make([]LookupInSpec, 1)
	ops[0] = GetSpec(fmt.Sprintf("[%d]", index), nil)
	result, err := cl.collection.LookupIn(cl.id, ops, &LookupInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return result.ContentAt(0, valuePtr)
}

// RemoveAt removes the value specified at the given index from the list.
func (cl *CouchbaseList) RemoveAt(index int) error {
	span := cl.collection.startKvOpTrace("list_remove_at", nil, false)
	defer span.End()
	ops := make([]MutateInSpec, 1)
	ops[0] = RemoveSpec(fmt.Sprintf("[%d]", index), nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return nil
}

// Append appends an item to the list.
func (cl *CouchbaseList) Append(val interface{}) error {
	span := cl.collection.startKvOpTrace("list_append", nil, false)
	defer span.End()
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayAppendSpec("", val, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{
		StoreSemantic: StoreSemanticsUpsert,
		ParentSpan:    span,
	})
	if err != nil {
		return err
	}

	return nil
}

// Prepend prepends an item to the list.
func (cl *CouchbaseList) Prepend(val interface{}) error {
	span := cl.collection.startKvOpTrace("list_prepend", nil, false)
	defer span.End()

	return dsListPrepend(span, cl.collection, cl.id, val)
}

func dsListPrepend(span RequestSpan, collection *Collection, id string, val interface{}) error {
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayPrependSpec("", val, nil)
	_, err := collection.MutateIn(id, ops, &MutateInOptions{
		StoreSemantic: StoreSemanticsUpsert,
		ParentSpan:    span,
	})
	if err != nil {
		return err
	}

	return nil
}

// IndexOf gets the index of the item in the list.
func (cl *CouchbaseList) IndexOf(val interface{}) (int, error) {
	span := cl.collection.startKvOpTrace("list_index_of", nil, false)
	defer span.End()
	content, err := cl.collection.Get(cl.id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return 0, err
	}

	var listContents []interface{}
	err = content.Content(&listContents)
	if err != nil {
		return 0, err
	}

	for i, item := range listContents {
		if item == val {
			return i, nil
		}
	}

	return -1, nil
}

// Size returns the size of the list.
func (cl *CouchbaseList) Size() (int, error) {
	span := cl.collection.startKvOpTrace("list_size", nil, false)
	defer span.End()

	return dsListSize(span, cl.collection, cl.id)
}

func dsListSize(span RequestSpan, collection *Collection, id string) (int, error) {
	ops := make([]LookupInSpec, 1)
	ops[0] = CountSpec("", nil)
	result, err := collection.LookupIn(id, ops, &LookupInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return 0, err
	}

	var count int
	err = result.ContentAt(0, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Clear clears a list, also removing it.
func (cl *CouchbaseList) Clear() error {
	span := cl.collection.startKvOpTrace("list_clear", nil, false)
	defer span.End()

	return dsListClear(span, cl.collection, cl.id)
}

func dsListClear(span RequestSpan, collection *Collection, id string) error {
	_, err := collection.Remove(id, &RemoveOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return nil
}

// CouchbaseMap represents a map document.
type CouchbaseMap struct {
	collection *Collection
	id         string
}

// Map returns a new CouchbaseMap.
func (c *Collection) Map(id string) *CouchbaseMap {
	return &CouchbaseMap{
		collection: c,
		id:         id,
	}
}

// Iterator returns an iterable for all items in the map.
func (cl *CouchbaseMap) Iterator() (map[string]interface{}, error) {
	span := cl.collection.startKvOpTrace("map_iterator", nil, false)
	defer span.End()
	content, err := cl.collection.Get(cl.id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	return mapContents, nil
}

// At retrieves the item for the given id from the map.
func (cl *CouchbaseMap) At(id string, valuePtr interface{}) error {
	span := cl.collection.startKvOpTrace("map_at", nil, false)
	defer span.End()
	ops := make([]LookupInSpec, 1)
	ops[0] = GetSpec(id, nil)
	result, err := cl.collection.LookupIn(cl.id, ops, &LookupInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return result.ContentAt(0, valuePtr)
}

// Add adds an item to the map.
func (cl *CouchbaseMap) Add(id string, val interface{}) error {
	span := cl.collection.startKvOpTrace("map_add", nil, false)
	defer span.End()
	ops := make([]MutateInSpec, 1)
	ops[0] = UpsertSpec(id, val, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{
		StoreSemantic: StoreSemanticsUpsert,
		ParentSpan:    span,
	})
	if err != nil {
		return err
	}

	return nil
}

// Remove removes an item from the map.
func (cl *CouchbaseMap) Remove(id string) error {
	span := cl.collection.startKvOpTrace("map_remove", nil, false)
	defer span.End()
	ops := make([]MutateInSpec, 1)
	ops[0] = RemoveSpec(id, nil)
	_, err := cl.collection.MutateIn(cl.id, ops, &MutateInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return nil
}

// Exists verifies whether or a id exists in the map.
func (cl *CouchbaseMap) Exists(id string) (bool, error) {
	span := cl.collection.startKvOpTrace("map_exists", nil, false)
	defer span.End()
	ops := make([]LookupInSpec, 1)
	ops[0] = ExistsSpec(id, nil)
	result, err := cl.collection.LookupIn(cl.id, ops, nil)
	if err != nil {
		return false, err
	}

	return result.Exists(0), nil
}

// Size returns the size of the map.
func (cl *CouchbaseMap) Size() (int, error) {
	span := cl.collection.startKvOpTrace("map_size", nil, false)
	defer span.End()
	ops := make([]LookupInSpec, 1)
	ops[0] = CountSpec("", nil)
	result, err := cl.collection.LookupIn(cl.id, ops, &LookupInOptions{
		ParentSpan: span,
	})
	if err != nil {
		return 0, err
	}

	var count int
	err = result.ContentAt(0, &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Keys returns all of the keys within the map.
func (cl *CouchbaseMap) Keys() ([]string, error) {
	span := cl.collection.startKvOpTrace("map_keys", nil, false)
	defer span.End()
	content, err := cl.collection.Get(cl.id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	var keys []string
	for id := range mapContents {
		keys = append(keys, id)
	}

	return keys, nil
}

// Values returns all of the values within the map.
func (cl *CouchbaseMap) Values() ([]interface{}, error) {
	span := cl.collection.startKvOpTrace("map_values", nil, false)
	defer span.End()
	content, err := cl.collection.Get(cl.id, nil)
	if err != nil {
		return nil, err
	}

	var mapContents map[string]interface{}
	err = content.Content(&mapContents)
	if err != nil {
		return nil, err
	}

	var values []interface{}
	for _, val := range mapContents {
		values = append(values, val)
	}

	return values, nil
}

// Clear clears a map, also removing it.
func (cl *CouchbaseMap) Clear() error {
	span := cl.collection.startKvOpTrace("map_clear", nil, false)
	defer span.End()
	_, err := cl.collection.Remove(cl.id, &RemoveOptions{
		ParentSpan: span,
	})
	if err != nil {
		return err
	}

	return nil
}

// CouchbaseSet represents a set document.
type CouchbaseSet struct {
	id         string
	collection *Collection
}

// Set returns a new CouchbaseSet.
func (c *Collection) Set(id string) *CouchbaseSet {
	return &CouchbaseSet{
		id:         id,
		collection: c,
	}
}

// Iterator returns an iterable for all items in the set.
func (cs *CouchbaseSet) Iterator() ([]interface{}, error) {
	span := cs.collection.startKvOpTrace("set_iterator", nil, false)
	defer span.End()
	return dsListIterator(span, cs.collection, cs.id)
}

// Add adds a value to the set.
func (cs *CouchbaseSet) Add(val interface{}) error {
	span := cs.collection.startKvOpTrace("set_add", nil, false)
	defer span.End()
	ops := make([]MutateInSpec, 1)
	ops[0] = ArrayAddUniqueSpec("", val, nil)
	_, err := cs.collection.MutateIn(cs.id, ops, &MutateInOptions{
		StoreSemantic: StoreSemanticsUpsert,
		ParentSpan:    span,
	})
	if err != nil {
		return err
	}

	return nil
}

// Remove removes an value from the set.
func (cs *CouchbaseSet) Remove(val string) error {
	span := cs.collection.startKvOpTrace("set_remove", nil, false)
	defer span.End()
	for i := 0; i < 16; i++ {
		content, err := cs.collection.Get(cs.id, &GetOptions{
			ParentSpan: span,
		})
		if err != nil {
			return err
		}

		cas := content.Cas()

		var setContents []interface{}
		err = content.Content(&setContents)
		if err != nil {
			return err
		}

		indexToRemove := -1
		for i, item := range setContents {
			if item == val {
				indexToRemove = i
			}
		}

		if indexToRemove > -1 {
			ops := make([]MutateInSpec, 1)
			ops[0] = RemoveSpec(fmt.Sprintf("[%d]", indexToRemove), nil)
			_, err = cs.collection.MutateIn(cs.id, ops, &MutateInOptions{
				Cas:        cas,
				ParentSpan: span,
			})
			if errors.Is(err, ErrCasMismatch) {
				continue
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("failed to perform operation after 16 retries")
}

// Values returns all of the values within the set.
func (cs *CouchbaseSet) Values() ([]interface{}, error) {
	span := cs.collection.startKvOpTrace("set_values", nil, false)
	defer span.End()
	content, err := cs.collection.Get(cs.id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return nil, err
	}

	var setContents []interface{}
	err = content.Content(&setContents)
	if err != nil {
		return nil, err
	}

	return setContents, nil
}

// Contains verifies whether or not a value exists within the set.
func (cs *CouchbaseSet) Contains(val string) (bool, error) {
	span := cs.collection.startKvOpTrace("set_contains", nil, false)
	defer span.End()
	content, err := cs.collection.Get(cs.id, &GetOptions{
		ParentSpan: span,
	})
	if err != nil {
		return false, err
	}

	var setContents []interface{}
	err = content.Content(&setContents)
	if err != nil {
		return false, err
	}

	for _, item := range setContents {
		if item == val {
			return true, nil
		}
	}

	return false, nil
}

// Size returns the size of the set
func (cs *CouchbaseSet) Size() (int, error) {
	span := cs.collection.startKvOpTrace("set_size", nil, false)
	defer span.End()
	return dsListSize(span, cs.collection, cs.id)
}

// Clear clears a set, also removing it.
func (cs *CouchbaseSet) Clear() error {
	span := cs.collection.startKvOpTrace("set_clear", nil, false)
	defer span.End()
	return dsListClear(span, cs.collection, cs.id)
}

// CouchbaseQueue represents a queue document.
type CouchbaseQueue struct {
	id         string
	collection *Collection
}

// Queue returns a new CouchbaseQueue.
func (c *Collection) Queue(id string) *CouchbaseQueue {
	return &CouchbaseQueue{
		id:         id,
		collection: c,
	}
}

// Iterator returns an iterable for all items in the queue.
func (cs *CouchbaseQueue) Iterator() ([]interface{}, error) {
	span := cs.collection.startKvOpTrace("queue_iterator", nil, false)
	defer span.End()
	return dsListIterator(span, cs.collection, cs.id)
}

// Push pushes a value onto the queue.
func (cs *CouchbaseQueue) Push(val interface{}) error {
	span := cs.collection.startKvOpTrace("queue_push", nil, false)
	defer span.End()
	return dsListPrepend(span, cs.collection, cs.id, val)
}

// Pop pops an items off of the queue.
func (cs *CouchbaseQueue) Pop(valuePtr interface{}) error {
	span := cs.collection.startKvOpTrace("queue_pop", nil, false)
	defer span.End()
	for i := 0; i < 16; i++ {
		ops := make([]LookupInSpec, 1)
		ops[0] = GetSpec("[-1]", nil)
		content, err := cs.collection.LookupIn(cs.id, ops, &LookupInOptions{
			ParentSpan: span,
		})
		if err != nil {
			return err
		}

		cas := content.Cas()
		err = content.ContentAt(0, valuePtr)
		if err != nil {
			return err
		}

		mutateOps := make([]MutateInSpec, 1)
		mutateOps[0] = RemoveSpec("[-1]", nil)
		_, err = cs.collection.MutateIn(cs.id, mutateOps, &MutateInOptions{
			Cas:        cas,
			ParentSpan: span,
		})
		if errors.Is(err, ErrCasMismatch) {
			continue
		}
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("failed to perform operation after 16 retries")
}

// Size returns the size of the queue.
func (cs *CouchbaseQueue) Size() (int, error) {
	span := cs.collection.startKvOpTrace("queue_size", nil, false)
	defer span.End()
	return dsListSize(span, cs.collection, cs.id)
}

// Clear clears a queue, also removing it.
func (cs *CouchbaseQueue) Clear() error {
	span := cs.collection.startKvOpTrace("queue_clear", nil, false)
	defer span.End()
	return dsListClear(span, cs.collection, cs.id)
}
