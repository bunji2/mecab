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
}

var defaultCommands = []string{
	"C:/Program Files (x86)/MeCab/bin/mecab.exe",
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

	ret = strings.Split(r, "\r\n")
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

/*
// 関数：PyRun_from_string
// 概要：Pythonスクリプトを実行する。
// 引数：pyscript --- Pythonスクリプトの文字列
// 戻り値：result --- 0:成功、1:失敗
// 戻り値：out_text --- スクリプト実行で得られた標準出力の文字列
// 戻り値：err_text --- スクリプト実行で得られた標準エラーの文字列
// メモ：外部コマンド python.exe を実行し、パイプでスクリプトを
//       標準入力で受け渡し、その標準出力と標準エラーを取得する。
func PyRun_from_string(pyscript string) (int, string, string) {
	result := 0
	out_string := ""
	err_string := ""

	// Pythonコマンドのオブジェクト。
	// 引数 "-" は標準入力からスクリプトを取得する。
	cmd := exec.Command(conf.Python_exe, "-")

	// 標準入力のパイプの取得とスクリプトの送出
	if stdin, err := cmd.StdinPipe(); err == nil {
		if stdin != nil {
		    io.WriteString(stdin, pyscript)
    		stdin.Close()
		} else {
			return 1,"","nil pointer"
		}
	} else {
		return 1, "", "failed in opening stdin"
	}

	// 標準出力のパイプの取得
	stdout, err2 := cmd.StdoutPipe()
	if err2 != nil {
		return 1, "", "failed in opening stdout"
	}

	// 標準エラーのパイプの取得
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 1, "", "failed in opening stderr "
	}

	// プロセスの起動
	if err := cmd.Start(); err != nil {
		return 1,"","failed in starting"
	}

	done := make(chan error, 1)

	go func() { // [GoRoutine#1]

		// 標準出力の読み出し
		if b, err := ioutil.ReadAll(stdout); err == nil {
			out_string = string(b)
		} else {
			result = 1
			err_string = "failed in reading stdout"
			done <- nil
			return
		}

		// 標準エラーの読み出し
		if b, err := ioutil.ReadAll(stderr); err == nil {
			err_string = string(b)
		} else {
			result = 1
			err_string = "failed in reading stderr"
			done <- nil
			return
		}

		// プロセスの待ち合わせ
		err := cmd.Wait()
		if err != nil {
			//return 1,"","failed in wating"
			result = 1
			if exiterr, ok := err.(*exec.ExitError); ok {

				// プロセスの「終了コード」の取り出し
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					err_string = fmt.Sprintf(
						"%s\nexection failed (exit code=%d)\n",
						err_string,
						status.ExitStatus())
				}
			}
		}
		done <- err
	}() // [GoRoutine#1]

	// プロセス終了の同期待ち
	select {
		// 所定の処理時間内に終了した場合
		case <- done:
			break
			// result は 0 のまま。

		// 所定の処理時間を超過した場合(処理時間>time_out_sec)
		case <- time.After(time.Duration(conf.TimeOutSec) * time.Second):
			result = 1
			err_string = "timeout."
        	//log.Fatal(err_string, "pid =", cmd.Process.Pid)

			// 所定の処理時間を超過したプロセスを強制終了
			if err := cmd.Process.Kill(); err != nil {
				result = 1
				err_string += "\nfailed to kill:" + err.Error()
        			log.Fatal("failed to kill: ", err)
			}

	}

	// result --- 0:成功、1:失敗
	return result, out_string, err_string

	// TODO: 返り値の "result" は error 型にすべき。(GoLang の流儀)
}
*/
