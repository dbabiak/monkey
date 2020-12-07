package main

import (
	"fmt"
	"github.com/dbabiak/dbgo"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	LineNum int
	Col     int
}

const (
	LET         = TokenType("LET")
	FUNCTION    = TokenType("FUNCTION")
	IDENTIFIER  = TokenType("IDENTIFIER")
	EQUALS      = TokenType("EQUALS")
	PLUS        = TokenType("PLUS")
	SEMICOLON   = TokenType("SEMICOLON")
	NUMBER      = TokenType("NUMBER")
	OPEN_PAREN  = TokenType("OPEN_PAREN")
	CLOSE_PAREN = TokenType("CLOSE_PAREN")
	OPEN_BRACE  = TokenType("OPEN_BRACE")
	CLOSE_BRACE = TokenType("CLOSE_BRACE")
	COMMA       = TokenType("COMMA")
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
}

func isLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func Identifier(xs []rune, i int) (nextPos int, literal string, matched bool) {
	assert(i < len(xs), "i must be in bounds")
	j := i
	for {
		if !isLetter(xs[j]) {
			break
		}
		j += 1
	}
	if i == j {
		return 0, "", false
	}
	return j, string(xs[i:j]), true
}

func Number(xs []rune, i int) (nextPos int, literal string, matched bool) {
	assert(i < len(xs), "i must be in bounds")
	j := i
	for {
		if !isDigit(xs[j]) {
			break
		}
		j += 1
	}
	if i == j {
		return 0, "", false
	}
	return j, string(xs[i:j]), true
}


func TokType(ident string) TokenType {
	// have we lexed a keyword or a regular identifier?
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENTIFIER
}

func nextToken(line []rune, linenum, col int) Token {
	assert(col < len(line), "i must be in bounds")
	switch line[col] {
	case '(':
		return Token{OPEN_PAREN, "(", linenum, col}
	case ')':
		return Token{CLOSE_PAREN, ")", linenum, col}
	case '{':
		return Token{OPEN_BRACE, "{", linenum, col}
	case '}':
		return Token{CLOSE_BRACE, "}", linenum, col}
	case '=':
		return Token{EQUALS, "=", linenum, col}
	case ',':
		return Token{COMMA, ",", linenum, col}
	case ';':
		return Token{SEMICOLON, ";", linenum, col}
	case '+':
		return Token{PLUS, "+", linenum, col}
	}

	_, literal, match := Identifier(line, col)
	if match {
		return Token{TokType(literal), literal, linenum, col}
	}

	_, literal, match = Number(line, col)
	if match {
		return Token{NUMBER, literal, linenum, col}
	}

	panic(fmt.Sprintf("line %d col %d cannot lex %#v", linenum, col, string(line[col:])))
}

func Lex(filename string) []Token {
	lines := dbgo.PathToLines(filename)

	tokens := []Token{}
	for linenum, line := range lines {
		// this is safe b/c no way for a token to be split by a new line
		runes := []rune(line)
		NextRune:
		for i := 0; i < len(runes); {
			if runes[i] == '\\' || runes[i] == ' ' {
				i += 1
				continue NextRune
			}

			token := nextToken(runes, linenum, i)
			i += len(token.Literal)
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func main() {
	tokens := Lex("data/demo.monkey")
	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}
}
