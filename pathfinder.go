package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

var spells map[string]Spell

func main() {
	sourceFile := "spells.json"
	argsWithoutProg := os.Args[1:]

	jsonFile, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	var tmp []Spell
	json.Unmarshal(jsonByte, &tmp)

	spells = make(map[string]Spell)
	for _, s := range tmp {
		spells[s.Name] = s
	}

	for _, c := range argsWithoutProg {
		switch c {
		case "save":
			backup(sourceFile)
			save(sourceFile)
		case "denull":
			reader := bufio.NewReader(os.Stdin)
			for _, s := range spells {
				if s.Description != "null" {
					continue
				}

				fmt.Println(s.Name)
				fmt.Println(s.Url)
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				text = strings.Replace(text, "\n", "", -1)
				if text == "" {
					continue
				}
				if text == "e" {
					break
				}

				spell := spells[s.Name]
				spell.Description = text
				spells[s.Name] = spell

				fmt.Println(spells[s.Name].Description)
			}
		}
	}
}

func backup(sourceFile string) {
	currentTime := time.Now()
	input, _ := ioutil.ReadFile(sourceFile)
	destinationFile := sourceFile + "-bkp-" + currentTime.Format("2006-01-02-15:04:05")
	ioutil.WriteFile(destinationFile, input, 0644)
}

func save(sourceFile string) {
	var output []Spell
	for _, spell := range spells {
		output = append(output, spell)
	}

	// Sort
	sort.Slice(output, func(i, j int) bool {
		return output[i].Name < output[j].Name
	})

	// Write
	writeFile, _ := json.MarshalIndent(output, "", " ")
	ioutil.WriteFile("spells.json", writeFile, 0644)
}

type School struct {
	School      string   `json:"school"`
	SubSchool   *string  `json:"sub_school"`
	Descriptors []string `json:"descriptors"`
}

type Components struct {
	Verbal      bool    `json:"verbal"`
	Somatic     bool    `json:"somatic"`
	Material    *string `json:"material"`
	Focus       *string `json:"focus"`
	DivineFocus bool    `json:"divine_focus"`
}

type Effect struct {
	Range       *string `json:"range"`
	Area        *string `json:"area"`
	Target      *string `json:"target"`
	Duration    *string `json:"duration"`
	Dismissible bool    `json:"dismissible"`
}

type SavingThrow struct {
	Fortitude   bool    `json:"fortitude"`
	Reflex      bool    `json:"reflex"`
	Will        bool    `json:"will"`
	Description *string `json:"description"`
}

type SpellResistance struct {
	Applies     bool    `json:"applies"`
	Description *string `json:"description"`
}

type Spell struct {
	Name              string          `json:"name"`
	Url               string          `json:"url"`
	School            School          `json:"school"`
	Classes           map[string]int  `json:"classes"`
	CastingTime       string          `json:"casting_time"`
	Components        Components      `json:"components"`
	Effect            Effect          `json:"effect"`
	SavingThrow       SavingThrow     `json:"saving_throw"`
	SpellResistance   SpellResistance `json:"spell_resistance"`
	Description       string          `json:"description"`
	SourceBook        string          `json:"source_book"`
	RelatedSpellNames []string        `json:"related_spell_names"`
}
