package simplekvstore

import (
	"reflect"
	"sort"
	"testing"
)

func assert(expected interface{}, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v"+"\n"+"Actual: %v", expected, actual)
	}
}

func TestSet(t *testing.T) {
	store := New()
	key, val := "foo", "bar"

	store.Set(key, val)

	assert(val, store.m[key].Value, t)
}

func TestGet(t *testing.T) {
	store := New()
	key, val := "foo", "bar"

	store.Set(key, val)

	// it should get a value from store
	kv, _ := store.Get(key)
	assert(kv.Value, val, t)

	// it should return an error if the key does not exist
	_, err := store.Get("baz")
	assert(KeyError{"baz", ErrNotFound}, err, t)
}

func TestGetAll(t *testing.T) {
	store := New()

	// it should return empty list of array if no valus in store
	valuesEmpty := store.GetAll()

	assert(make([]KeyValue, 0), valuesEmpty, t)

	// it should return array with all KeyValue
	key1, val1 := "foo", "bar"
	key2, val2 := "baz", "doe"

	store.Set(key1, val1)
	store.Set(key2, val2)

	values := store.GetAll()

	sort.Slice(values, func(i, j int) bool {
		return values[i].Key < values[j].Key
	})

	assert(
		[]KeyValue{KeyValue{key2, val2}, KeyValue{key1, val1}},
		values,
		t,
	)
}

func TestGetClone(t *testing.T) {
	store := New()

	// setting a value directly will change it in store
	store.m["foo"] = KeyValue{"foo", "baz"}
	got, _ := store.Get("foo")
	assert(KeyValue{"foo", "baz"}, got, t)

	// updating a clone should not update the original
	storeClone := store.GetClone()
	storeClone["foo"] = KeyValue{"bar", "bar"}

	got, _ = store.Get("foo")
	assert(KeyValue{"foo", "baz"}, got, t)
}

func TestGetWithPrefix(t *testing.T) {
	store := New()
	valuesList := map[string]string{
		"pre-foo": "1",
		"pre-bar": "2",
		"foo-bar": "4",
		"bazbar":  "5",
	}

	for key, val := range valuesList {
		store.Set(key, val)
	}

	res := store.GetWithPrefix("pre")
	sort.Slice(res, func(i, j int) bool {
		return res[i].Value < res[j].Value
	})

	assert(
		[]KeyValue{KeyValue{"pre-foo", "1"}, KeyValue{"pre-bar", "2"}},
		res,
		t,
	)
}

func TestGetKeysContaining(t *testing.T) {
	store := New()
	valuesList := map[string]string{
		"megatron": "1",
		"lol":      "5",
	}

	for key, val := range valuesList {
		store.Set(key, val)
	}

	assert(
		[]KeyValue{KeyValue{"megatron", "1"}},
		store.GetKeysContaining("tr"),
		t,
	)
}

func TestGetValuesContaining(t *testing.T) {
	store := New()
	valuesList := map[string]string{
		"1": "foo_bar",
		"2": "s_foo_b",
	}

	for key, val := range valuesList {
		store.Set(key, val)
	}

	res := store.GetValuesContaining("foo")
	sort.Slice(res, func(i, j int) bool {
		return res[i].Key < res[j].Key
	})

	assert(
		[]KeyValue{KeyValue{"1", "foo_bar"}, KeyValue{"2", "s_foo_b"}},
		res,
		t,
	)
}
