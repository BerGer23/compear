package main

import (
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func tokenize(entry string, separator string, trim bool) []token {
	var sep = "\n"
	if separator == "Comma" {
		sep = ","
	} else if separator == "Space" {
		sep = " "
	}
	tokenStrings := strings.Split(entry, sep)
	tokens := make([]token, 0)

	for i := 0; i < len(tokenStrings); i++ {
		if trim {
			tokenStrings[i] = strings.TrimSpace(tokenStrings[i])
		}
		tokens = append(tokens, token{i, tokenStrings[i]})
	}
	sort.Slice(tokens, func(i, j int) bool {
		return strings.Compare(tokens[i].content, tokens[j].content) < 0
	})

	return tokens
}

func detectSeparator(entry string) string {
	var count = strings.Count(entry, "\n")
	log.Println("found " + strconv.Itoa(count) + " newlines")
	if count > 1 {
		return "Newline"
	}
	count = strings.Count(entry, ",")
	log.Println("found " + strconv.Itoa(count) + " newlines")
	if count > 1 {
		return "Comma"
	}
	return "Space"
}

func compareTokens(leftEntry string, rightEntry string, leftSeparator string, rightSeparator string, trim bool) analysis {
	leftTokens := tokenize(leftEntry, leftSeparator, trim)
	rightTokens := tokenize(rightEntry, rightSeparator, trim)

	if !reflect.DeepEqual(leftTokens, rightTokens) {
		findingsLeft := make([]token, 0)
		findingsRight := make([]token, 0)
		if len(leftTokens) > 0 && len(rightTokens) > 0 {
			i := 0
			j := 0
			for {
				result := strings.Compare(leftTokens[i].content, rightTokens[j].content)
				log.Println(result)
				if result == 0 {
					if i < len(leftTokens)-1 {
						i = i + 1
					}
					if j < len(rightTokens)-1 {
						j = j + 1
					}
				} else if result == -1 {
					if i < len(leftTokens)-1 {
						findingsLeft = append(findingsLeft, leftTokens[i])
						i = i + 1
					} else if j < len(rightTokens)-1 {
						findingsRight = append(findingsRight, rightTokens[j])
						j = j + 1
					}
				} else if result == 1 {
					if j < len(rightTokens)-1 {
						findingsRight = append(findingsRight, rightTokens[j])
						j = j + 1
					} else if i < len(leftTokens)-1 {
						findingsLeft = append(findingsLeft, leftTokens[i])
						i = i + 1
					}
				}
				log.Println("i: " + strconv.Itoa(i) + ", j: " + strconv.Itoa(j))
				if i == len(leftTokens)-1 && j == len(rightTokens)-1 {
					if strings.Compare(leftTokens[i].content, rightTokens[j].content) != 0 {
						findingsLeft = append(findingsLeft, leftTokens[i])
						findingsRight = append(findingsRight, rightTokens[j])
					}
					break
				}
			}
		}
		sort.Slice(findingsLeft, func(i, j int) bool {
			return findingsLeft[i].index < findingsLeft[j].index
		})
		sort.Slice(findingsRight, func(i, j int) bool {
			return findingsRight[i].index < findingsRight[j].index
		})
		return analysis{findingsLeft, findingsRight}
	} else {
		return analysis{}
	}
}
