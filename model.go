package main

type analysis struct {
	findingsLeft []token
	findingsRight []token
}

type token struct {
	index int
	content string
}