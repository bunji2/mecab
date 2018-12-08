package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/bunji2/mecab"
)

func main() {
	os.Exit(run())
}

func run() int {
	conf := mecab.Config{
		//Separator:       "\n",
		UseStopWords:    false,
		StopWords:       []string{"名詞_(", "名詞_)"},
		StopWordClasses: []string{"助詞", "助動詞", "記号"},
		Commands: []string{
			"mecab",
			//"-F%f[0]_%m\n",
			"-F%m\n",
			"-E ",
		},
	}
	mecab.Init(conf)
	mp, err := mecab.NewProc()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	text := `当該製品の管理画面にアクセスされ任意のコマンドを実行される - CVE-2018-0676
	同一 LAN 内の当該製品に管理者権限でアクセス可能なユーザによって、任意の OS コマンドを実行される - CVE-2018-0677
	同一 LAN 内の当該製品に管理者権限でアクセス可能なユーザによって、任意のコードを実行されたり、サービス運用妨害 (DoS) 攻撃を受けたりする - CVE-2018-0678`
	r := mp.Write(text)
	if mp.Error() != nil {
		fmt.Println(mp.Error())
		return 2
	}
	fmt.Println(r)

	r2 := mp.Write(text)
	if mp.Error() != nil {
		fmt.Println(mp.Error())
		return 2
	}

	fmt.Println(len(r2))
	doc := mecab.MakeDoc(r2)
	//fmt.Println(mecab.MakeDoc(r2))

	//for _, wid := range doc.Words {
	//	fmt.Println(doc.Dic[wid])
	//}

	/*
		counts := map[int]int{}
		numWords := len(doc.Dic)
		for i := 0; i < len(doc.Words)-2; i++ {
			//id := doc.Words[i]*numWords + doc.Words[i+1]
			id := foo(numWords, []int{doc.Words[i+2], doc.Words[i+1], doc.Words[i]})
			counts[id] = counts[id] + 1
		}

		for id, count := range counts {
			fmt.Printf("%d %v %d\n", id, toWords(doc.Dic, bar(numWords, id)), count)
		}
	*/
	//fmt.Println(doc)

	ng := NewNGram(2, &doc)
	ng.dump()
	ng = NewNGram(3, &doc)
	ng.dump()

	aa := []int{0, 1, 2, 3}
	id := seqToID(4, aa)
	fmt.Println(aa, "==>", id, "==>", idToSeq(4, len(aa), id))

	aa = []int{3, 2, 1, 0}
	id = seqToID(4, aa)
	fmt.Println(aa, "==>", id, "==>", idToSeq(4, len(aa), id))

	return 0
}

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
