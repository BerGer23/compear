package main

type Analysis struct {
	FindingsLeft  []Token
	FindingsRight []Token
}

type Token struct {
	Index   int
	Content string
	Added   bool
}

func NewToken(index int, content string) Token {
	token := Token{}
	token.Index = index
	token.Content = content
	token.Added = false
	return token
}
