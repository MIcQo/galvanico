package building

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

type Building uint

// String returns the string representation of a Building
func (b Building) String() string {
	if int(b) >= len(buildingNames) {
		return "unknown"
	}
	return buildingNames[b]
}

func (b Building) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

const (
	CityHall Building = iota
	GovernmentPalace
	ColonialOffice
	PublicHouse
	NaturalHistoryMuseum
	University
	Factory
	Cathedral
	ForeignMinistry
	IndustrialWarehouse
	Landfill
	IndustrialPort
	StockExchange
	FortifiedCityWall
	IntelligenceBureau
	MilitaryBarracks
	NavalDockyard
	PrivateerBase
	ForestryOffice
	Sawmill
	Winery
	WineDistillery
	StoneQuarry
	ArchitecturalBureau
	GlassFactory
	OpticalLaboratory
	ChemicalPlant
	ExplosivesTestingSite
)

// FromString converts a string to a Building type
func FromString(s string) (Building, error) {
	for i, name := range buildingNames {
		if strings.EqualFold(name, s) {
			return Building(i), nil
		}
	}
	return 0, fmt.Errorf("invalid building name: %s", s)
}

var buildingNames = []string{
	"city_hall",
	"government_palace",
	"colonial_office",
	"public_house",
	"natural_history_museum",
	"university",
	"factory",
	"cathedral",
	"foreign_ministry",
	"industrial_warehouse",
	"landfill",
	"industrial_port",
	"stock_exchange",
	"fortified_city_wall",
	"intelligence_bureau",
	"military_barracks",
	"naval_dockyard",
	"privateer_base",
	"forestry_office",
	"sawmill",
	"winery",
	"wine_distillery",
	"stone_quarry",
	"architectural_bureau",
	"glass_factory",
	"optical_laboratory",
	"chemical_plant",
	"explosives_testing_site",
}
