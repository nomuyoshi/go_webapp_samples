package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"webapp_samples/chapter4/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	bigHuge := thesaurus.NewBigHuge(apiKey)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := bigHuge.Synonyms(word)
		if err != nil {
			log.Fatalln("類語検索に失敗しました。err: ", err)
		}
		if len(syns) == 0 {
			fmt.Printf("「%s」の類語はありませんでした。\n", word)
		}

		for i := range syns {
			fmt.Println(syns[i])
		}
	}
}
