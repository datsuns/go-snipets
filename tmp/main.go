package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 1: N テンプレートの文字列数
// 2: テンプレート。空白区切りにN個の文字列
// 3: n 取引先数
// 4: li 各取引先ごとにある取引先情報数
// 5~
// データ不足 := "Error: Lack of data"

type ConvPattern struct {
	A string
	B string
}

type BP struct {
	Patterns []ConvPattern
}

type Config struct {
	N        int
	template []string
	n        int
	BPs      []BP
}

func parse_input(lines []string) (*Config, error) {
	ret := &Config{}
	var err error
	if len(lines) < 4 {
		return nil, errors.New("invalid header lines")
	}
	ret.N, err = strconv.Atoi(lines[0])
	if err != nil {
		fmt.Println("parse N error")
		return nil, err
	}
	ret.template = strings.Split(lines[1], " ")
	ret.n, err = strconv.Atoi(lines[2])
	if err != nil {
		fmt.Println("parse error: template")
		return nil, err
	}

	n := 0
	for nline := 3; nline < len(lines); {
		if n >= ret.n {
			break
		}
		bpLines, err := strconv.Atoi(lines[nline])
		if err != nil {
			fmt.Println("parse error. ", nline, lines[nline])
			return nil, err
		}
		ret.BPs = append(ret.BPs, BP{})
		ret.BPs[n].Patterns = []ConvPattern{}
		nline++
		for i := 0; i < bpLines; i++ {
			param := strings.Split(lines[nline], " ")
			p := ConvPattern{A: param[0], B: param[1]}
			ret.BPs[n].Patterns = append(ret.BPs[n].Patterns, p)
			nline++
		}
		n++
	}

	return ret, nil
}

func getStdin() []string {
	f, _ := os.Open("./test/in/basic/01_mini_01.in")
	stdin, _ := ioutil.ReadAll(f)
	//stdin, _ := ioutil.ReadAll(os.Stdin)
	return strings.Split(strings.TrimRight(string(stdin), "\n"), "\n")
}

func dump(cfg *Config) {
	fmt.Printf(" N: %v\n", cfg.N)
	fmt.Printf(" template: %v\n", cfg.template)
	fmt.Printf(" n: %v\n", cfg.n)
	for i, bp := range cfg.BPs {
		fmt.Printf(" BP[%v]: %v\n", i, bp.Patterns)
	}
}

func replace_completed(input []string) bool {
	for _, s := range input {
		if s[0] == '#' {
			return false
		}
	}
	return true
}

func issue_replace(cfg *Config) []string {
	ret := []string{}
	for i := 0; i < cfg.n; i++ {
		replaced := []string{}
		for _, s := range cfg.template {
			str := s
			for _, p := range cfg.BPs[i].Patterns {
				if str == p.A {
					str = p.B
					break
				}
			}
			replaced = append(replaced, str)
		}

		if replace_completed(replaced) {
			ret = append(ret, strings.Join(replaced, " "))
		} else {
			ret = append(ret, "Error: Lack of data")
		}
	}
	return ret
}

func main() {
	lines := getStdin()
	cfg, err := parse_input(lines)
	if err != nil {
		fmt.Println(err)
		return
	}
	//dump(cfg)
	out := issue_replace(cfg)
	for _, str := range out {
		fmt.Println(str)
	}
}
