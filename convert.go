package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	var spells []OldSpell
	newSpells := make(map[string]NewSpell)

	json.Unmarshal(jsonByte, &spells)

	//formatSpells(spells, newSpells)

	var output []NewSpell
	for _, spell := range newSpells {
		output = append(output, spell)
	}

	writeFile, _ := json.MarshalIndent(output, "", " ")

	ioutil.WriteFile("spells.json", writeFile, 0644)
}

func formatSpells(spells []OldSpell, newSpells map[string]NewSpell) {
	//comp := make(map[string]int)
	errors := make(map[string]string)

	for _, spell := range spells {
		newSpell := NewSpell{}
		newSpell.Name = spell.Name
		newSpell.Link = spell.Link

		{
			words := strings.Split(spell.School, " ")
			for idx, word := range words {
				reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
				word = reg.ReplaceAllString(word, "")
				switch idx {
				case 0:
					newSpell.School.School = word
				case 1:
					newSpell.School.SubSchool = word
				default:
					newSpell.School.Descriptors = append(newSpell.School.Descriptors, word)
				}
			}
		}

		newSpell.Classes = make(map[string]int)
		for class, level := range spell.Classes {
			newSpell.Classes[class], _ = strconv.Atoi(level)
		}

		newSpell.CastingTime = spell.CastingTime

		for _, c := range spell.Components {

			switch c[0:1] {
			case "V":
				newSpell.Components.Verbal = true
			case "S":
				newSpell.Components.Somatic = true
			case "M":
				if strings.Contains(c, "M/DF") {
					c = strings.Replace(c, "M/DF", "M", 1)
					newSpell.Components.DivineFocus = true
				}
				newSpell.Components.Material = c
			case "F":
				newSpell.Components.Focus = c
			case "D":
				newSpell.Components.DivineFocus = true
				if strings.Contains(c, "DF/M") {
					c = strings.Replace(c, "DF/M", "M", 1)
					newSpell.Components.Material = c
				}
			}
		}

		newSpell.Effect.Range = spell.Range
		newSpell.Effect.Area = spell.Area
		newSpell.Effect.Target = spell.Target
		newSpell.Effect.Duration = spell.Duration
		newSpell.Effect.Description = spell.Description

		newSpell.SavingThrow.Description = spell.SavingThrow
		newSpell.SavingThrow.Fortitude = strings.Contains(spell.SavingThrow, "Fort")
		newSpell.SavingThrow.Reflex = strings.Contains(spell.SavingThrow, "Reflex")
		newSpell.SavingThrow.Will = strings.Contains(spell.SavingThrow, "Will")

		newSpell.SpellResistance.Description = spell.SpellResistance
		newSpell.SpellResistance.Applies = strings.Contains(spell.SavingThrow, "Yes")

		newSpell.Description = spell.Description

		newSpells[newSpell.Name] = newSpell
	}

	fmt.Println(errors)
	fmt.Println(newSpells["Geyser"])

}

type NewSpell struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	School struct {
		School      string   `json:"school"`
		SubSchool   string   `json:"subSchool"`
		Descriptors []string `json:"descriptors"`
	} `json:"school"`
	Classes     map[string]int `json:"classes"`
	CastingTime struct {
		Action string `json:"action"`
		Time   string `json:"time"`
	} `json:"castingTime"`
	Components struct {
		Verbal      bool   `json:"verbal"`
		Somatic     bool   `json:"somatic"`
		Material    string `json:"material"`
		Focus       string `json:"focus"`
		DivineFocus bool   `json:"divineFocus"`
	} `json:"components"`
	Effect struct {
		Range       string `json:"range"`
		Area        string `json:"area"`
		Target      string `json:"target"`
		Duration    string `json:"duration"`
		Description string `json:"description"`
	} `json:"effect"`
	SavingThrow struct {
		Fortitude   bool   `json:"fortitude"`
		Reflex      bool   `json:"reflex"`
		Will        bool   `json:"will"`
		Description string `json:"description"`
	} `json:"savingThrow"`
	SpellResistance struct {
		Applies     bool   `json:"applies"`
		Description string `json:"description"`
	} `json:"spellResistance"`
	Description string `json:"description"`
}

type OldSpell struct {
	Name        string            `json:"name"`
	Link        string            `json:"link"`
	School      string            `json:"school"`
	Classes     map[string]string `json:"classes"`
	CastingTime struct {
		Action string `json:"action"`
		Time   string `json:"time"`
	} `json:"castingTime"`
	Components      []string `json:"components"`
	Range           string   `json:"range"`
	Area            string   `json:"area"`
	Target          string   `json:"target"`
	Duration        string   `json:"duration"`
	Effect          string   `json:"effect"`
	SavingThrow     string   `json:"savingThrow"`
	SpellResistance string   `json:"spellResistance"`
	Description     string   `json:"description"`
}
