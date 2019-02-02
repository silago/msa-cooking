package locations

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestLocationsConfig_GetLocationByName(t *testing.T) {
	type args struct {
		name string
	}
	locationConfig:=&LocationsConfig{}
	locationConfig.Locations = make ([]Location, 2)
	location1:=Location{Name:"location1"}
	location2:=Location{Name:"location2"}

	locationConfig.Locations[0] = location1
	locationConfig.Locations[1] = location2
	tests := []struct {
		name string
		conf *LocationsConfig
		args args
		want *Location
	}{
		{name:"get existing location by name ", conf: locationConfig, args:args{name:"location1"},  want:&location1},
		{name:"get unexsisting location by name ", conf: locationConfig, args:args{name:"location3"},  want:nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.conf.GetLocationByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocationsConfig.GetLocationByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpgradeState_GetUpgradePriceByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		s    *UpgradeState
		args args
		want string
	}{
		{name:"get price by name 'coins'", s:&UpgradeState{Coins:"100",Crystals:"200"},args:args{name:"coins"}, want:"100" },
		{name:"get price by name 'crystals'", s:&UpgradeState{Coins:"100",Crystals:"200"},args:args{name:"crystals"}, want:"200" },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetUpgradePriceByName(tt.args.name); got != tt.want {
				t.Errorf("UpgradeState.GetUpgradePriceByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequirement_Compare(t *testing.T) {
	type args struct {
		value int
	}
	tests := []struct {
		name string
		r    *Requirement
		args args
		want bool
	}{
		{name:"eq test #1", r:&Requirement{Eq:"some field",Value:"20"}, args:args{value:20},want:true},
		{name:"eq test #2", r:&Requirement{Eq:"some field",Value:"20"}, args:args{value:21},want:false},

		{name:"gt test #1", r:&Requirement{Gt:"some field",Value:"20"}, args:args{value:20},want:false},
		{name:"gt test #2", r:&Requirement{Gt:"some field",Value:"20"}, args:args{value:21},want:true},

		{name:"ge test #1", r:&Requirement{Ge:"some field",Value:"20"}, args:args{value:20},want:true},
		{name:"ge test #2", r:&Requirement{Ge:"some field",Value:"20"}, args:args{value:21},want:true},
		{name:"ge test #3", r:&Requirement{Ge:"some field",Value:"20"}, args:args{value:19},want:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Compare(tt.args.value); got != tt.want {
				t.Errorf("Requirement.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequirement_GetType(t *testing.T) {
	tests := []struct {
		name string
		r    *Requirement
		want RequirementType
	}{
		{
			name: "gt",
			r:    &Requirement{Gt:"some_field"},
			want: "Gt",
		},
		{
			name: "ge",
			r:    &Requirement{Ge:"some_field"},
			want: "Ge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Requirement.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequirement_GetName(t *testing.T) {
	tests := []struct {
		name string
		r    *Requirement
		want string
	}{
		{
			name: "other_field",
			r:    &Requirement{Eq:"other_field"},
			want: "other_field",
		},
		{
			name: "some_field",
			r:    &Requirement{Gt:"some_field"},
			want: "some_field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetName(); got != tt.want {
				t.Errorf("Requirement.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequirement_ToInt(t *testing.T) {
	tests := []struct {
		name string
		r    *Requirement
		want int
	}{
		{
			name: "to 10 int",
			r: &Requirement{
				Gt:    "some_field",
				Eq:    "",
				Ge:    "",
				Value: "10",
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.ToInt(); got != tt.want {
				t.Errorf("Requirement.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_GetCurrencyByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		location *Location
		args     args
		want     int
	}{
		{
			name: "",
			location: &Location{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				Name:     "",
				Level:    "",
				Coins:    "100",
				Crystals: "",
				Key:      "",
				Requirement: Requirement{
					Gt:    "",
					Eq:    "",
					Ge:    "",
					Value: "",
				},
				Upgrades: nil,
			},
			args: args{
				name: "coins",
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.location.GetCurrencyByName(tt.args.name); got != tt.want {
				t.Errorf("Location.GetCurrencyByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_GetUpgradeByName(t *testing.T) {
	type args struct {
		name string
	}


	upgradeItem:=UpgradeItem{
		Name:   "TestUpgradeItem",
		Value:  "",
		States: nil,
	}

	tests := []struct {
		name     string
		location *Location
		args     args
		want     *UpgradeItem
	}{
		{
			name: "",
			location: &Location{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				Name:     "",
				Level:    "",
				Coins:    "",
				Crystals: "",
				Key:      "",
				Requirement: Requirement{
					Gt:    "",
					Eq:    "",
					Ge:    "",
					Value: "",
				},
				Upgrades: []UpgradeItem{
					upgradeItem,
				},
			},
			args: args{
				name: "TestUpgradeItem",
			},
			want: &upgradeItem,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.location.GetUpgradeByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Location.GetUpgradeByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
