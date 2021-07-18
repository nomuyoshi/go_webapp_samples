package main

import (
	"log"
	"os"
	"os/exec"
)

var cmdChain = []*exec.Cmd{
	exec.Command("lib/synonyms"),  // synonyms で類語リスト取得
	exec.Command("lib/sprinkle"),  // sprinkel で単語を改変（しない場合もある）
	exec.Command("lib/coolify"),   // coolify で母音字の数を増やす（かも）
	exec.Command("lib/domainify"), // domainify でドメインに使用可能文字のみにする
	exec.Command("lib/available"), // available で利用可能なドメインか判定
}

func main() {
	// １つ目のコマンドの標準入力をos.Stdin(domainfinderにとっての標準入力)にする
	cmdChain[0].Stdin = os.Stdin
	// 最後のコマンドの標準出力をos.Stdout(domainfinderにとっての標準出力)にする
	cmdChain[len(cmdChain)-1].Stdout = os.Stdout

	// それぞれのコマンドの標準出力を次のコマンドの標準入力に設定
	for i := 0; i < len(cmdChain)-1; i++ {
		thisCmd := cmdChain[i]
		nextCmd := cmdChain[i+1]
		stdout, err := thisCmd.StdoutPipe()
		if err != nil {
			log.Panicln(err)
		}
		nextCmd.Stdin = stdout
	}

	// コマンド実行
	for _, cmd := range cmdChain {
		if err := cmd.Start(); err != nil {
			log.Panicln(err)
		} else {
			defer cmd.Process.Kill()
		}
	}

	// コマンド完了待ち
	for _, cmd := range cmdChain {
		if err := cmd.Wait(); err != nil {
			log.Panicln(err)
		}
	}
}
