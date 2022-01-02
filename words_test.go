package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordsNoMatch(t *testing.T) {
	s, err := GetMatchingWords("zxqwvrt", 'z')
	assert.Nil(t, err)
	assert.Equal(t, len(s), 0)
}

func TestErrorIfNot7Chars(t *testing.T) {
	s, err := GetMatchingWords("zo", 'z') //expects error
	assert.NotNil(t, err)
	assert.Nil(t, s)
}

func TestErrorIfNot7CharsDublicates(t *testing.T) {
	s, err := GetMatchingWords("aabcdef", 'a') //expects error
	assert.NotNil(t, err)
	assert.Nil(t, s)
}

func TestErrorIfNotMandatoryCharContainedInLetters(t *testing.T) {
	s, err := GetMatchingWords("abcdefg", 'z')
	assert.NotNil(t, err)
	assert.Nil(t, s)
}

func TestWordsMultipleMatches(t *testing.T) {
	s, err := GetMatchingWords("adefpzi", 'a') //expects [add added dad dead deaf fade faded]
	assert.Nil(t, err)
	assert.Equal(t, len(s), 7)
	assert.Contains(t, s, "added")
	assert.Contains(t, s, "add")
	assert.Contains(t, s, "dad")
	assert.Contains(t, s, "dead")
	assert.Contains(t, s, "deaf")
	assert.Contains(t, s, "fade")
	assert.Contains(t, s, "faded")
}
