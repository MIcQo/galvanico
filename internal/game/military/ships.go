package military

import (
	"fmt"
	"strings"
)

// Ship represents different ship types using uint
type Ship uint

const (
	Ironclad Ship = iota
	TorpedoBoat
	SteamBattleship
	Gunboat
	Monitor
	HeavyCruiser
	Submarine
	Destroyer
	AircraftCarrier
	SupplyShip
	RocketBarge
)

// shipNames maps Ship values to their string representations
var shipNames = []string{
	"ironclad",
	"torpedo_boat",
	"steam_battleship",
	"gunboat",
	"monitor",
	"heavy_cruiser",
	"submarine",
	"destroyer",
	"aircraft_carrier",
	"supply_ship",
	"rocket_barge",
}

// String returns the string representation of a Ship
func (s Ship) String() string {
	if int(s) >= len(shipNames) {
		return unknown
	}
	return shipNames[s]
}

// ShipFromString converts a string to a Ship type
func ShipFromString(s string) (Ship, error) {
	for i, name := range shipNames {
		if strings.EqualFold(name, s) {
			return Ship(i), nil
		}
	}
	return 0, fmt.Errorf("invalid ship name: %s", s)
}
