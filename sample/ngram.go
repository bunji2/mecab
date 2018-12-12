package main

import (
	"fmt"
	"sort"

	"github.com/bunji2/doc"
)

type NGram struct {
	n        int
	numWords int
	//doc      *mecab.Doc
	doc    *doc.Data
	counts map[int]int
	data   []int
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

// NewGram
func NewNGram(n int, d *doc.Data) (r *NGram) {
	counts := map[int]int{}
	numWords := len(d.Dic)

	for i := 0; i < len(d.Seq)-n; i++ {
		id := seqToID(numWords, d.Seq[i:i+n])
		counts[id] = counts[id] + 1
	}

	r = &NGram{
		n:        n,
		numWords: numWords,
		doc:      d,
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
