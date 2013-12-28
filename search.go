package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var currentID = 0
var Mappings map[string]string
var allDocuments map[int]*Document
var allIndexes map[string]*DataIndex

type DataIndex struct {
	WordNodes     map[string]*WordNode
	LastWordNodes map[string]*WordNode
	Sort          []Sort
}

func generateIndexName(sort []Sort) string {
	name := ""
	for _, sort := range sort {
		name += fmt.Sprintf("(%v-%v)", sort.Field, sort.Ascending)
	}
	return name
}

func FindIndex(sort []Sort) (*DataIndex, error) {
	name := generateIndexName(sort)
	index, exists := allIndexes[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("Couldn't find index %q", name))
	}
	return index, nil
}

func CreateIndex(order []Sort) {
	if allIndexes == nil {
		allIndexes = map[string]*DataIndex{}
	}

	name := generateIndexName(order)
	index := &DataIndex{Sort: order}
	index.WordNodes = map[string]*WordNode{}
	index.LastWordNodes = map[string]*WordNode{}
	allIndexes[name] = index

	var documents []*Document
	for _, document := range allDocuments {
		documents = append(documents, document)
	}

	if len(order) != 0 {
		sorter := BySort{documents, order}
		sort.Sort(sorter)
	}

	for _, document := range documents {
		for field, mapping := range Mappings {
			if mapping == "text" {
				indexTextField(index, document, field)
			}
		}
	}
}

type WordNode struct {
	Document *Document
	Child    *WordNode
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
	for _, sortType := range s.Sort {
		if Mappings[sortType.Field] == "text" {
			a, b := s.Documents[i].Source[sortType.Field].(string), s.Documents[j].Source[sortType.Field].(string)
			a, b = strings.ToLower(a), strings.ToLower(b)
			if a != b {
				if sortType.Ascending {
					return a < b
				} else {
					return a > b
				}
			}
		} else if Mappings[sortType.Field] == "integer" {
			a, b := s.Documents[i].Source[sortType.Field].(int), s.Documents[j].Source[sortType.Field].(int)
			if a != b {
				if sortType.Ascending {
					return a < b
				} else {
					return a > b
				}
			}
		}
	}

	return false
}

func Search(query Query) Documents {
	index, err := FindIndex(query.Sort)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, text := range query.Text {
		results := searchTextField(index, text.Field, text.Value)
		return results
	}
	return nil
}

func get(id int) *Document {
	return allDocuments[id]
}

func searchTextField(index *DataIndex, field string, query string) Documents {
	var documents Documents

	words := splitWords(query)
	requiredMatches := len(words)
	matches := map[*Document]int{}

	for _, word := range words {
		if node, ok := index.WordNodes[word]; ok {
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
}

var stopWords map[string]bool

func isStopWord(word string) bool {
	if stopWords == nil {
		stopWords = map[string]bool{}
		words := []string{
			"a", "an", "and", "are", "as", "at", "be", "but", "by",
			"for", "if", "in", "into", "is", "it",
			"no", "not", "of", "on", "or", "such",
			"that", "the", "their", "then", "there", "these",
			"they", "this", "to", "was", "will", "with",
		}
		for _, word := range words {
			stopWords[word] = true
		}
	}

	return stopWords[word]
}

func splitWords(s string) []string {
	s = strings.ToLower(s)
	var words []string
	for _, word := range strings.Split(s, " ") {
		if word != "" && isStopWord(word) == false {
			words = append(words, word)
		}
	}

	return words
}

func indexTextField(index *DataIndex, document *Document, field string) {
	value := document.Source[field].(string)
	words := splitWords(value)

	for _, word := range words {
		node, exists := index.WordNodes[word]
		if exists {
			lastWordNode := index.LastWordNodes[word]
			lastWordNode.Child = &WordNode{Document: document}
			index.LastWordNodes[word] = lastWordNode.Child
		} else {
			node = &WordNode{Document: document}
			index.WordNodes[word] = node
			index.LastWordNodes[word] = node
		}
	}
}
