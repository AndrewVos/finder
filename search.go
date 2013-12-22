package main

import (
	"sort"
	"strings"
)

var currentID = 0
var Mappings map[string]FieldMapping
var allSearchables map[int]*Searchable
var wordIndexes map[string]*WordIndex

type FieldMapping struct {
	Type     string
	Sortable bool
}

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

type Sort struct {
	Field     string
	Ascending bool
}

type Query struct {
	Text []TextQuery
	Sort []Sort
}

type Searchables []*Searchable

func (s BySort) Len() int { return len(s.Searchables) }
func (s BySort) Swap(i, j int) {
	s.Searchables[i], s.Searchables[j] = s.Searchables[j], s.Searchables[i]
}

type BySort struct {
	Searchables Searchables
	Sort        []Sort
}

func (s BySort) Less(i, j int) bool {
	field := s.Sort[0].Field
	a, b := s.Searchables[i].Source[field].(string), s.Searchables[j].Source[field].(string)
	if s.Sort[0].Ascending {
		return a < b
	} else {
		return a > b
	}
}

func Search(query Query) Searchables {
	for _, text := range query.Text {
		results := searchTextField(text.Field, text.Value)
		if len(query.Sort) > 0 {
			sorter := BySort{results, query.Sort}
			sort.Sort(sorter)
		}
		return results
	}
	return nil
}

func get(id int) *Searchable {
	return allSearchables[id]
}

func searchTextField(field string, query string) Searchables {
	var things Searchables
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
		if mapping.Type == "text" {
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
