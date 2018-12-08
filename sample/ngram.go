package main

import (
	"fmt"
	"sort"

	"github.com/bunji2/mecab"
)

type NGram struct {
	n        int
	numWords int
	doc      *mecab.Doc
	counts   map[int]int
	data     []int
}

func (ng *NGram) Len() int {
	return len(ng.counts)
}

func (ng *NGram) Less(i, j int) bool {
	return ng.counts[ng.data[i]] > ng.counts[ng.data[j]]
}

func (ng *NGram) Swap(i, j int) {
	ng.data[i], ng.data[j] = ng.data[j], ng.data[i]
}

func (ng *NGram) sortedList() []int {
	ng.data = make([]int, len(ng.counts))
	i := 0
	for id := range ng.counts {
		ng.data[i] = id
		i++
	}
	sort.Sort(ng)
	return ng.data
}

func NewNGram(n int, doc *mecab.Doc) (r *NGram) {
	counts := map[int]int{}
	numWords := len(doc.Dic)

	for i := 0; i < len(doc.Words)-n; i++ {
		id := seqToID(numWords, doc.Words[i:i+n])
		counts[id] = counts[id] + 1
	}

	r = &NGram{
		n:        n,
		numWords: numWords,
		doc:      doc,
		counts:   counts,
	}
	return
}

func (ng *NGram) dump() {
	for _, id := range ng.sortedList() {
		fmt.Println(ng.idToStrs(id), ":", ng.counts[id])
	}
}

func (ng *NGram) idToStrs(id int) (r []string) {
	seq := idToSeq(ng.numWords, ng.n, id)
	r = make([]string, len(seq))
	for i, x := range seq {
		r[i] = ng.doc.Dic[x]
	}
	return
}

func seqToID(numWords int, seq []int) (id int) {
	for _, x := range seq {
		id = id*numWords + x
	}
	return
}

func idToSeq(numWords, numSeq, id int) (seq []int) {
	seq = make([]int, numSeq)
	for i := numSeq - 1; i >= 0; i-- {
		seq[i] = id % numWords
		id = id / numWords
	}
	return
}