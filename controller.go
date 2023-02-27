package main

import (
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func tokenize(entry string, separator string, trim bool) []Token {
	var sep = "\n"
	if separator == "Comma" {
		sep = ","
	} else if separator == "Space" {
		sep = " "
	}
	tokenStrings := strings.Split(entry, sep)
	tokens := make([]Token, 0)

	for i := 0; i < len(tokenStrings); i++ {
		if trim {
			tokenStrings[i] = strings.TrimSpace(tokenStrings[i])
		}
		tokens = append(tokens, NewToken(i, tokenStrings[i]))
	}
	sort.Slice(tokens, func(i, j int) bool {
		return strings.Compare(tokens[i].Content, tokens[j].Content) < 0
	})

	return tokens
}

func detectSeparator(entry string) string {
	if hasMoreThanOne(entry, "\n") {
		return "Newline"
	} else if hasMoreThanOne(entry, ",") {
		return "Comma"
	} else {
		return "Space"
	}
}

func hasMoreThanOne(text string, token string) bool {
	var count = strings.Count(text, token)
	log.Println("found " + strconv.Itoa(count) + token)
	return count > 1
}

func contains(s []Token, str string) bool {
	for _, v := range s {
		if v.Content == str {
			return true
		}
	}

	return false
}

func compareTokens(leftEntry string, rightEntry string, leftSeparator string, rightSeparator string, trim bool) Analysis {
	leftTokens := tokenize(leftEntry, leftSeparator, trim)
	rightTokens := tokenize(rightEntry, rightSeparator, trim)

	if !reflect.DeepEqual(leftTokens, rightTokens) {
		findingsLeft := make([]Token, 0)
		findingsRight := make([]Token, 0)
		if len(leftTokens) > 0 && len(rightTokens) > 0 {
			i := 0
			j := 0
			for {
				result := strings.Compare(leftTokens[i].Content, rightTokens[j].Content)
				log.Println(result)
				if result == 0 {
					if i < len(leftTokens)-1 {
						i = i + 1
					}
					if j < len(rightTokens)-1 {
						j = j + 1
					}
				} else if result == -1 {
					if !leftTokens[i].Added {
						findingsLeft = append(findingsLeft, leftTokens[i])
						leftTokens[i].Added = true
					}
					if i < len(leftTokens)-1 {
						i = i + 1
					} else if j < len(rightTokens)-1 {
						if !rightTokens[j].Added {
							findingsRight = append(findingsRight, rightTokens[j])
							rightTokens[j].Added = true
						}
						j = j + 1
					}
				} else if result == 1 {
					if !rightTokens[j].Added {
						findingsRight = append(findingsRight, rightTokens[j])
						rightTokens[j].Added = true
					}
					if j < len(rightTokens)-1 {
						j = j + 1
					} else if i < len(leftTokens)-1 {
						if !leftTokens[i].Added {
							findingsLeft = append(findingsLeft, leftTokens[i])
							leftTokens[i].Added = true
						}
						i = i + 1
					}
				}
				log.Println("i: " + strconv.Itoa(i) + ", j: " + strconv.Itoa(j))
				if i == len(leftTokens)-1 && j == len(rightTokens)-1 {
					if strings.Compare(leftTokens[i].Content, rightTokens[j].Content) != 0 {
						if !contains(rightTokens, leftTokens[i].Content) && !leftTokens[i].Added {
							findingsLeft = append(findingsLeft, leftTokens[i])
							leftTokens[i].Added = true
						}
						if !contains(leftTokens, rightTokens[j].Content) && !rightTokens[j].Added {
							findingsRight = append(findingsRight, rightTokens[j])
							rightTokens[j].Added = true
						}
					}
					break
				}
			}
		}
		sort.Slice(findingsLeft, func(i, j int) bool {
			return findingsLeft[i].Index < findingsLeft[j].Index
		})
		sort.Slice(findingsRight, func(i, j int) bool {
			return findingsRight[i].Index < findingsRight[j].Index
		})
		return Analysis{findingsLeft, findingsRight}
	} else {
		return Analysis{}
	}
}
