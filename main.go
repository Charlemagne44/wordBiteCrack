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
	}

	// gathers user input for vertical and horizontal chunks and load them into the game object
	game.LoadWords()
	fmt.Printf("Horizontal chunks: %v\n", game.HorizontalChunks)
	fmt.Printf("Vertical chunks: %v\n", game.VerticalChunks)

	// backgrack to find longest and highest scoring combinations

}
