package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"
)

var ZoneLookup = []string{"indoor", "outdoor"}
var BonusLookup = []string{"aa", "coin", "experience", "faction", "loot", "none", "rare", "respawn", "skill", "unconfirmed"}

// Not prepopulating this slice to avoid compiling in IP.
var ExpansionLookup []string

type zone struct {
	Name      string
	Expansion int8
	ZoneType  int8
	Bonus     int8
	MinLevel  uint8
	MaxLevel  uint8
}

type configuration struct {
	sortlevel string
	bonus     int8
	minlevel  uint8
	maxlevel  uint8
	expansion int8
	zonetype  int8
}

type sortedZones map[int8][]zone

type LevelSorter []zone

func (a LevelSorter) Len() int           { return len(a) }
func (a LevelSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LevelSorter) Less(i, j int) bool { return a[i].MinLevel < a[j].MinLevel }

func uniqueAppend(val string, l *[]string) int8 {
	lowval := strings.ToLower(val)
	i := slices.Index(*l, lowval)
	if i == -1 {
		*l = append(*l, lowval)
		i = len(*l) - 1
	}
	return int8(i)
}

func (z *zone) UnmarshalJSON(data []byte) error {
	var temp struct {
		Name      string `json:"name"`
		Expansion string `json:"expansion"`
		ZoneType  string `json:"zoneType"`
		Bonus     string `json:"bonus"`
		MinLevel  uint8  `json:"minLevel"`
		MaxLevel  uint8  `json:"maxLevel"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	z.Name = temp.Name
	z.Expansion = uniqueAppend(temp.Expansion, &ExpansionLookup)
	z.ZoneType = uniqueAppend(temp.ZoneType, &ZoneLookup)
	z.Bonus = uniqueAppend(temp.Bonus, &BonusLookup)
	z.MinLevel = temp.MinLevel
	z.MaxLevel = temp.MaxLevel
	return nil
}

func getJson() []zone {
	const url = "https://fangbreaker.zone/api/bonuses/today"
	myClient := &http.Client{Timeout: 10 * time.Second}
	fmt.Print("Fetching from ", url, "...")
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println(" Error.")
		log.Fatalln(err)
	}
	defer r.Body.Close()
	var zones []zone
	json.NewDecoder(r.Body).Decode(&zones)
	if zones == nil {
		log.Fatalln("Something went wrong, no zones were decoded.")
	}
	fmt.Println(" Done!")
	return zones
}

func processZones(conf configuration, allZones []zone) sortedZones {
	bonusZonesByType := make(sortedZones)
	none := int8(slices.Index(BonusLookup, "none"))
	for _, zone := range allZones {
		if zone.Bonus == none {
			continue
		}
		checkLevel := zone.MinLevel >= conf.minlevel && zone.MinLevel <= conf.maxlevel
		checkExpansion := conf.expansion == -1 || conf.expansion == zone.Expansion
		checkZoneType := conf.zonetype == -1 || conf.zonetype == zone.ZoneType
		if checkLevel && checkExpansion && checkZoneType {
			bonusZonesByType[zone.Bonus] = append(bonusZonesByType[zone.Bonus], zone)
		}
	}
	for _, zones := range bonusZonesByType {
		if conf.sortlevel == "desc" {
			sort.Sort(sort.Reverse(LevelSorter(zones)))
		} else {
			sort.Sort(LevelSorter(zones))
		}
	}
	return bonusZonesByType
}

func getConfig() configuration {
	var bonus string
	var zonetype string
	var expansion string
	var sortlevel string
	var minlevel int
	var maxlevel int
	flag.StringVar(&bonus, "bonus", "", "What type of bonus to filter on (experience, loot etc) Empty for all")
	flag.IntVar(&minlevel, "minlevel", 1, "Minimum level")
	flag.IntVar(&maxlevel, "maxlevel", 255, "Maximum level")
	flag.StringVar(&sortlevel, "sortbylevel", "asc", "Sort (ASC or DESC)")
	flag.StringVar(&expansion, "expansion", "", "Expansion name")
	flag.StringVar(&zonetype, "zonetype", "", "Indoor or outdoor")
	flag.Parse()
	return configuration{
		sortlevel: strings.ToLower(sortlevel),
		bonus:     int8(slices.Index(BonusLookup, strings.ToLower(bonus))),
		expansion: int8(slices.Index(ExpansionLookup, strings.ToLower(expansion))),
		zonetype:  int8(slices.Index(ZoneLookup, strings.ToLower(zonetype))),
		minlevel:  uint8(minlevel),
		maxlevel:  uint8(maxlevel),
	}
}

func displayZones(conf configuration, zonesByType sortedZones) {
	displayed := false
	none := int8(slices.Index(BonusLookup, "none"))
	unconfirmed := int8(slices.Index(BonusLookup, "unconfirmed"))
	for key, zones := range zonesByType {
		if key != none && key != unconfirmed && (conf.bonus == -1 || key == conf.bonus) {
			fmt.Println("\n---", strings.ToUpper(BonusLookup[key]), "---")
			for _, zone := range zones {
				fmt.Println(zone.Name, ExpansionLookup[zone.Expansion], zone.MinLevel, zone.MaxLevel)
				displayed = true
			}
		}
	}

	if !displayed {
		fmt.Println("\nNo zones found. Check your filters.")
	}
}

func main() {
	allZones := getJson()
	conf := getConfig()
	zonesByType := processZones(conf, allZones)
	displayZones(conf, zonesByType)

	fmt.Println("\nThis app is not affiliated with the site but please support it at https://fangbreaker.zone")
}
