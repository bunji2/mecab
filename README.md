# Mecab Wrapper for GoLang

# Reuqurement
- Mecab Binary package for MS-Windows http://taku910.github.io/mecab/#install
- The charset of dictionaries must be "UTF-8".

# Sample

```
	import "github.com/bunji2/mecab"

	conf := mecab.Config{
		Commands: []string{
			"C:/Program Files (x86)/MeCab/bin/mecab.exe",
			"-F%f[0]_%m\n",
			"-E ",
		},
		StopWords:       []string{"名詞_(", "名詞_)", ... },
		StopWordClasses: []string{"助詞", "助動詞", "記号", ... },
	}
	mecab.Init(conf)
	mp, err := mecab.NewProc()
	if err != nil {
		// error handling
	}
  
  	text := "Your text"
  
  	result := mp.Write(text)
	if mp.Error() != nil {
		// error handling 2
	}
	fmt.Println(result)
  
```

# Document
- https://godoc.org/github.com/bunji2/mecab
