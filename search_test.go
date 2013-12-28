package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"
)

func expectAmountOfResults(t *testing.T, results []*Document, expected int) {
	if len(results) != expected {
		t.Fatalf("Expected %d results, but got %d\n", expected, len(results))
	}
}

func expectDocumentWithName(t *testing.T, results []*Document, index int, expectedName string) {
	if actual := results[index].Source["name"].(string); actual != expectedName {
		t.Errorf("Expected element %d to be %q, but was %q\n", index, expectedName, actual)
	}
}

func cleanup() {
	currentID = 0
	Mappings = nil
	allDocuments = nil
	allIndexes = nil
}

func createNameMapping() {
	Mappings = map[string]string{
		"name": "text",
	}
}

func indexProductWithName(name string) {
	Index(map[string]interface{}{"name": name})
}

func indexProductWithPopularity(name string, popularity int) {
	Index(map[string]interface{}{"name": name, "popularity": popularity})
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

	expectDocumentWithName(t, results, 0, "batman spiderman superman")
	expectDocumentWithName(t, results, 1, "spiderman superman")
}

func TestSortTextAscending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]string{
		"name": "text",
	}
	indexProductWithName("c thing")
	indexProductWithName("a thing")
	indexProductWithName("z thing")

	query := Query{
		Text: []TextQuery{{"name", "thing"}},
		Sort: []Sort{{Field: "name", Ascending: true}},
	}

	results := Search(query)
	expectDocumentWithName(t, results, 0, "a thing")
	expectDocumentWithName(t, results, 1, "c thing")
	expectDocumentWithName(t, results, 2, "z thing")
}

func TestSortTextDescending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]string{"name": "text"}
	indexProductWithName("c thing")
	indexProductWithName("a thing")
	indexProductWithName("z thing")

	query := Query{
		Text: []TextQuery{{"name", "thing"}},
		Sort: []Sort{{Field: "name", Ascending: false}},
	}

	results := Search(query)
	expectDocumentWithName(t, results, 0, "z thing")
	expectDocumentWithName(t, results, 1, "c thing")
	expectDocumentWithName(t, results, 2, "a thing")
}

func TestSortIntAscending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]string{"name": "text", "popularity": "integer"}
	indexProductWithPopularity("name 1", 10)
	indexProductWithPopularity("name 2", 5)
	indexProductWithPopularity("name 3", 1)

	query := Query{
		Text: []TextQuery{{"name", "name"}},
		Sort: []Sort{{Field: "popularity", Ascending: true}},
	}

	results := Search(query)
	expectDocumentWithName(t, results, 0, "name 3")
	expectDocumentWithName(t, results, 1, "name 2")
	expectDocumentWithName(t, results, 2, "name 1")
}
func TestSortIntDescending(t *testing.T) {
	defer cleanup()

	Mappings = map[string]string{"name": "text", "popularity": "integer"}
	indexProductWithPopularity("name 1", 10)
	indexProductWithPopularity("name 2", 5)
	indexProductWithPopularity("name 3", 1)

	query := Query{
		Text: []TextQuery{{"name", "name"}},
		Sort: []Sort{{Field: "popularity", Ascending: false}},
	}

	results := Search(query)
	expectDocumentWithName(t, results, 0, "name 1")
	expectDocumentWithName(t, results, 1, "name 2")
	expectDocumentWithName(t, results, 2, "name 3")
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
		if fileCount == 10 {
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

	texts := []string{
		"blue dress",
		"car",
		"monkey",
		"dog house",
	}

	for _, text := range texts {
		fmt.Printf("Searching for %q\n", text)

		start := time.Now()
		query := Query{Text: []TextQuery{{"name", text}}, Sort: []Sort{{Field: "name", Ascending: true}}}
		results := Search(query)
		elapsed := time.Since(start)

		fmt.Printf("Found %d things out of a total of %d\n", len(results), thingCount)
		fmt.Printf("Search took %s\n", elapsed)
		if len(results) > 10 {
			results = results[:10]
		}
		for _, document := range results {
			fmt.Println(document.Source["name"])
		}
		fmt.Println()
	}

	printMemoryStats := func() {
		stats := &runtime.MemStats{}
		runtime.ReadMemStats(stats)
		fmt.Printf("Current memory usage: %d\n", stats.Alloc)
	}
	printMemoryStats()
}
