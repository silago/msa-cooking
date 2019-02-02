package locations

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	XMLName xml.Name `xml:"config"`
	LocationsConfig   LocationsConfig   `xml:"locations"`
}

type LocationsConfig struct {
	XMLName xml.Name       `xml:"locations"`
	Locations   []Location `xml:"location"`
}

func (conf *LocationsConfig) GetLocationByName(name string) *Location {
	for _, location := range conf.Locations {
		if location.Name == name {
			return  &location
		}
	}
	return nil
}

type Location struct {
	XMLName     xml.Name   `xml:"location"`
	Name        string     `xml:"name,attr"`
	Level       string     `xml:"level,attr"`
	Coins       string     `xml:"coins,attr"`
	Crystals       string     `xml:"crystals,attr"`
	Key         string     `xml:"key,attr"`
	Requirement Requirement `xml:"req"`
	Upgrades    []UpgradeItem `xml:"upgrade"`
}

type UpgradeItem struct {
	Name     	string     `xml:"name,attr"`
	Value		string     `xml:"val,attr"`
	States      []UpgradeState `xml:"state"`
}

type UpgradeState struct {
	Crystals    string     `xml:"coins,attr"`
	Coins		string     `xml:"crystals,attr"`
	Requirement Requirement `xml:"req"`
}

func (s *UpgradeState) GetUpgradePriceByName(name string) string {
	switch  name {
		case "crystals":
			return s.Crystals
		case "coins":
			return s.Coins
	}
	return ""
}

type Requirement struct {
	Gt     string     `xml:"gt,attr"`
	Eq     string     `xml:"eq,attr"`
	Ge     string     `xml:"ge,attr"`
	Value  string     `xml:"val,attr"`
}

type RequirementType string
const Gt RequirementType = "Gt"
const Eq RequirementType = "Eq"
const Ge RequirementType = "Ge"

func (r *Requirement) Compare(value int) bool {
	requirement:=r.ToInt()
	log.Printf("comparing %s:  %d(%s) %s %d(%s)",r.GetName(),value,"value", r.GetType(), requirement,"requirement")
	switch r.GetType() {
		case Gt:
			return value > requirement
		case Eq:
			return value == requirement
		case Ge:
			return value >= requirement
	}
	return  false
}

func (r *Requirement) GetType() RequirementType {
	if r.Gt!="" { return Gt }
	if r.Eq!="" { return Eq }
	if r.Ge!="" { return Ge }
	return ""
}

func (r *Requirement) GetName() string {
	if r.Gt!="" { return r.Gt }
	if r.Eq!="" { return r.Eq }
	if r.Ge!="" { return r.Ge }
	return ""
}

func (r *Requirement) ToInt() int {
	if val, err := strconv.Atoi(r.Value); err!=nil {
		log.Printf("cannot decode requirement {%s}", r)
		return -1
	} else {
		return val
	}
}

func LoadConfig(path string) ( Config, error ) {
	xmlFile, err := os.Open(path)
	var config Config

	if err == nil {
		defer xmlFile.Close()
		byteValue, _ := ioutil.ReadAll(xmlFile)
		err= xml.Unmarshal(byteValue, &config)
	}
	return config,  err
}

func (location *Location) GetCurrencyByName(name string)  int {
	switch name {
	case "crystals":
		val, _ := strconv.Atoi(location.Crystals)
		return val
	case "coins":
		val, _ := strconv.Atoi(location.Coins)
		return val
	default:
		return -1
	}
}




func (location *Location) GetUpgradeByName(name string) *UpgradeItem {
	for _, item:= range location.Upgrades {
		if item.Name== name {
			return &item
		}
	}
	return nil
}
