package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type zone struct {
	ID        int    `json:"id"`
	ZoneID    int    `json:"zoneId"`
	Name      string `json:"name"`
	Expansion string `json:"expansion"`
	MinLevel  int    `json:"minLevel"`
	MaxLevel  int    `json:"maxLevel"`
	ZoneType  string `json:"zoneType"`
	Bonus     string `json:"bonus"`
}

type configuration struct {
	bonus     *string
	minlevel  *int
	maxlevel  *int
	sortlevel *string
	expansion *string
	quiet     *bool
	mock      *bool
	url       string
}

type LevelSorter []zone

func (a LevelSorter) Len() int           { return len(a) }
func (a LevelSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LevelSorter) Less(i, j int) bool { return a[i].MinLevel < a[j].MinLevel }

func mockZones() []zone {
	return []zone{
		{
			ID:        1,
			ZoneID:    1,
			Name:      "Pizza",
			Expansion: "Green",
			MinLevel:  1,
			MaxLevel:  50,
			ZoneType:  "Indoor",
			Bonus:     "respawn",
		},
		{
			ID:        2,
			ZoneID:    2,
			Name:      "Chicken",
			Expansion: "Blue",
			MinLevel:  1,
			MaxLevel:  20,
			ZoneType:  "Indoor",
			Bonus:     "coin",
		},
		{
			ID:        3,
			ZoneID:    3,
			Name:      "Corn",
			Expansion: "Blue",
			MinLevel:  20,
			MaxLevel:  45,
			ZoneType:  "Outdoor",
			Bonus:     "respawn",
		},
		{
			ID:        4,
			ZoneID:    4,
			Name:      "Burger",
			Expansion: "Blue",
			MinLevel:  1,
			MaxLevel:  25,
			ZoneType:  "Indoor",
			Bonus:     "respawn",
		},
		{
			ID:        5,
			ZoneID:    5,
			Name:      "Hotdog",
			Expansion: "Green",
			MinLevel:  1,
			MaxLevel:  50,
			ZoneType:  "Outdoor",
			Bonus:     "none",
		},
	}
}

func getJson(conf configuration, zones *[]zone) {
	myClient := &http.Client{Timeout: 10 * time.Second}
	if !*conf.quiet {
		fmt.Print("Fetching from ", conf.url, "...")
	}
	r, err := myClient.Get(conf.url)
	if err != nil {
		if !*conf.quiet {
			fmt.Println(" Error.")
		}
		log.Fatalln(err)
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(zones)
	if !*conf.quiet {
		fmt.Println(" Done!")
	}
}

func splitFilterZones(zones []zone, conf configuration) map[string][]zone {
	bonusZonesByType := make(map[string][]zone)
	for _, zone := range zones {
		if zone.MinLevel >= *conf.minlevel && zone.MaxLevel <= *conf.maxlevel && (*conf.expansion == "" || strings.ToLower(*conf.expansion) == strings.ToLower(zone.Expansion)) {
			bonusZonesByType[zone.Bonus] = append(bonusZonesByType[zone.Bonus], zone)
		}
	}
	return bonusZonesByType
}

func getConfig() configuration {
	var conf configuration
	conf.bonus = flag.String("bonus", "", "What type of bonus to filter on (experience, loot etc) Empty for all")
	conf.minlevel = flag.Int("minlevel", 1, "Minimum level")
	conf.maxlevel = flag.Int("maxlevel", 500, "Maximum level")
	conf.sortlevel = flag.String("sortbylevel", "asc", "Sort (ASC or DESC)")
	conf.expansion = flag.String("expansion", "", "Expansion name")
	conf.quiet = flag.Bool("quiet", false, "Suppress header messages")
	conf.mock = flag.Bool("mock", false, "Use mocked data instead of live JSON.")
	conf.url = "https://fangbreaker.zone/api/bonuses/today"
	flag.Parse()
	return conf
}

func main() {
	conf := getConfig()

	var allZones []zone
	if *conf.mock {
		allZones = mockZones()
	} else {
		getJson(conf, &allZones)
	}
	zonesByType := splitFilterZones(allZones, conf)

	displayCount := 0
	for key, zones := range zonesByType {
		if key != "none" && key != "unconfirmed" && (*conf.bonus == "" || key == *conf.bonus) {
			if *conf.sortlevel == "desc" {
				sort.Sort(sort.Reverse(LevelSorter(zones)))
			} else {
				sort.Sort(LevelSorter(zones))
			}
			fmt.Println("\n---", key, "---")
			for _, zone := range zones {
				fmt.Println(zone.Name, zone.Expansion, zone.MinLevel, zone.MaxLevel)
				displayCount++
			}
		}
	}

	if displayCount == 0 && !*conf.quiet {
		fmt.Println("\nNo zones found. Check your filters.")
	}
	if !*conf.quiet {
		fmt.Println("\nThis app is not affiliated with the site but please support it at https://fangbreaker.zone")
	}

}
