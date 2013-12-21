package main

type Trie struct {
	Value byte
	IDs   []int
	Nodes map[byte]*Trie
}

func (t *Trie) Find(word []byte) (*Trie, bool) {
	current := t
	for _, c := range word {
		if found, ok := current.Nodes[c]; ok {
			current = found
		} else {
			return nil, false
		}
	}
	return current, true
}

func (t *Trie) Add(id int, word []byte) {
	if t.Nodes == nil {
		t.Nodes = map[byte]*Trie{}
	}

	next, ok := t.Nodes[word[0]]
	if !ok {
		next = &Trie{Value: word[0], Nodes: map[byte]*Trie{}}
		t.Nodes[word[0]] = next
	}

	if len(word) == 1 {
		next.IDs = append(next.IDs, id)
	} else {
		next.Add(id, word[1:])
	}
}
