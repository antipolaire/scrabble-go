package config

import (
	"github.com/spf13/viper"
	"log"
)

// read config file using viper

const (
	// Special field types:
	NF = iota // NormalField
	DL        // DoubleLetterField
	TL        // TripleLetterField
	DW        // DoubleWordField
	TW        // TripleWordField
	CS        // CenterStarField
	// Special field type colors:
	NFColor = 0xffffff // White
	DLColor = 0x0000ff // Light blue
	TLColor = 0x000080 // Blue
	DWColor = 0xff00ff // Pink
	TWColor = 0xff0000 // Red
	CSColor = 0x101010
)

// Special field matrix:
var specialFields = [][]int{
	{TW, NF, NF, DL, NF, NF, NF, TW, NF, NF, NF, DL, NF, NF, TW},
	{NF, DW, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, DW, NF},
	{NF, NF, DW, NF, NF, NF, DL, NF, DL, NF, NF, NF, DW, NF, NF},
	{DL, NF, NF, DW, NF, NF, NF, DL, NF, NF, NF, DW, NF, NF, DL},
	{NF, NF, NF, NF, DW, NF, NF, NF, NF, NF, DW, NF, NF, NF, NF},
	{NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF},
	{NF, NF, DL, NF, NF, NF, DL, NF, DL, NF, NF, NF, DL, NF, NF},
	{TW, NF, NF, DL, NF, NF, NF, CS, NF, NF, NF, DL, NF, NF, TW},
	{NF, NF, DL, NF, NF, NF, DL, NF, DL, NF, NF, NF, DL, NF, NF},
	{NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF},
	{NF, NF, NF, NF, DW, NF, NF, NF, NF, NF, DW, NF, NF, NF, NF},
	{DL, NF, NF, DW, NF, NF, NF, DL, NF, NF, NF, DW, NF, NF, DL},
	{NF, NF, DW, NF, NF, NF, DL, NF, DL, NF, NF, NF, DW, NF, NF},
	{NF, DW, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, DW, NF},
	{TW, NF, NF, DL, NF, NF, NF, TW, NF, NF, NF, DL, NF, NF, TW},
}

var letterScores = map[string]int{
	"A": 1,
	"B": 3,
	"C": 3,
	"D": 2,
	"E": 1,
	"F": 4,
	"G": 2,
	"H": 4,
	"I": 1,
	"J": 8,
	"K": 5,
	"L": 1,
	"M": 3,
	"N": 1,
	"O": 1,
	"P": 3,
	"Q": 10,
	"R": 1,
	"S": 1,
	"T": 1,
	"U": 1,
	"V": 4,
	"W": 4,
	"X": 8,
	"Y": 4,
	"Z": 10,
	"*": 0,
}

var tileDistribution = map[string]int{
	"A": 9,
	"B": 2,
	"C": 2,
	"D": 4,
	"E": 12,
	"F": 2,
	"G": 3,
	"H": 2,
	"I": 9,
	"J": 1,
	"K": 1,
	"L": 4,
	"M": 2,
	"N": 6,
	"O": 8,
	"P": 2,
	"Q": 1,
	"R": 6,
	"S": 4,
	"T": 6,
	"U": 4,
	"V": 2,
	"W": 2,
	"X": 1,
	"Y": 2,
	"Z": 1,
	"*": 2,
}

type Theme struct {
	NFColor int
	DLColor int
	TLColor int
	DWColor int
	TWColor int
	CSColor int
	// more to come ...
}

type Config struct {
	SpecialFields [][]int
	LetterScores  map[string]int
	TileDist      map[string]int
	Theme         Theme
}

func NewConfig() *Config {
	return &Config{
		SpecialFields: specialFields,
		LetterScores:  letterScores,
		TileDist:      tileDistribution,
		Theme: Theme{
			NFColor: NFColor,
			DLColor: DLColor,
			TLColor: TLColor,
			DWColor: DWColor,
			TWColor: TWColor,
			CSColor: CSColor,
		},
	}
}

func (c *Config) GetSpecialFields() [][]int {
	return c.SpecialFields
}

func (c *Config) GetLetterScores() map[string]int {
	return c.LetterScores
}

func (c *Config) GetTileDist() map[string]int {
	return c.TileDist
}

func (c *Config) GetTheme() Theme {
	return c.Theme
}

// ReadConfig reads the config file and returns a Config struct using viper
func ReadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return &config
}

// WriteConfig writes the config file using viper
func WriteConfig(config *Config) {
	viper.Set("SpecialFields", config.SpecialFields)
	viper.Set("LetterScores", config.LetterScores)
	viper.Set("TileDist", config.TileDist)
	viper.Set("Theme", config.Theme)

	if err := viper.SafeWriteConfigAs("config.json"); err != nil {
		log.Printf("Failed to write config file: %v", err)
	}
}
