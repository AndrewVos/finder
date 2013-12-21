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

func createTextQuery(field string, value string) Query {
	return Query{Text: []TextQuery{{field, value}}}
}

func TestFindsSimpleMatches(t *testing.T) {
	createNameMapping()
	indexProductWithName("some  thing")
	indexProductWithName("some other thing")
	indexProductWithName("other")

	results := Search(createTextQuery("name", "thing"))
	expectAmountOfResults(t, results, 2)
}

func TestFindsMultipleWordsInQuery(t *testing.T) {
	createNameMapping()
	indexProductWithName("batman spiderman superman")
	indexProductWithName("spiderman")
	indexProductWithName("spiderman superman")

	results := Search(createTextQuery("name", "spiderman superman"))
	expectAmountOfResults(t, results, 2)

	expectedThingWithName(t, results, 0, "batman spiderman superman")
	expectedThingWithName(t, results, 1, "spiderman superman")
}

func TestLargeFile(t *testing.T) {
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

	start := time.Now()
	results := Search(createTextQuery("name", "blue dress"))
	elapsed := time.Since(start)

	for _, result := range results {
		fmt.Println(result["name"])
	}

	fmt.Println()
	fmt.Printf("Found %d things out of a total of %d\n", len(results), thingCount)

	log.Printf("Search took %s", elapsed)
}
