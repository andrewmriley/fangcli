package main

import (
	"reflect"
	"testing"
)

var mockZones = []zone{
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

func Test_getConfig(t *testing.T) {
	tests := []struct {
		name string
		want configuration
	}{
		{name: "Defaults",
			want: configuration{
				bonus:     "",
				minlevel:  1,
				maxlevel:  500,
				sortlevel: "asc",
				expansion: "",
				quiet:     false,
				url:       "https://fangbreaker.zone/api/bonuses/today",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processZones(t *testing.T) {
	type args struct {
		conf     configuration
		allZones []zone
	}
	defaultConf := getConfig()
	tests := []struct {
		name string
		args args
		want sortedZones
	}{
		{
			name: "Sort_ASC",
			args: args{defaultConf, mockZones},
			want: sortedZones{
				"coin": []zone{
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
				},
				"respawn": []zone{
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
						ID:        3,
						ZoneID:    3,
						Name:      "Corn",
						Expansion: "Blue",
						MinLevel:  20,
						MaxLevel:  45,
						ZoneType:  "Outdoor",
						Bonus:     "respawn",
					},
				},
			},
		},
		{
			name: "Sort_DESC",
			args: args{configuration{sortlevel: "desc", minlevel: 1, maxlevel: 500}, mockZones},
			want: sortedZones{
				"coin": []zone{
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
				},
				"respawn": []zone{
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
						ID:        4,
						ZoneID:    4,
						Name:      "Burger",
						Expansion: "Blue",
						MinLevel:  1,
						MaxLevel:  25,
						ZoneType:  "Indoor",
						Bonus:     "respawn",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processZones(tt.args.conf, tt.args.allZones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processZones() = %v, want %v", got, tt.want)
			}
		})
	}
}
