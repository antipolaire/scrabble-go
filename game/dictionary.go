package game

// This represents a dictionary in the game of scrabble. It allows loading a dictionary from csv and checking if a word is valid.

import (
	"bufio"
	"fmt"
	"github.com/smhanov/dawg"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

type WordStats struct {
	WordLength        int
	WordScore         int
	WordRelativeScore float64
}

type Dictionary struct {
	WordFinder dawg.Finder
	// Store length by word
	WordStats map[string]WordStats
}

type DictionaryActions interface {
	IsWord(word string) bool
	FindWords(prefix string) []string
	GetWordStats(word string) *WordStats
	ExportStatsCSV(path string)
	FindNLongestWords(n int) []string
	FindNHighestScoringWords(n int) []string
	FindNHighestRelativeScoringWords(n int) []string
	initWordsStats()
}

func filenameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// LoadDictionary loads a dictionary. It first checks if a dawg file exists in the same directory as the csv file.
// If it does, it loads the dawg file, otherwise it creates it and saves it.
func LoadDictionary(path string) *Dictionary {
	// Check if the dawg file exists, if not, create it
	dawgPath := filenameWithoutExtension(path) + ".dawg"
	zap.S().Debugf("Checking if dawg file exists: %s", dawgPath)
	if _, err := os.Stat(dawgPath); os.IsExist(err) {
		return NewDictionaryFromDAWG(dawgPath)
	} else {
		return NewDictionaryFromCSV(path)
	}
}

func NewDictionaryFromDAWG(path string) *Dictionary {
	dictionary := &Dictionary{}

	finder, err := dawg.Load(path)
	if err != nil {
		zap.S().Errorf("Error loading dawg file: %s", err)
		fmt.Println(err)
	}

	dictionary.WordFinder = finder

	dictionary.initWordsStats()

	return dictionary
}

// NewDictionaryFromCSV loads a dictionary from a csv file, creates a DAWG from it and returns a Dictionary
// If the according DAWG file in the same directory does not exist, it will be created. All entries are stored as
// lowercase words.
func NewDictionaryFromCSV(path string) *Dictionary {
	dictionary := &Dictionary{}
	file, err := os.Open(path)
	if err != nil {
		zap.S().Errorf("Error opening file: %s", err)
		log.Fatal(err)
	}
	defer file.Close()

	words := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	// might fail for lines longer than 64K!
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		words[word] = true
	}

	// sort words alphabetically

	keys := make([]string, 0, len(words))
	for k := range words {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create a dawg from the words
	dawgBuilder := dawg.New()
	for _, word := range keys {
		dawgBuilder.Add(word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	dictionary.WordFinder = dawgBuilder.Finish()

	// Save the dawg to a file
	dawgPath := filenameWithoutExtension(path) + ".dawg"
	zap.S().Debugf("Saving dawg file: %s", dawgPath)
	_, err = dictionary.WordFinder.Save(dawgPath)
	if err != nil {
		zap.S().Errorf("Error saving dawg file: %s", err)
	}

	dictionary.initWordsStats()

	return dictionary
}

func (dictionary *Dictionary) FindNLongestWords(n int) []string {
	words := make([]string, 0, len(dictionary.WordStats))
	for word := range dictionary.WordStats {
		words = append(words, word)
	}
	slices.SortFunc(words, func(a, b string) int {
		if dictionary.WordStats[a].WordLength < dictionary.WordStats[b].WordLength {
			return 1
		} else if dictionary.WordStats[a].WordLength > dictionary.WordStats[b].WordLength {
			return -1
		}
		return 0
	})
	return words[:n]
}

func (dictionary *Dictionary) FindNHighestScoringWords(n int) []string {
	words := make([]string, 0, len(dictionary.WordStats))
	for word := range dictionary.WordStats {
		words = append(words, word)
	}
	slices.SortFunc(words, func(a, b string) int {
		if dictionary.WordStats[a].WordScore < dictionary.WordStats[b].WordScore {
			return 1
		} else if dictionary.WordStats[a].WordScore > dictionary.WordStats[b].WordScore {
			return -1
		}
		return 0
	})
	return words[:n]
}

func (dictionary *Dictionary) FindNHighestRelativeScoringWords(n int) []string {
	words := make([]string, 0, len(dictionary.WordStats))
	for word := range dictionary.WordStats {
		words = append(words, word)
	}
	slices.SortFunc(words, func(a, b string) int {
		if dictionary.WordStats[a].WordRelativeScore < dictionary.WordStats[b].WordRelativeScore {
			return 1
		} else if dictionary.WordStats[a].WordRelativeScore > dictionary.WordStats[b].WordRelativeScore {
			return -1
		}
		return 0
	})
	return words[:n]
}

func (dictionary *Dictionary) ExportStatsCSV(path string) {
	file, err := os.Create(path)
	if err != nil {
		zap.S().Errorf("Error creating file: %s", err)
		log.Fatal(err)
	}
	defer file.Close()

	// Write header
	_, err = file.WriteString("word,word_length,word_score,word_relative_score\n")
	if err != nil {
		zap.S().Errorf("Error writing to file: %s", err)
		log.Fatal(err)
	}

	// Write stats
	for word, stats := range dictionary.WordStats {
		_, err = file.WriteString(fmt.Sprintf("%s,%d,%d,%f\n", word, stats.WordLength, stats.WordScore, stats.WordRelativeScore))
		if err != nil {
			zap.S().Errorf("Error writing to file: %s", err)
			log.Fatal(err)
		}
	}
}

func (dictionary *Dictionary) GetWordStats(word string) *WordStats {
	stats, ok := dictionary.WordStats[word]
	if !ok {
		return nil
	}
	return &stats
}

func (dictionary *Dictionary) IsWord(word string) bool {
	return dictionary.WordFinder.IndexOf(strings.ToLower(word)) != -1
}

func (dictionary *Dictionary) FindWords(prefix string) []string {
	findResults := dictionary.WordFinder.FindAllPrefixesOf(strings.ToLower(prefix))
	words := make([]string, len(findResults))
	for i, result := range findResults {
		fmt.Printf("%d: %s\n", i, result.Word)
		words[i] = result.Word
	}
	return words
}

func (dictionary *Dictionary) initWordsStats() {
	dictionary.WordStats = make(map[string]WordStats)
	dictionary.WordFinder.Enumerate(
		func(index int, word []rune, final bool) int {
			if final {
				wordLength := len(word)
				wordScore := CalculateWordScore(string(word))
				dictionary.WordStats[string(word)] = WordStats{
					WordLength:        wordLength,
					WordScore:         wordScore,
					WordRelativeScore: float64(wordScore) / float64(wordLength),
				}
				return dawg.Continue
			}
			return dawg.Continue
		})
}
