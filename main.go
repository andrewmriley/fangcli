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
	bonus     string
	minlevel  int
	maxlevel  int
	sortlevel string
	expansion string
	zonetype  string
	quiet     bool
	url       string
}

type sortedZones map[string][]zone

type LevelSorter []zone

func (a LevelSorter) Len() int           { return len(a) }
func (a LevelSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LevelSorter) Less(i, j int) bool { return a[i].MinLevel < a[j].MinLevel }

func getJson(conf configuration) []zone {
	myClient := &http.Client{Timeout: 10 * time.Second}
	if !conf.quiet {
		fmt.Print("Fetching from ", conf.url, "...")
	}
	r, err := myClient.Get(conf.url)
	if err != nil {
		if !conf.quiet {
			fmt.Println(" Error.")
		}
		log.Fatalln(err)
	}
	defer r.Body.Close()
	var zones []zone
	json.NewDecoder(r.Body).Decode(&zones)
	if zones == nil {
		log.Fatalln("Something went wrong, no zones were decoded.")
	}
	if !conf.quiet {
		fmt.Println(" Done!")
	}
	return zones
}

func compareLower(value string, confExpected string) bool {
	if confExpected == "" {
		return true
	}
	return strings.ToLower(value) == strings.ToLower(confExpected)
}

func processZones(conf configuration, allZones []zone) sortedZones {
	bonusZonesByType := make(sortedZones)
	for _, zone := range allZones {
		if zone.Bonus == "none" {
			continue
		}
		if zone.MinLevel >= conf.minlevel && zone.MaxLevel <= conf.maxlevel && compareLower(zone.Expansion, conf.expansion) && compareLower(zone.ZoneType, conf.zonetype) {
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
	var conf = configuration{
		url: "https://fangbreaker.zone/api/bonuses/today",
	}
	flag.StringVar(&conf.bonus, "bonus", "", "What type of bonus to filter on (experience, loot etc) Empty for all")
	flag.IntVar(&conf.minlevel, "minlevel", 1, "Minimum level")
	flag.IntVar(&conf.maxlevel, "maxlevel", 500, "Maximum level")
	flag.StringVar(&conf.sortlevel, "sortbylevel", "asc", "Sort (ASC or DESC)")
	flag.StringVar(&conf.expansion, "expansion", "", "Expansion name")
	flag.StringVar(&conf.zonetype, "zonetype", "", "Indoor or outdoor")
	flag.BoolVar(&conf.quiet, "quiet", false, "Suppress header messages")
	flag.Parse()
	return conf
}

func displayZones(conf configuration, zonesByType sortedZones) {
	displayCount := 0
	for key, zones := range zonesByType {
		if key != "none" && key != "unconfirmed" && compareLower(key, conf.bonus) {
			fmt.Println("\n---", key, "---")
			for _, zone := range zones {
				fmt.Println(zone.Name, zone.Expansion, zone.MinLevel, zone.MaxLevel)
				displayCount++
			}
		}
	}

	if displayCount == 0 && !conf.quiet {
		fmt.Println("\nNo zones found. Check your filters.")
	}
}

func main() {
	conf := getConfig()
	allZones := getJson(conf)
	zonesByType := processZones(conf, allZones)
	displayZones(conf, zonesByType)

	if !conf.quiet {
		fmt.Println("\nThis app is not affiliated with the site but please support it at https://fangbreaker.zone")
	}
}
