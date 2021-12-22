package main

import (
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
