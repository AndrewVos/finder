package main

import (
	"sort"
	"strings"
)

var currentID = 0
var Mappings map[string]FieldMapping
var allSearchables map[int]*Searchable

type WordNode struct {
	Document *Searchable
	Child    *WordNode
}

var wordNodes map[string]*WordNode
var lastWordNodes map[string]*WordNode

type FieldMapping struct {
	Type     string
	Sortable bool
}

type Searchable struct {
	ID     int
	Source map[string]interface{}
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
	var documents Searchables

	words := splitWords(query)
	requiredMatches := len(words)
	matches := map[*Searchable]int{}

	for _, word := range words {
		if node, ok := wordNodes[word]; ok {
			for {
				matches[node.Document] += 1
				if matches[node.Document] == requiredMatches {
					documents = append(documents, node.Document)
				}
				if node.Child == nil {
					break
				}
				node = node.Child
			}
		} else {
			return documents
		}
	}

	return documents
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

func splitWords(s string) []string {
	s = strings.ToLower(s)
	var words []string
	for _, word := range strings.Split(s, " ") {
		if word != "" {
			words = append(words, word)
		}
	}

	return words
}

func indexTextField(document *Searchable, field string) {
	if wordNodes == nil {
		wordNodes = map[string]*WordNode{}
	}
	if lastWordNodes == nil {
		lastWordNodes = map[string]*WordNode{}
	}

	value := document.Source[field].(string)
	words := splitWords(value)
	for _, word := range words {
		node, exists := wordNodes[word]
		if exists {
			lastWordNode := lastWordNodes[word]
			lastWordNode.Child = &WordNode{Document: document}
			lastWordNodes[word] = lastWordNode.Child
		} else {
			node = &WordNode{Document: document}
			wordNodes[word] = node
			lastWordNodes[word] = node
		}
	}
}
