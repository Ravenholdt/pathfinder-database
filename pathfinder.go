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

	jsonFile, err := os.Open("spells.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	var tmp []Spell

	//fmt.Println(jsonFile)
	json.Unmarshal(jsonByte, &tmp)

	//fmt.Println(tmp)
}

type Spell struct {
	Name        string         `json:"name"`
	Link        string         `json:"link"`
	School      string         `json:"school"`
	Classes     map[string]int `json:"classes"`
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
