package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Response struct {
	Letters         string
	Words           []string
	MandatoryLetter string
	ExecutionTimeMs int
	Created         time.Time
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println(uniqueLettersSorted("Starting"))
	//dat, err := os.Open("./short.txt")

	fmt.Println(GetMatchingWords("wniougk", 'w'))
	fmt.Println(GetMatchingWordsResponse("wniougk", 'w'))

}
func GetMatchingWordsResponse(letters string, mandatoryChar rune) string {
	start := time.Now()
	words := GetMatchingWords(letters, mandatoryChar)
	response := Response{
		Letters:         letters,
		MandatoryLetter: string(mandatoryChar),
		Words:           words,
		Created:         time.Now(),
		ExecutionTimeMs: int(time.Since(start).Milliseconds()),
	}
	j, err := json.Marshal(response)
	check(err)

	//fmt.Println(string(j))
	return string(j)

}

func GetMatchingWords(letters string, mandatoryChar rune) []string {
	dat, err := os.Open("./corncob_lowercase.txt")
	check(err)

	s := make([]string, 0)

	scanner := bufio.NewScanner(dat)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}

	dat.Close()

	//fmt.Println(len(s))

	//make sure chars are unique and sorted
	letters = uniqueLettersSorted(letters)

	result := make([]string, 0)

	for i := 0; i < len(s); i++ {
		if strings.Contains(letters, uniqueLettersSorted(s[i])) {
			//fmt.Println("Found a word: ", s[i])
			if strings.ContainsRune(s[i], mandatoryChar) {
				result = append(result, s[i])
			}
		}
	}
	return result
}

/*
 s is an arbitrary string
 returns a string containing each unique character of s sorted alfabetically and converted to lowecase.
 e.g. uniqueLettersSorted("Hello") -> "ehlo"
*/
func uniqueLettersSorted(s string) string {

	s = strings.ToLower(s)

	//map for unique letters
	m := make(map[string]bool)

	for i := 0; i < len(s); i++ {
		m[string(s[i])] = true
	}

	//array for letters as strings
	keys := make([]string, len(m))

	i := 0

	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	//	fmt.Println("--->", keys)

	var sb strings.Builder

	//convert strings of letters to 1 string
	for _, c := range keys {
		sb.WriteString(c)
	}
	return sb.String()
}
