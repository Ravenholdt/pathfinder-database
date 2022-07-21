package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var spells map[string]Spell

func main() {
	spells = make(map[string]Spell)

	fmt.Println("Hello!")

	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	var tmp []Spell

	//fmt.Println(jsonFile)
	json.Unmarshal(jsonByte, &tmp)

	//fmt.Println(tmp)
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
