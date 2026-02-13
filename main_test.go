package main

import (
	"reflect"
	"slices"
	"testing"
)

var mockZones = []zone{
	{
		Name:      "Pizza",
		Expansion: uniqueAppend("green", &ExpansionLookup),
		MinLevel:  1,
		MaxLevel:  50,
		ZoneType:  uniqueAppend("Indoor", &ZoneLookup),
		Bonus:     uniqueAppend("respawn", &BonusLookup),
	},
	{
		Name:      "Chicken",
		Expansion: uniqueAppend("blue", &ExpansionLookup),
		MinLevel:  1,
		MaxLevel:  20,
		ZoneType:  uniqueAppend("Indoor", &ZoneLookup),
		Bonus:     uniqueAppend("coin", &BonusLookup),
	},
	{
		Name:      "Corn",
		Expansion: uniqueAppend("blue", &ExpansionLookup),
		MinLevel:  20,
		MaxLevel:  45,
		ZoneType:  uniqueAppend("Outdoor", &ZoneLookup),
		Bonus:     uniqueAppend("respawn", &BonusLookup),
	},
	{
		Name:      "Burger",
		Expansion: uniqueAppend("blue", &ExpansionLookup),
		MinLevel:  1,
		MaxLevel:  25,
		ZoneType:  uniqueAppend("Indoor", &ZoneLookup),
		Bonus:     uniqueAppend("respawn", &BonusLookup),
	},
	{
		Name:      "Hotdog",
		Expansion: uniqueAppend("green", &ExpansionLookup),
		MinLevel:  1,
		MaxLevel:  50,
		ZoneType:  uniqueAppend("Outdoor", &ZoneLookup),
		Bonus:     uniqueAppend("none", &BonusLookup),
	},
}

func Test_getConfig(t *testing.T) {
	tests := []struct {
		name string
		want configuration
	}{
		{name: "Defaults",
			want: configuration{
				sortlevel: "asc",
				bonus:     -1,
				minlevel:  1,
				maxlevel:  255,
				expansion: -1,
				zonetype:  -1,
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
	tests := []struct {
		name string
		args args
		want sortedZones
	}{
		{
			name: "Sort_ASC",
			args: args{configuration{minlevel: 1, maxlevel: 255, expansion: -1, zonetype: -1}, mockZones},
			want: sortedZones{
				int8(slices.Index(BonusLookup, "coin")): []zone{
					{
						Name:      "Chicken",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  1,
						MaxLevel:  20,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "coin")),
					},
				},
				int8(slices.Index(BonusLookup, "respawn")): []zone{
					{
						Name:      "Pizza",
						Expansion: int8(slices.Index(ExpansionLookup, "green")),
						MinLevel:  1,
						MaxLevel:  50,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
					{
						Name:      "Burger",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  1,
						MaxLevel:  25,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
					{
						Name:      "Corn",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  20,
						MaxLevel:  45,
						ZoneType:  int8(slices.Index(ZoneLookup, "outdoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
				},
			},
		},
		{
			name: "Sort_DESC",
			args: args{configuration{sortlevel: "desc", minlevel: 1, maxlevel: 255, expansion: -1, zonetype: -1}, mockZones},
			want: sortedZones{
				int8(slices.Index(BonusLookup, "coin")): []zone{
					{
						Name:      "Chicken",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  1,
						MaxLevel:  20,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "coin")),
					},
				},
				int8(slices.Index(BonusLookup, "respawn")): []zone{
					{
						Name:      "Corn",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  20,
						MaxLevel:  45,
						ZoneType:  int8(slices.Index(ZoneLookup, "outdoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
					{
						Name:      "Pizza",
						Expansion: int8(slices.Index(ExpansionLookup, "green")),
						MinLevel:  1,
						MaxLevel:  50,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
					{
						Name:      "Burger",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  1,
						MaxLevel:  25,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
				},
			},
		},
		{
			name: "Expansion",
			args: args{configuration{expansion: int8(slices.Index(ExpansionLookup, "green")), minlevel: 1, maxlevel: 255}, mockZones},
			want: sortedZones{
				int8(slices.Index(BonusLookup, "respawn")): []zone{
					{
						Name:      "Pizza",
						Expansion: int8(slices.Index(ExpansionLookup, "green")),
						MinLevel:  1,
						MaxLevel:  50,
						ZoneType:  int8(slices.Index(ZoneLookup, "indoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
					},
				},
			},
		},
		{
			name: "Zonetype",
			args: args{configuration{zonetype: int8(slices.Index(ZoneLookup, "outdoor")), minlevel: 1, maxlevel: 255, expansion: -1}, mockZones},
			want: sortedZones{
				int8(slices.Index(BonusLookup, "respawn")): []zone{
					{
						Name:      "Corn",
						Expansion: int8(slices.Index(ExpansionLookup, "blue")),
						MinLevel:  20,
						MaxLevel:  45,
						ZoneType:  int8(slices.Index(ZoneLookup, "outdoor")),
						Bonus:     int8(slices.Index(BonusLookup, "respawn")),
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
