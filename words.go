package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
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

func main_old() {
	fmt.Println(uniqueLettersSorted("Starting"))
	//dat, err := os.Open("./short.txt")

	fmt.Println(GetMatchingWords("wniougk", 'w'))
	fmt.Println(GetMatchingWordsResponse("wniougk", 'w'))

}
func GetMatchingWordsResponse(letters string, mandatoryChar rune) string {
	start := time.Now()
	words, err := GetMatchingWords(letters, mandatoryChar)
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

func GetMatchingWords(letters string, mandatoryChar rune) ([]string, error) {
	root := "."
	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		root = os.Getenv("LAMBDA_TASK_ROOT")
	}
	fn := root + "/corncob_lowercase.txt"
	err := verify(letters, mandatoryChar)
	if err != nil {
		return nil, err
	}

	dat, err := os.Open(fn) //"./corncob_lowercase.txt")

	if err != nil {
		return nil, errors.New("could not open dictionary. cause:" + err.Error())
	}

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
	return result, nil
}

/* verifies preconditions */
func verify(letters string, mandatoryChar rune) error {
	if len(letters) != 7 {
		return errors.New("letters parameter must be exactly 7 characters. it was " + strconv.Itoa(len(letters)) + ".")
	}

	if !strings.Contains(letters, string(mandatoryChar)) {
		return errors.New("mandatoryChar parameter must be one of the characters in letters string. " + string(mandatoryChar) + " is not in " + letters)
	}

	return nil
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
