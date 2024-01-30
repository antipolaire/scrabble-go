package game

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDictionary(t *testing.T) {
	// Test the ParseFlags function
	t.Run("Test Load Dictionary", func(t *testing.T) {
		dict := LoadDictionary("../assets/dicts/en.csv")
		if dict == nil {
			t.Errorf("Dictionary not loaded")
		}
		assert.Equal(t, dict.IsWord("hello"), true)
		prefixes := dict.WordFinder.FindAllPrefixesOf("abortivenesses")
		expectedPrefixes := []string{"ab", "abo", "abort", "abortive", "abortiveness", "abortivenesses"}
		for _, prefix := range prefixes {
			assert.Contains(t, expectedPrefixes, prefix.Word)
		}
	})
	t.Run("Test Load Dictionary from DAWG", func(t *testing.T) {
		dict := NewDictionaryFromDAWG("../assets/dicts/en.dawg")
		if dict == nil {
			t.Errorf("Dictionary not loaded")
		}
		assert.Equal(t, dict.IsWord("hello"), true)
		prefixes := dict.WordFinder.FindAllPrefixesOf("abortivenesses")
		expectedPrefixes := []string{"ab", "abo", "abort", "abortive", "abortiveness", "abortivenesses"}
		for _, prefix := range prefixes {
			assert.Contains(t, expectedPrefixes, prefix.Word)
		}
	})
	// Load all dicts
	t.Run("Test Load All Dictionaries", func(t *testing.T) {
		languages := []string{"en", "de", "fr", "es"}
		for _, language := range languages {
			dict := LoadDictionary("../assets/dicts/" + language + ".csv")
			if dict == nil {
				t.Errorf("Failed to load dictionary for language: %s", language)
			}
			dict.ExportStatsCSV("../assets/dicts/" + language + ".stats.csv")
		}
	})

	t.Run("Get Some Word Stats", func(t *testing.T) {
		dict := NewDictionaryFromDAWG("../assets/dicts/en.dawg")
		wordStats := dict.GetWordStats("axolotl")
		fmt.Println(wordStats)
	})

	t.Run("Find 20 Highest Scoring Words For Each Language", func(t *testing.T) {
		// languages := []string{"en", "de", "fr", "es"}
		languages := []string{"de2"}
		for _, language := range languages {
			dict := LoadDictionary("../assets/dicts/" + language + ".csv")
			if dict == nil {
				t.Errorf("Failed to load dictionary for language: %s", language)
			}
			n := 20
			words := dict.FindNHighestScoringWords(n)
			fmt.Printf("Top %d scoring words for language %s:\n", n, language)
			for _, word := range words {
				fmt.Println(fmt.Sprintf("\t%s: %d", word, CalculateWordScore(word)))
			}

			words = dict.FindNLongestWords(n)
			fmt.Printf("Top %d longest words for language %s:\n", n, language)
			for _, word := range words {
				fmt.Println(fmt.Sprintf("\t%s: %d", word, CalculateWordScore(word)))
			}

			words = dict.FindNHighestRelativeScoringWords(n)
			fmt.Printf("Top %d words with highest relative score for language %s:\n", n, language)
			for _, word := range words {
				fmt.Println(fmt.Sprintf("\t%s: %d", word, CalculateWordScore(word)))
			}
		}
	})
}
