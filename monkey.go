package main

import (
	"fmt"
	"github.com/dbabiak/dbgo"
	"strings"
)

type TokenType string

type Token struct {
	Type TokenType
	Literal string
	LineNum int
	Col int
}

func Let(xs []rune, i int) (nextPos int, literal string, matched bool) {
	const token = "let"
	if len(token) <= len(xs) - i  {
		if literal := string(xs[i:i+len(token)]); literal == token {
			return i + len(token), literal, true
		}
	}
	return 0, "", false
}

func main() {
	dbgo.System("cat data/demo.monkey")
	println("\n\n---------------\n")

	lines := dbgo.PathToLines("data/demo.monkey")

	tokens := []Token{}
	for linenum, line := range lines {
		// this is safe b/c no way for a token to be split by a new line
		runes := []rune(strings.TrimSpace(line))
		println(linenum, line)
		for i := 0; i < len(runes); i++ {
			if runes[i] == '\\' {
				continue
			}

			next, literal, matched := Let(runes, i)
			if matched {
				tokens = append(tokens, Token{
					Type: TokenType("LET"),
					Literal: literal,
					LineNum: linenum,
					Col: i,
				})
				i = next
			}
		}
	}
	fmt.Printf("%v\n", tokens)
}
