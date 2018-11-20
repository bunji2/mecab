package mecab

import (
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

// Config :
type Config struct {
	Commands        []string
	TimeOutSec      int
	StopWords       []string
	StopWordClasses []string
	Separator       string
}

var defaultCommands = []string{
	"mecab",
	"-F%f[0]_%m\n",
	"-E ",
}

var conf Config

// Init :
func Init(c Config) (err error) {
	conf = c

	if len(conf.Commands) < 1 {
		conf.Commands = defaultCommands
	}

	if conf.TimeOutSec < 1 {
		conf.TimeOutSec = 10
	}

	if conf.Separator == "" {
		conf.Separator = "\r\n"
	}
	return
}

// Proc :
type Proc struct {
	err    error
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// NewProc :
func NewProc() (mp *Proc, err error) {
	mp = &Proc{}
	mp.Reset()
	err = mp.Error()
	return
}

// Reset :
func (mp *Proc) Reset() {
	mp.cmd = exec.Command(conf.Commands[0], conf.Commands[1:]...)

	// 標準入力のパイプの取得
	mp.stdin, mp.err = mp.cmd.StdinPipe()
	if mp.err != nil {
		return
	}

	// 標準出力のパイプの取得
	mp.stdout, mp.err = mp.cmd.StdoutPipe()
}

// Close :
func (mp *Proc) Close() {
	if mp.stdin != nil {
		mp.stdin.Close()
		mp.stdin = nil
	}
	if mp.stdout != nil {
		mp.stdout.Close()
		mp.stdout = nil
	}

}

// Write : 与えた文字列を品詞に分解する関数
func (mp *Proc) Write(text string) (ret []string) {

	if mp.cmd == nil {
		mp.Reset()
	}

	if mp.err != nil {
		return
	}

	// プロセスの起動
	mp.err = mp.cmd.Start()
	if mp.err != nil {
		return
	}

	if mp.stdin == nil {
		return
	}

	_, mp.err = io.WriteString(mp.stdin, text)
	if mp.err != nil {
		return
	}

	mp.stdin.Close()
	mp.stdin = nil

	r := ""
	done := make(chan error, 1)

	go func() { // [GoRoutine#1]

		// 標準出力の読み出し
		var b []byte
		b, mp.err = ioutil.ReadAll(mp.stdout)
		if mp.err == nil {
			r = string(b)
		} else {
			done <- nil
			return
		}

		// プロセスの待ち合わせ
		done <- mp.cmd.Wait()
	}() // [GoRoutine#1]

	// プロセス終了の同期待ち
	select {
	// 所定の処理時間内に終了した場合
	case <-done:
		break

	// 所定の処理時間を超過した場合(処理時間>time_out_sec)
	case <-time.After(time.Duration(conf.TimeOutSec) * time.Second):

		// 所定の処理時間を超過したプロセスを強制終了
		mp.err = mp.cmd.Process.Kill()
	}
	if r == "" {
		return
	}

	if mp.stdout != nil {
		mp.stdout.Close()
		mp.stdout = nil
	}
	if mp.cmd != nil {
		mp.cmd = nil
	}

	ret = strings.Split(r, conf.Separator)
	return
}

// Error :
func (mp *Proc) Error() (err error) {
	err = mp.err
	return
}

// Doc :
type Doc struct {
	Dic   []string
	Words []int
}

// MakeDoc :
func MakeDoc(words []string) (r Doc) {
	bag := map[string]int{}
	w := []int{}
	i := 0
	for _, word := range words {
		if IsStopWord(word) {
			continue
		}
		_, ok := bag[word]
		if !ok {
			bag[word] = i
			i++
		}
		idx := bag[word]
		w = append(w, idx)
	}

	dic := make([]string, len(bag))
	for word, idx := range bag {
		dic[idx] = word
	}
	r = Doc{
		Dic:   dic,
		Words: w,
	}
	return
}

// IsStopWord :
func IsStopWord(x string) (r bool) {
	for _, word := range conf.StopWords {
		if x == word {
			r = true
			return
		}
	}
	for _, wordClass := range conf.StopWordClasses {
		if strings.HasPrefix(x, wordClass) {
			r = true
			break
		}
	}
	return
}
