package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

func main() {
	sourceFile := "spells.json"
	backup(sourceFile)
	jsonFile, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	var spells []Spell

	json.Unmarshal(jsonByte, &spells)

	spellList := make(map[string]Spell)
	for _, spell := range spells {
		spellList[spell.Name] = spell
	}

	//formatSpells(spells, newSpells)
	//addCopyright(spells, newSpells)
	//addClass(spells, newSpells)
	//updateOldToNew(spells, newSpells)

	for name, _ := range spellList {
		for c, level := range spellList[name].Classes {
			switch c {
			case "sorcerer/wizard":
				spellList[name].Classes["sorcerer"] = level
				spellList[name].Classes["wizard"] = level
				delete(spellList[name].Classes, "sorcerer/wizard")

			case "summoner/unchained":
				spellList[name].Classes["summoner"] = level
				spellList[name].Classes["unchained_summoner"] = level
				delete(spellList[name].Classes, "summoner/unchained")
			}

		}
	}

	var output []Spell
	for _, spell := range spellList {
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

func backup(sourceFile string) {
	currentTime := time.Now()
	input, _ := ioutil.ReadFile(sourceFile)
	destinationFile := sourceFile + "-bkp-" + currentTime.Format("2006-01-02-15:04:05")
	ioutil.WriteFile(destinationFile, input, 0644)
}

func updateOldToNew(spells []OldSpell, newSpells map[string]Spell) {
	for _, old := range spells {
		newSpell := copySpell(old)
		newSpells[newSpell.Name] = newSpell
	}
}

func nilOrString(in string) *string {
	if in == "" {
		return nil
	}
	return &in
}

func copySpell(old OldSpell) Spell {
	return Spell{
		Name: old.Name,
		Url:  old.Link,
		School: School{
			School:      old.School.School,
			SubSchool:   nilOrString(old.School.SubSchool),
			Descriptors: old.School.Descriptors,
		},
		Classes:     old.Classes,
		CastingTime: old.CastingTime.Time + " " + old.CastingTime.Unit,
		Components: Components{
			Verbal:      old.Components.Verbal,
			Somatic:     old.Components.Somatic,
			Material:    nilOrString(old.Components.Material),
			Focus:       nilOrString(old.Components.Focus),
			DivineFocus: old.Components.DivineFocus,
		},
		Effect: Effect{
			Range:       nilOrString(old.Effect.Range),
			Area:        nilOrString(old.Effect.Area),
			Target:      nilOrString(old.Effect.Target),
			Duration:    nilOrString(old.Effect.Duration),
			Dismissible: old.Effect.Dismissible,
		},
		SavingThrow: SavingThrow{
			Fortitude:   old.SavingThrow.Fortitude,
			Reflex:      old.SavingThrow.Reflex,
			Will:        old.SavingThrow.Will,
			Description: nilOrString(old.SavingThrow.Description),
		},
		SpellResistance: SpellResistance{
			Applies:     old.SpellResistance.Applies,
			Description: nilOrString(old.SpellResistance.Description),
		},
		Description:       old.Description,
		SourceBook:        old.SourceBook,
		RelatedSpellNames: nil,
	}
}

func addClass(spells []OldSpell, newSpells map[string]Spell) {
	jsonFile, err := os.Open("class-spells.json")
	if err != nil {
		fmt.Println(err)
	}

	type nC struct {
		Name  string
		Level int
		Class string
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)
	var newClass []nC
	json.Unmarshal(jsonByte, &newClass)

	fmt.Println(newClass)

	for _, old := range spells {
		newSpell := copySpell(old)
		newSpells[newSpell.Name] = newSpell
	}

	for _, c := range newClass {
		if _, ok := newSpells[c.Name]; ok {
			newSpells[c.Name].Classes[c.Class] = c.Level
		} else {
			fmt.Println(c.Name)
		}

	}

	//fmt.Println(newSpells)
}

func addCopyright(spells []OldSpell, newSpells map[string]Spell) {
	jsonFile, err := os.Open("spells-copyright.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonByte, _ := ioutil.ReadAll(jsonFile)
	var source = make(map[string]string)
	json.Unmarshal(jsonByte, &source)

	for _, old := range spells {
		newSpell := copySpell(old)
		newSpell.SourceBook = source[newSpell.Name]
		newSpells[newSpell.Name] = newSpell
	}
}

/*
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
*/

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

type OldSpell struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	School struct {
		School      string   `json:"school"`
		SubSchool   string   `json:"subSchool"`
		Descriptors []string `json:"descriptors"`
	} `json:"school"`
	Classes     map[string]int `json:"classes"`
	CastingTime struct {
		Unit string `json:"unit"`
		Time string `json:"time"`
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
		Dismissible bool   `json:"dismissible"`
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
	Description       string   `json:"description"`
	SourceBook        string   `json:"sourceBook"`
	RelatedSpellNames []string `json:"relatedSpellNames"`
}
