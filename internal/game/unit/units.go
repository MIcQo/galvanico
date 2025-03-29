package unit

import (
	"fmt"
	"strings"
)

const unknown = "unknown"

// Unit represents different military units using uint
type Unit uint

const (
	LineInfantry Unit = iota
	ArmoredAutomaton
	BayonetInfantry
	ShockTrooper
	Rifleman
	Sharpshooter
	Dragoon
	ArmoredTrain
	FieldArtillery
	Howitzer
	ReconBiplane
	ZeppelinBomber
	CombatMedic
)

// unitNames maps Unit values to their string representations
var unitNames = []string{
	"line_infantry",
	"armored_automaton",
	"bayonet_infantry",
	"shock_trooper",
	"rifleman",
	"sharpshooter",
	"dragoon",
	"armored_train",
	"field_artillery",
	"howitzer",
	"recon_biplane",
	"zeppelin_bomber",
	"combat_medic",
}

// String returns the string representation of a Unit
func (u Unit) String() string {
	if int(u) >= len(unitNames) {
		return unknown
	}
	return unitNames[u]
}

// UnitFromString converts a string to a Unit type
func UnitFromString(s string) (Unit, error) {
	for i, name := range unitNames {
		if strings.EqualFold(name, s) {
			return Unit(i), nil
		}
	}
	return 0, fmt.Errorf("invalid unit name: %s", s)
}
