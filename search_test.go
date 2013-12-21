package main

import (
	"testing"
)

func TestFindsSimpleMatches(t *testing.T) {
	Mappings = map[string]string{
		"name": "text",
	}
	Index(map[string]interface{}{"name": "some thing"})
	Index(map[string]interface{}{"name": "some other thing"})
	Index(map[string]interface{}{"name": "other"})

	results := Search([]QueryPart{{"name", "thing"}})
	if len(results) != 2 {
		t.Errorf("Expected 2 results")
	}
}

func TestFindsMultipleWordsInQuery(t *testing.T) {
	Mappings = map[string]string{
		"name": "text",
	}
	Index(map[string]interface{}{"name": "batman spiderman superman"})
	Index(map[string]interface{}{"name": "spiderman"})
	Index(map[string]interface{}{"name": "spiderman superman"})

	results := Search([]QueryPart{{"name", "spiderman superman"}})
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, but got %d\n", len(results))
	}

	expectedNames := []string{"batman spiderman superman", "spiderman superman"}
	for i, expectedName := range expectedNames {
		if actual := results[i]["name"].(string); actual != expectedName {
			t.Errorf("Expected first element to be %q, but was %q\n", expectedName, actual)
		}
	}
}
