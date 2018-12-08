package mecab_test

func Example() {
	conf := mecab.Config{}
	mecab.Init(conf)
	mp, err := mecab.NewProc()
	if err != nil {
		fmt.Println(err)
		return
	}
	text := `それでも暮らしは続くから全てを今忘れてしまう為には全てを今知っている事が条件で僕にはとても無理だから一つずつ忘れて行く為に愛する人達と手を取り分け合ってせめて思い出さないように暮らしを続けて行くのです`
	r := mp.Write(text)
	if mp.Error() != nil {
		fmt.Println(mp.Error())
		return
	}
	fmt.Println(r)

}
