package main

import (
	"strings"
)

var stems map[string][]string
var currentID = 0
var Mappings map[string]string
var tries map[string]*Trie
var allThings map[int]map[string]interface{}

type TextQuery struct {
	Field string
	Value string
}

type Query struct {
	Text []TextQuery
}

func Search(query Query) []map[string]interface{} {
	for _, text := range query.Text {
		return searchTextField(text.Field, text.Value)
	}
	return nil
}

func get(id int) map[string]interface{} {
	return allThings[id]
}

func searchTextField(name string, query string) []map[string]interface{} {
	trie := tries[name]
	var things []map[string]interface{}
	foundIDsCount := map[int]int{}

	words := strings.Split(query, " ")
	for _, word := range words {
		node, found := trie.Find([]byte(word))
		if found == true {
			for _, id := range node.IDs {
				foundIDsCount[id] += 1
			}
		}
	}

	expectedMatches := len(words)
	for id, count := range foundIDsCount {
		if count == expectedMatches {
			things = append(things, get(id))
		}
	}
	return things
}

func getNextId() int {
	i := currentID
	currentID += 1
	return i
}

func Index(thing map[string]interface{}) {
	id := getNextId()
	if allThings == nil {
		allThings = map[int]map[string]interface{}{}
	}
	allThings[id] = thing

	for name, mapping := range Mappings {
		if mapping == "text" {
			indexTextField(id, name, thing[name].(string))
		}
	}
}

func indexTextField(id int, name string, value string) {
	if tries == nil {
		tries = map[string]*Trie{}
	}
	if _, ok := tries[name]; !ok {
		tries[name] = &Trie{}
	}
	trie := tries[name]

	value = strings.ToLower(value)
	words := strings.Split(value, " ")
	for _, word := range words {
		if word != "" {
			trie.Add(id, []byte(word))
		}
	}
}
