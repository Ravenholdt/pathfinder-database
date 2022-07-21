package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	parseSpells()
}

func parseSpells() {
	file, err := os.Open("spellsToAdd.txt")
	if err != nil {
		fmt.Println(err)
	}

	var spells []Spell

	s := Spell{}

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	findDesc := false
	for scanner.Scan() {
		row := scanner.Text()

		row = strings.Trim(row, " ")

		fmt.Println(row)
		fmt.Println(findDesc)

		if len(row) == 0 {
			continue
		}

		if findDesc {
			row = "DESCRIPTION " + row
			fmt.Println(row)
		}

		fmt.Println(strings.Split(row, " ")[0])

		switch strings.Split(row, " ")[0] {
		case "CASTING":
		case "EFFECT":
		case "Effect":
		default:
			row = "padd " + row //Name removes the first 4 letters
			fallthrough
		case "Name":
			if s.Name != "" {
				spells = append(spells, s)
			}
			s = Spell{Effect: Effect{}}

			s.Name = strings.Trim(row[4:], " ")
		case "School":
			rows := strings.Split(strings.Trim(row, " "), ";")
			row = rows[0]
			words := strings.Split(strings.Trim(row[6:], " "), " ")
			for idx, word := range words {

				fmt.Println(idx)
				fmt.Println(word)

				switch idx {
				case 0:
					s.School.School = word
				case 1:
					if strings.Contains(word, "(") {
						tmp := strings.Trim(word, "()")
						s.School.SubSchool = &tmp
						continue
					}
					fallthrough
				default:
					w := strings.Trim(word, "[]")
					s.School.Descriptors = append(s.School.Descriptors, w)
				}
			}
			if len(rows) > 1 {
				row = strings.Trim(rows[1], " ")
				fmt.Println(row)
				goto splitSchoolAndLevel
			}
			continue
		splitSchoolAndLevel:
			fallthrough

		case "Level":
			row = strings.ReplaceAll(row, "/ ", "/")
			words := strings.Split(strings.Trim(row[5:], " "), ",")
			fmt.Println(words)
			s.Classes = make(map[string]int)
			for _, w := range words {
				w = strings.Trim(w, " ")
				if w == "unchained" {
					w = "unchained summoner"
				}
				fmt.Println(w)
				s.Classes[strings.Split(w, " ")[0]], _ = strconv.Atoi(w[len(w)-1:])
			}
		case "Casting":
			s.CastingTime = strings.Trim(row[12:], " ")
		case "Components":
			comp := strings.Split(row[10:], ",")
			for _, c := range comp {
				c = strings.Trim(c, " ")
				fmt.Println(c)
				if strings.Contains(c, "/DF") {
					c = strings.ReplaceAll(c, "/DF", "")
					s.Components.DivineFocus = true
				}
				switch c[0:1] {
				case "V":
					s.Components.Verbal = true
				case "S":
					s.Components.Somatic = true
				case "M":
					s.Components.Material = &c
				case "F":
					s.Components.Focus = &c
				case "D":
					s.Components.DivineFocus = true
				}
			}
		case "Range":
			tmp := strings.Split(strings.Trim(row[5:], " "), " ")
			s.Effect.Range = &tmp[0]
		case "Target":
			tmp := strings.Split(strings.Trim(row[6:], " "), " ")
			s.Effect.Target = &tmp[0]
		case "Duration":
			tmp := strings.Trim(row[8:], " ")
			s.Effect.Duration = &tmp
			s.Effect.Dismissible = strings.Contains(*s.Effect.Duration, "(D)")
		case "Saving":
			tmp := strings.Trim(row[12:], " ")
			s.SavingThrow.Description = &tmp
			s.SavingThrow.Fortitude = strings.Contains(*s.SavingThrow.Description, "Fort")
			s.SavingThrow.Reflex = strings.Contains(*s.SavingThrow.Description, "Reflex")
			s.SavingThrow.Will = strings.Contains(*s.SavingThrow.Description, "Will")
		case "Spell":
			tmp := strings.Trim(row[16:], " ")
			s.SpellResistance.Description = &tmp
			s.SpellResistance.Applies = strings.Contains(*s.SpellResistance.Description, "yes")
		case "DESCRIPTION":
			fmt.Println("!!!!!!!!!!!!!")
			desc := strings.Trim(row[11:], " ")
			fmt.Println(len(desc))
			if len(desc) == 0 {
				findDesc = true
				fmt.Println(findDesc)
				continue
			} else {
				fmt.Println("WEEEEEEE!!!")
				fmt.Println(len(desc))
			}
			fmt.Println(desc)
			s.Description = strings.Trim(desc, " ")
			findDesc = false
		}

	}

	spells = append(spells, s)

	writeFile, _ := json.MarshalIndent(spells, "", " ")
	ioutil.WriteFile("spellsToAdd.json", writeFile, 0644)
}

func parseCopyright() {
	jsonMap := make(map[string]string)

	csvFile, _ := os.Open("d20-spell-copyright.csv")
	csvReader := csv.NewReader(csvFile)
	data, _ := csvReader.ReadAll()
	for idx := range data {
		if idx == 0 {
			continue
		}

		var copyright string

		if strings.Contains(data[idx][5], "©") {
			copyright = strings.Split(data[idx][5], "©")[0]
		} else if strings.Contains(data[idx][5], "Copyright") {
			copyright = strings.Split(data[idx][5], "Copyright")[0]
		} else {
			if data[idx][5] != "null" {
				copyright = data[idx][5]
			}
		}

		if strings.Contains(copyright, "15") {
			fmt.Println(data[idx][4])
		}

		copyright = strings.Trim(copyright, ". ")

		jsonMap[data[idx][4]] = copyright

		writeFile, _ := json.MarshalIndent(jsonMap, "", " ")

		ioutil.WriteFile("spells-copyright.json", writeFile, 0644)
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
