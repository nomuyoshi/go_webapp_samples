package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var marks = map[bool]string{true: "OK", false: "NG"}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		domain := s.Text()
		// まだ改行したくないのでPrintを使う
		fmt.Print(domain, " ")
		exist, err := exists(domain)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(marks[!exist])
		time.Sleep(1 * time.Second)
	}
}

func exists(domain string) (bool, error) {
	const whoisServer string = "com.whois-servers.net"
	// whoisServerで指定したサーバーのポート43にたして net.Dial で接続を開く
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return false, err
	}
	// deferで最終的に接続を閉じる
	defer conn.Close()
	// ドメイン名と\r\nを送信
	// （WHOISの仕様で定義されているのはここまで）
	conn.Write([]byte(domain + "\r\n"))
	// レスポンスを１行ずつ読み込む
	s := bufio.NewScanner(conn)
	for s.Scan() {
		// あるWHOISサーバーではドメインに関する情報がなければ「No match」というメッセージが含まれる
		// 大文字・小文字の区別をしないように、全て小文字にして比較する
		if strings.Contains(strings.ToLower(s.Text()), "no match") {
			return false, nil
		}
	}

	return true, nil
}
