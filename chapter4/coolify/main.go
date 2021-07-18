package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func main() {
	rand.Seed(time.Now().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := []byte(s.Text())
		if !randBool() {
			fmt.Println(string(word))
			continue
		}
		var vI int = -1
		for i, char := range word {
			switch char {
			case 'a', 'i', 'u', 'e', 'o', 'A', 'I', 'U', 'E', 'O':
				if randBool() {
					vI = i
				}
			}
		}
		if vI >= 0 {
			if randBool() {
				word = duplicateVowel(word, vI)
			} else {
				word = removeVowel(word, vI)
			}
		}
		fmt.Println(string(word))
	}
}

func duplicateVowel(word []byte, i int) []byte {
	return append(word[:i+1], word[i:]...)
}

func removeVowel(word []byte, i int) []byte {
	return append(word[:i], word[i+1:]...)
}
