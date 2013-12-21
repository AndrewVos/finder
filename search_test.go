package main

import (
	"testing"
)

func expectAmountOfResults(t *testing.T, results []map[string]interface{}, expected int) {
	if len(results) != expected {
		t.Fatalf("Expected %d results, but got %d\n", expected, len(results))
	}
}

func expectedThingWithName(t *testing.T, results []map[string]interface{}, index int, expectedName string) {
	if actual := results[index]["name"].(string); actual != expectedName {
		t.Errorf("Expected first element to be %q, but was %q\n", expectedName, actual)
	}
}

func createNameMapping() {
	Mappings = map[string]string{"name": "text"}
}

func indexProductWithName(name string) {
	Index(map[string]interface{}{"name": name})
}

func TestFindsSimpleMatches(t *testing.T) {
	createNameMapping()
	indexProductWithName("some thing")
	indexProductWithName("some other thing")
	indexProductWithName("other")

	results := Search([]QueryPart{{"name", "thing"}})
	expectAmountOfResults(t, results, 2)
}

func TestFindsMultipleWordsInQuery(t *testing.T) {
	createNameMapping()
	indexProductWithName("batman spiderman superman")
	indexProductWithName("spiderman")
	indexProductWithName("spiderman superman")

	results := Search([]QueryPart{{"name", "spiderman superman"}})
	expectAmountOfResults(t, results, 2)

	expectedThingWithName(t, results, 0, "batman spiderman superman")
	expectedThingWithName(t, results, 1, "spiderman superman")
}
