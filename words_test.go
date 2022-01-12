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
	assert.Equal(t, 18, len(s))
	assert.Contains(t, s, "added")
	assert.Contains(t, s, "dead")
	assert.Contains(t, s, "deaf")
	assert.Contains(t, s, "fade")
	assert.Contains(t, s, "faded")
}

func TestWordsMultipleMatches2(t *testing.T) {
	s, err := GetMatchingWords("tiplvxe", 'e')
	assert.Nil(t, err)
	assert.Equal(t, 35, len(s))
	assert.Contains(t, s, "peel")
	assert.Contains(t, s, "expletive")
	assert.Contains(t, s, "pile")
	assert.Contains(t, s, "tipple")
}
