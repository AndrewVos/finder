package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func expectAmountOfResults(t *testing.T, results []*Searchable, expected int) {
	if len(results) != expected {
		t.Fatalf("Expected %d results, but got %d\n", expected, len(results))
	}
}

func expectedSearchableWithName(t *testing.T, results []*Searchable, index int, expectedName string) {
	if actual := results[index].Source["name"].(string); actual != expectedName {
		t.Errorf("Expected element %d to be %q, but was %q\n", index, expectedName, actual)
	}
}

func cleanup() {
	currentID = 0
	Mappings = nil
	allSearchables = nil
	wordNodes = nil
	lastWordNodes = nil
}

func createNameMapping() {
	Mappings = map[string]FieldMapping{
		"name": FieldMapping{Type: "text", Sortable: false},
	}
}

func indexProductWithName(name string) {
	Index(map[string]interface{}{"name": name})
}

func createTextQuery(field string, value string) Query {
	return Query{Text: []TextQuery{{field, value}}}
}

func TestFindsSimpleMatches(t *testing.T) {
	defer cleanup()

	createNameMapping()
	indexProductWithName("some  thing")
	indexProductWithName("some other thing")
	indexProductWithName("other")

	results := Search(createTextQuery("name", "thing"))
	expectAmountOfResults(t, results, 2)
}

func TestFindsMultipleWordsInQuery(t *testing.T) {
	defer cleanup()

	createNameMapping()
	indexProductWithName("batman spiderman superman")
	indexProductWithName("spiderman")
	indexProductWithName("spiderman superman")

	results := Search(createTextQuery("name", "spiderman superman"))
	expectAmountOfResults(t, results, 2)

	expectedSearchableWithName(t, results, 0, "batman spiderman superman")
	expectedSearchableWithName(t, results, 1, "spiderman superman")
}

func TestResultsAreSortedAscending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]FieldMapping{
		"name": FieldMapping{Type: "text"},
	}
	indexProductWithName("c thing")
	indexProductWithName("a thing")
	indexProductWithName("z thing")

	query := Query{
		Text: []TextQuery{{"name", "thing"}},
		Sort: []Sort{{Field: "name", Ascending: true}},
	}

	results := Search(query)
	expectedSearchableWithName(t, results, 0, "a thing")
	expectedSearchableWithName(t, results, 1, "c thing")
	expectedSearchableWithName(t, results, 2, "z thing")
}

func TestResultsAreSortedDescending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]FieldMapping{
		"name": FieldMapping{Type: "text"},
	}
	indexProductWithName("c thing")
	indexProductWithName("a thing")
	indexProductWithName("z thing")

	query := Query{
		Text: []TextQuery{{"name", "thing"}},
		Sort: []Sort{{Field: "name", Ascending: false}},
	}

	results := Search(query)
	expectedSearchableWithName(t, results, 0, "z thing")
	expectedSearchableWithName(t, results, 1, "c thing")
	expectedSearchableWithName(t, results, 2, "a thing")
}

func TestLargeFile(t *testing.T) {
	defer cleanup()

	jsonPath := os.Getenv("JSON_PATH")
	if jsonPath == "" {
		return
	}

	createNameMapping()
	files, err := ioutil.ReadDir(jsonPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileCount := 0
	thingCount := 0
	for _, file := range files {
		if fileCount == 100 {
			break
		}
		fileCount += 1
		fmt.Println("indexing " + file.Name())
		b, _ := ioutil.ReadFile(jsonPath + "/" + file.Name())
		var things []map[string]interface{}
		err = json.Unmarshal(b, &things)
		if err != nil {
			fmt.Println(err)
			continue
		}
		thingCount += len(things)
		for _, thing := range things {
			Index(thing)
		}
	}
	fmt.Println("indexing complete")

	queries := []string{
		"blue dress",
		"car",
		"monkey",
		"dog house",
	}

	for _, query := range queries {
		log.Printf("Searching for %q", query)

		start := time.Now()
		results := Search(createTextQuery("name", query))
		elapsed := time.Since(start)

		log.Printf("Found %d things out of a total of %d\n", len(results), thingCount)
		log.Printf("Search took %s", elapsed)
		log.Println()
	}
}
