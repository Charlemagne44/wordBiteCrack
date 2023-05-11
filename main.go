package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"wordBiteCrack/trie"
)

//go:embed resources/scrabble.json
var f embed.FS

type Game struct {
	Trie             trie.Trie
	HorizontalChunks []string
	VerticalChunks   []string
	SingleChunks     []string
	ValidWords       []string
}

func (g *Game) LoadTestWords() {
	g.VerticalChunks = []string{"ta", "ng"}
	g.HorizontalChunks = []string{"se", "fo", "iv"}
	g.SingleChunks = []string{"i", "l", "b", "a", "c"}
}

func (g *Game) LoadWords() {
	fmt.Println("Enter vertical chunks followed by enter per vertical chunk")
	fmt.Println("Hit enter with no characters to conlude the vertical chunk entry")
	fmt.Println()

	verticalEntry := true
	i := 1
	for verticalEntry {
		var entry string
		fmt.Printf("Vertical word %d: ", i)
		fmt.Scanf("%s", &entry)
		if entry != "" {
			g.VerticalChunks = append(g.VerticalChunks, entry)
		} else {
			verticalEntry = false
		}
		fmt.Println()
		i += 1
	}

	fmt.Println()
	fmt.Println("Enter horizontal chunks followed by enter per horizontal chunk")
	fmt.Println("Hit enter with no characters to conlude the horizontal chunk entry")
	fmt.Println()

	horizontalEntry := true
	i = 1
	for horizontalEntry {
		var entry string
		fmt.Printf("Horizontal word %d: ", i)
		fmt.Scanf("%s", &entry)
		if entry != "" {
			g.HorizontalChunks = append(g.HorizontalChunks, entry)
		} else {
			horizontalEntry = false
		}
		i += 1
	}
}

func (g *Game) Backtrack(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	// if word is horizontal
	if chunk_orientation == 'h' {
		// try all horizontal chunks before the beginning and end
		for _, horiz_chunk := range g.HorizontalChunks {

			// insert chunk before and check validity
			new_word := horiz_chunk + chunk
			// append a valid scoring word to the games valid words
			if g.Trie.Search(new_word) && len(new_word) >= 3 {
				g.ValidWords = append(g.ValidWords, new_word)
			}
			// is its valid, and more could be constructed -> recurse
			if g.Trie.ValidPath(new_word) {
				// create a new limited list of horizontal chunks excluding the one used and recurse
				g.Backtrack(new_word, remove(remainingHorizontalChunks, horiz_chunk),
					remainingVerticalChunks, remainingSingleChunks, 'h')
			}

			// insert chunk after and check validity
			new_word = chunk + horiz_chunk
			// append a valid scoring word to the games valid words
			if g.Trie.Search(new_word) && len(new_word) >= 3 {
				g.ValidWords = append(g.ValidWords, new_word)
			}
			// if its valid, and more could be constructed -> recurse
			if g.Trie.ValidPath(new_word) {
				// create a new limited list of horizontal chunks excluding the one used and recurse
				g.Backtrack(new_word, remove(remainingHorizontalChunks, horiz_chunk),
					remainingVerticalChunks, remainingSingleChunks, 'h')
			}
		}

		// try all vertical chunks before and after each letter

		// try all single chunks before and after chunk, and before and after each letter

	} else { // if word is vertical

		// try all vertical chunks before the beginning and the end

		// try all horizontal chunks before and after each letter

		// try all single chunks before and after chunk, and before and after each letter

	}

}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func main() {
	// load in the english dict into a trie
	data, _ := f.ReadFile("resources/scrabble.json")

	var dictionary []string
	err := json.Unmarshal(data, &dictionary)
	if err != nil {
		fmt.Printf("Unmarshal %v\n", err)
	}

	trie := trie.InitTrie()
	for _, word := range dictionary {
		if !strings.Contains(word, "-") {
			trie.Insert(strings.ToLower(word))
		}
	}

	// create game object and load in the english dict trie
	game := Game{
		Trie:             *trie,
		VerticalChunks:   make([]string, 0),
		HorizontalChunks: make([]string, 0),
		SingleChunks:     make([]string, 0),
		ValidWords:       make([]string, 0),
	}

	// gathers user input for vertical and horizontal chunks and load them into the game object
	// game.LoadWords()
	game.LoadTestWords()
	fmt.Printf("Horizontal chunks: %v\n", game.HorizontalChunks)
	fmt.Printf("Vertical chunks: %v\n", game.VerticalChunks)
	fmt.Printf("Single chunks: %v\n", game.SingleChunks)

	// create the initial copy of the chunksthat can be removed during recursion
	remainingHorizontalChunks := game.HorizontalChunks
	remainingVerticalChunks := game.VerticalChunks
	remainingSingleChunks := game.SingleChunks

	// backtrack to find longest and highest scoring combinations
	for _, horizontal_chunk := range game.HorizontalChunks {
		game.Backtrack(horizontal_chunk, remove(game.HorizontalChunks, horizontal_chunk),
			remainingVerticalChunks, remainingSingleChunks, 'h')
	}

	for _, vertical_chunk := range game.VerticalChunks {
		game.Backtrack(vertical_chunk, remainingHorizontalChunks,
			remove(game.VerticalChunks, vertical_chunk), remainingSingleChunks, 'v')
	}

	fmt.Printf("Valid words: %v\n", game.ValidWords)

	// fmt.Println(trie.Search("sese"))
	// fmt.Println(trie.ValidPath("ashen"))

}
