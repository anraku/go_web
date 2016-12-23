package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

func main() {
	// 文字列変換後のパターンをテキストファイルから読み込み
	filename := "statics/patterns.txt"
	text, err := ioutil.ReadFile(filename)
	// 読み込んだテキストデータを配列に変換
	var strs []string
	strs = strings.Split(string(text), "\n")
	if err != nil {
		fmt.Println("error")
		os.Exit(0)
	}

	// 入力値を用意したパターンに応じて変換
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := strs[rand.Intn(len(strs))]
		fmt.Println(strings.Replace(t, "*", s.Text(), -1))
	}
}
