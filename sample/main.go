// MeCab で形態素解析し、NGram 解析してみるサンプル

package main

import (
	"fmt"
	"os"

	"github.com/bunji2/mecab"
)

func main() {
	os.Exit(run())
}

func run() int {
	conf := mecab.Config{
		Commands: []string{
			"mecab",
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

	doc := mecab.MakeDoc(r)

	ng := NewNGram(2, &doc)
	ng.dump()
	ng = NewNGram(3, &doc)
	ng.dump()

	return 0
}
