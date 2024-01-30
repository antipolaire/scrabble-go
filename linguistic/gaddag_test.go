package linguistic

import (
	"github.com/adject1/macondo/gaddagmaker"
	"testing"
)

func TestGaddag(t *testing.T) {
	t.Run("Test", func(t *testing.T) {
		dictionary_filepath := "../assets/dicts/en.csv"
		gd := gaddagmaker.GenerateGaddag(dictionary_filepath, false, false)
		gd.SerializeElements()
		gd.Alphabet.Val('a')
	})

}
