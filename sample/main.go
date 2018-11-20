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
		StopWords:       []string{"名詞_(", "名詞_)"},
		StopWordClasses: []string{"助詞", "助動詞", "記号"},
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
	fmt.Println(r2)

	fmt.Println(mecab.MakeDoc(r2))
	return 0
}
