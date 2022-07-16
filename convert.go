package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
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

	formatSpells(spellList)
	//addCopyright(spells, spellList)
	//addClass(spells, spellList)
	//updateOldToNew(spells, spellList)

	verifyClasses(spellList)

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

func verifyClasses(spellList map[string]Spell) {
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

func formatSpells(spells map[string]Spell) {
	/*for spell := range spells {
		tmpSpell := spells[spell]
		if tmpSpell.SpellResistance.Description != nil {
			tmpSpell.SpellResistance.Applies = strings.Contains(*tmpSpell.SpellResistance.Description, "yes")
		}
		spells[spell] = tmpSpell
	}*/

	keywords := make(map[string][]string)
	keywords["abjuration"] = []string{}
	keywords["conjuration"] = []string{"calling", "creation", "healing", "summoning", "teleportation"}
	keywords["divination"] = []string{"scrying"}
	keywords["enchantment"] = []string{"charm", "compulsion"}
	keywords["evocation"] = []string{}
	keywords["illusion"] = []string{"figment", "glamer", "pattern", "phantasm", "shadow"}
	keywords["necromancy"] = []string{}
	keywords["transmutation"] = []string{"polymorph"}

	for spell := range spells {
		falseSubschool := true
		tmpSpell := spells[spell]
		if tmpSpell.School.SubSchool != nil && !strings.Contains(*tmpSpell.School.SubSchool, "or") {
			for _, val := range keywords[tmpSpell.School.School] {
				if val == *tmpSpell.School.SubSchool {
					falseSubschool = false
					break
				}
			}

			if falseSubschool {
				desc := *tmpSpell.School.SubSchool
				tmpSpell.School.Descriptors = append(tmpSpell.School.Descriptors, desc)
				tmpSpell.School.SubSchool = nil
			}
		}

		spells[spell] = tmpSpell
	}
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
