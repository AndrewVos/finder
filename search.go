package main

import (
	"sort"
	"strings"
)

var currentID = 0
var Mappings map[string]FieldMapping
var allDocuments map[int]*Document

type WordNode struct {
	Document *Document
	Child    *WordNode
}

var wordNodes map[string]*WordNode
var lastWordNodes map[string]*WordNode

type FieldMapping struct {
	Type     string
	Sortable bool
}

type Document struct {
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

type Documents []*Document

func (s BySort) Len() int { return len(s.Documents) }
func (s BySort) Swap(i, j int) {
	s.Documents[i], s.Documents[j] = s.Documents[j], s.Documents[i]
}

type BySort struct {
	Documents Documents
	Sort      []Sort
}

func (s BySort) Less(i, j int) bool {
	field := s.Sort[0].Field
	a, b := s.Documents[i].Source[field].(string), s.Documents[j].Source[field].(string)
	if s.Sort[0].Ascending {
		return a < b
	} else {
		return a > b
	}
}

func Search(query Query) Documents {
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

func get(id int) *Document {
	return allDocuments[id]
}

func searchTextField(field string, query string) Documents {
	var documents Documents

	words := splitWords(query)
	requiredMatches := len(words)
	matches := map[*Document]int{}

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
	if allDocuments == nil {
		allDocuments = map[int]*Document{}
	}
	document := &Document{ID: id, Source: source}
	allDocuments[id] = document

	for field, mapping := range Mappings {
		if mapping.Type == "text" {
			indexTextField(document, field)
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

func indexTextField(document *Document, field string) {
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
