package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordsNoMatch(t *testing.T) {
	s := GetMatchingWords("zxq", 'z')
	assert.Equal(t, len(s), 0)
}

func TestWordsOneMatch(t *testing.T) {
	s := GetMatchingWords("zo", 'z') //expects [zoo]
	assert.Equal(t, len(s), 1)
}
func TestWordsMultipleMatches(t *testing.T) {
	s := GetMatchingWords("adefpei", 'a') //expects [add added dad dead deaf fade faded]
	assert.Equal(t, len(s), 7)
	assert.Contains(t, s, "added")
	assert.Contains(t, s, "add")
	assert.Contains(t, s, "dad")
	assert.Contains(t, s, "dead")
	assert.Contains(t, s, "deaf")
	assert.Contains(t, s, "fade")
	assert.Contains(t, s, "faded")
}

func TestWordsMultipleMatchesJson(t *testing.T) {
	s := GetMatchingWordsResponse("adefpei", 'a') //expects [add added dad dead deaf fade faded]
	var r Response
	json.Unmarshal([]byte(s), &r)

	assert.Equal(t, r.Letters, "adefpei")
	assert.Equal(t, r.MandatoryLetter, "a")
	assert.GreaterOrEqual(t, r.ExecutionTimeMs, 0)
	assert.Equal(t, len(r.Words), 7)
	assert.Contains(t, r.Words, "added")
	assert.Contains(t, r.Words, "add")
	assert.Contains(t, r.Words, "dad")
	assert.Contains(t, r.Words, "dead")
	assert.Contains(t, r.Words, "deaf")
	assert.Contains(t, r.Words, "fade")
	assert.Contains(t, r.Words, "faded")

	//	fmt.Println(s)
}
