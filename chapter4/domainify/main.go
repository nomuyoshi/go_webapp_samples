package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var tlds = []string{"com", "net"}

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func main() {
	rand.Seed(time.Now().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			// 使えない文字が入っていたら無視する
			if !strings.ContainsRune(allowedChars, r) {
				fmt.Printf("警告: 「%q」 は使用できません\n", r)
				continue
			}
			newText = append(newText, r)
		}

		if len(newText) == 0 {
			fmt.Println("利用できる文字が１つもありません。")
		} else {
			fmt.Println(string(newText) + "." + tlds[rand.Intn(len(tlds))])
		}
	}
}
