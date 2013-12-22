package main

import (
	"strings"
)

var currentID = 0
var Mappings map[string]string
var allSearchables map[int]*Searchable
var wordIndexes map[string]*WordIndex

type Searchable struct {
	ID     int
	Source map[string]interface{}
}

type WordIndex struct {
	Words map[string][]WordCount
}

type WordCount struct {
	Searchable *Searchable
	Count      int
}

type TextQuery struct {
	Field string
	Value string
}

type Query struct {
	Text []TextQuery
}

func Search(query Query) []*Searchable {
	for _, text := range query.Text {
		return searchTextField(text.Field, text.Value)
	}
	return nil
}

func get(id int) *Searchable {
	return allSearchables[id]
}

func searchTextField(field string, query string) []*Searchable {
	var things []*Searchable
	words := strings.Split(query, " ")
	wordIndex := wordIndexes[field]
	matchCount := map[*Searchable]int{}
	for _, word := range words {
		for _, wordCount := range wordIndex.Words[word] {
			matchCount[wordCount.Searchable] += 1
		}
	}
	for thing, count := range matchCount {
		if count == len(words) {
			things = append(things, thing)
		}
	}
	return things
}

func getNextId() int {
	i := currentID
	currentID += 1
	return i
}

func Index(source map[string]interface{}) {
	id := getNextId()
	if allSearchables == nil {
		allSearchables = map[int]*Searchable{}
	}
	thing := &Searchable{ID: id, Source: source}
	allSearchables[id] = thing

	for field, mapping := range Mappings {
		if mapping == "text" {
			indexTextField(thing, field)
		}
	}
}

func indexTextField(thing *Searchable, field string) {
	if wordIndexes == nil {
		wordIndexes = map[string]*WordIndex{}
	}
	wordIndex := wordIndexes[field]
	if wordIndex == nil {
		wordIndex = &WordIndex{Words: map[string][]WordCount{}}
		wordIndexes[field] = wordIndex
	}

	value := strings.ToLower(thing.Source[field].(string))
	words := strings.Split(value, " ")
	wordCount := map[string]int{}
	for _, word := range words {
		if word != "" {
			wordCount[word] += 1
		}
	}

	for word, count := range wordCount {
		wc := WordCount{
			Searchable: thing,
			Count:      count,
		}
		wordIndex.Words[word] = append(wordIndex.Words[word], wc)
	}
}
