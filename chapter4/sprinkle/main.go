package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

// あえてスペースを入れている。domainifyで除外するため
var transforms = []string{
	otherWord,
	otherWord + " app",
	otherWord + " site",
	otherWord + " time",
	"get " + otherWord,
	"go " + otherWord,
	"lets " + otherWord,
}

func main() {
	rand.Seed(time.Now().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
	if s.Err() != nil {
		fmt.Println("Error: ", s.Err())
	}
}
