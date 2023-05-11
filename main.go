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
	Trie                 trie.Trie
	HorizontalChunks     []string
	VerticalChunks       []string
	SingleChunks         []string
	ValidHorizontalWords []string
	ValidVerticalWords   []string
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

func (g *Game) HorizontalxHorizontal(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	for _, horiz_chunk := range remainingHorizontalChunks {

		// insert chunk before and check validity
		new_word := horiz_chunk + chunk
		// append a valid scoring word to the games valid words
		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
			g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
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
		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
			g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
		}
		// if its valid, and more could be constructed -> recurse
		if g.Trie.ValidPath(new_word) {
			// create a new limited list of horizontal chunks excluding the one used and recurse
			g.Backtrack(new_word, remove(remainingHorizontalChunks, horiz_chunk),
				remainingVerticalChunks, remainingSingleChunks, 'h')
		}
	}
}

func (g *Game) VerticalxVertical(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	for _, vert_chunk := range remainingVerticalChunks {

		// insert chunk before and check validity
		new_word := vert_chunk + chunk
		// append a valid scoring word to the games valid words
		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidVerticalWords, new_word) {
			g.ValidHorizontalWords = append(g.ValidVerticalWords, new_word)
		}
		// is its valid, and more could be constructed -> recurse
		if g.Trie.ValidPath(new_word) {
			// create a new limited list of vertical chunks excluding the one used and recurse
			g.Backtrack(new_word, remainingHorizontalChunks,
				remove(remainingVerticalChunks, vert_chunk), remainingSingleChunks, 'v')
		}

		// insert chunk after and check validity
		new_word = chunk + vert_chunk
		// append a valid scoring word to the games valid words
		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidVerticalWords, new_word) {
			g.ValidVerticalWords = append(g.ValidVerticalWords, new_word)
		}
		// if its valid, and more could be constructed -> recurse
		if g.Trie.ValidPath(new_word) {
			// create a new limited list of horizontal chunks excluding the one used and recurse
			g.Backtrack(new_word, remainingHorizontalChunks,
				remove(remainingVerticalChunks, vert_chunk), remainingSingleChunks, 'v')
		}
	}
}

func (g *Game) HorizontalxVertical(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	// for each letter in the horizontal word, try each vertical chunk above and below each letter (result vertical)
	for _, letter := range chunk {
		for _, vert_chunk := range remainingVerticalChunks {
			// insert before letter (vertically)
			new_word := vert_chunk + string(letter)

			if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidVerticalWords, new_word) {
				g.ValidVerticalWords = append(g.ValidVerticalWords, new_word)
			}

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remove(remainingVerticalChunks, vert_chunk),
					remainingSingleChunks, 'v')
			}

			// insert after letter (vertically)
			new_word = string(letter) + vert_chunk

			if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidVerticalWords, new_word) {
				g.ValidVerticalWords = append(g.ValidVerticalWords, new_word)
			}

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remove(remainingVerticalChunks, vert_chunk),
					remainingSingleChunks, 'v')
			}
		}
	}

	// for each letter in each vertical chunk, try it before and after the horizontal chunk (result horizontal)
	for _, vert_chunk := range remainingVerticalChunks {
		for _, letter := range vert_chunk {
			// insert letter before chunk
			new_word := string(letter) + chunk

			if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
				g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
			}

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remove(remainingVerticalChunks, vert_chunk),
					remainingSingleChunks, 'h')
			}

			// insert letter after chunk
			new_word = vert_chunk + string(letter)

			if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
				g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
			}

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remove(remainingVerticalChunks, vert_chunk),
					remainingSingleChunks, 'h')
			}
		}
	}
}

func (g *Game) HorizontalxSingle(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	// try each single before and after the horizontal chunk (horizontal result)
	for _, single := range remainingSingleChunks {
		new_word := single + chunk

		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
			g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
		}

		if g.Trie.ValidPath(new_word) {
			g.Backtrack(new_word, remainingHorizontalChunks, remainingVerticalChunks,
				remove(remainingSingleChunks, single), 'h')
		}

		new_word = chunk + single

		if g.Trie.Search(new_word) && len(new_word) >= 3 && !contains(g.ValidHorizontalWords, new_word) {
			g.ValidHorizontalWords = append(g.ValidHorizontalWords, new_word)
		}

		if g.Trie.ValidPath(new_word) {
			g.Backtrack(new_word, remainingHorizontalChunks, remainingVerticalChunks,
				remove(remainingSingleChunks, single), 'h')
		}
	}

	// try each single before and after each letter of the horizontal chunk (vertical result)
	for _, single := range remainingSingleChunks {
		for _, letter := range chunk {
			new_word := single + string(letter)

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remainingVerticalChunks,
					remove(remainingSingleChunks, single), 'v')
			}

			new_word = string(letter) + single

			if g.Trie.ValidPath(new_word) {
				g.Backtrack(new_word, remainingHorizontalChunks, remainingVerticalChunks,
					remove(remainingSingleChunks, single), 'v')
			}
		}
	}
}

func (g *Game) Backtrack(chunk string, remainingHorizontalChunks, remainingVerticalChunks,
	remainingSingleChunks []string, chunk_orientation rune) {

	// if word is horizontal
	if chunk_orientation == 'h' {
		// as there are several unique permutations for horizontal x vertical x single they all will contain
		// their own functions as they are all quite long, and we want to avoid an untestable single service
		// function

		// try all horizontal chunks before the beginning and end
		g.HorizontalxHorizontal(chunk, remainingHorizontalChunks, remainingVerticalChunks,
			remainingSingleChunks, chunk_orientation)

		// try all vertical chunks before and after each letter
		g.HorizontalxVertical(chunk, remainingHorizontalChunks, remainingVerticalChunks,
			remainingSingleChunks, chunk_orientation)

		// try all single chunks before and after each chunk, and before and after each letter
		g.HorizontalxSingle(chunk, remainingHorizontalChunks, remainingVerticalChunks,
			remainingSingleChunks, chunk_orientation)

	} else if chunk_orientation == 'v' {

		// try all vertical chunks before the beginning and the end4
		g.VerticalxVertical(chunk, remainingHorizontalChunks, remainingVerticalChunks,
			remainingSingleChunks, 'v')

		// try all horizontal chunks before and after each letter

		// try all single chunks before and after chunk, and before and after each letter

	} else { //  single chunk

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

func contains(slice []string, word string) bool {
	for _, element := range slice {
		if element == word {
			return true
		}
	}
	return false
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
		Trie:                 *trie,
		VerticalChunks:       make([]string, 0),
		HorizontalChunks:     make([]string, 0),
		SingleChunks:         make([]string, 0),
		ValidHorizontalWords: make([]string, 0),
		ValidVerticalWords:   make([]string, 0),
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

	for _, single_chunk := range game.SingleChunks {
		game.Backtrack(single_chunk, remainingHorizontalChunks,
			remainingVerticalChunks, remove(remainingSingleChunks, single_chunk), 's')
	}

	fmt.Printf("Valid horizontal words: %v\n", game.ValidHorizontalWords)
	fmt.Printf("Valid vertical words: %v\n", game.ValidVerticalWords)

	// fmt.Println(trie.Search("sese"))
	// fmt.Println(trie.ValidPath("ashen"))

}
