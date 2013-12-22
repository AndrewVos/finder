package main

import (
	"strings"
)

var stems map[string][]string
var currentID = 0
var Mappings map[string]string
var tries map[string]*Trie
var allThings map[int]*Thing

var wordIndexes map[string]*WordIndex

type Thing struct {
	ID     int
	Source map[string]interface{}
}

type TextQuery struct {
	Field string
	Value string
}

type Query struct {
	Text []TextQuery
}

func Search(query Query) []*Thing {
	for _, text := range query.Text {
		return searchTextField(text.Field, text.Value)
	}
	return nil
}

func get(id int) *Thing {
	return allThings[id]
}

func searchTextField(field string, query string) []*Thing {
	var things []*Thing
	words := strings.Split(query, " ")
	wordIndex := wordIndexes[field]
	matchCount := map[*Thing]int{}
	for _, word := range words {
		for _, wordCount := range wordIndex.Words[word] {
			matchCount[wordCount.Thing] += 1
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
	if allThings == nil {
		allThings = map[int]*Thing{}
	}
	thing := &Thing{ID: id, Source: source}
	allThings[id] = thing

	for field, mapping := range Mappings {
		if mapping == "text" {
			indexTextField(thing, field)
		}
	}
}

type WordIndex struct {
	Words map[string][]WordCount
}

func indexTextField(thing *Thing, field string) {
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
			Thing: thing,
			Count: count,
		}
		wordIndex.Words[word] = append(wordIndex.Words[word], wc)
	}
}

type WordCount struct {
	Thing *Thing
	Count int
}
