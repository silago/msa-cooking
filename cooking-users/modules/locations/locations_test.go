package locations

import (
	"encoding/xml"
	"testing"
)


func createLocationModule() *LocationsModule {
	testProvider:=NewTestDataProvider()
	locationsModule:=LocationsModule{
			Config: Config{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				LocationsConfig: LocationsConfig{
					XMLName: xml.Name{
						Space: "",
						Local: "",
					},
					Locations: []Location{
						{
							XMLName: xml.Name{
								Space: "",
								Local: "",
							},
							Name:     "loc",
							Level:    "",
							Coins:    "10",
							Crystals: "",
							Key:      "key",
							Requirement: Requirement{
								Gt:    "coins",
								Eq:    "",
								Ge:    "",
								Value: "10",
							},
							Upgrades: []UpgradeItem{
								{
									Name:   "up",
									Value:  "",
									States: []UpgradeState{
										{
											Crystals: "",
											Coins:    "10",
											Requirement: Requirement{
												Gt:    "",
												Eq:    "baker_level",
												Ge:    "",
												Value: "1",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Provider: testProvider,
		}
	return &locationsModule;
}

func TestLocationsModule_Upgrade(t *testing.T) {
	type args struct {
		userId      string
		requestData ItemUpgradeRequest
	}


	tests := []struct {
		name   string
		module *LocationsModule
		args   args
		want   string
	}{
		{
			name:   "",
			module: createLocationModule(),
			args: args{
				userId: "10",
				requestData: ItemUpgradeRequest{
					LocationName: "loc",
					UpgradeName:  "up",
					ResourceName: "coins",
				},
			},
			want: "{\"baker_level\":\"1\",\"coins\":\"2990\",\"key\":true,\"loc_upgrades\":\"1\",\"up_level\":\"1\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.module.Upgrade(tt.args.userId, tt.args.requestData); got != tt.want {
				t.Errorf("LocationsModule.Upgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationsModule_isRequirementSatisfied(t *testing.T) {
	type args struct {
		userId      string
		requirement *Requirement
	}
	tests := []struct {
		name   string
		module *LocationsModule
		args   args
		want   bool
	}{
		{
			name:   "",
			module: createLocationModule(),
			args:   args{
				userId: "10",
				requirement: &Requirement{
					Gt:    "",
					Eq:    "baker_level",
					Ge:    "",
					Value: "1",
				},
			},
			want:   true,
		},
		{
			name:   "",
			module: createLocationModule(),
			args:   args{
				userId: "10",
				requirement: &Requirement{
					Gt:    "",
					Eq:    "baker_level",
					Ge:    "",
					Value: "10",
				},
			},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.module.isRequirementSatisfied(tt.args.userId, tt.args.requirement); got != tt.want {
				t.Errorf("LocationsModule.isRequirementSatisfied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocationsModule_Unlock(t *testing.T) {
	type args struct {
		userId      string
		requestData LocationUnlockRequest
	}
	tests := []struct {
		name   string
		module *LocationsModule
		args   args
		want   string
	}{
		{
			name:   "",
			module: createLocationModule(),
			args:   args{
				userId: "10",
				requestData: LocationUnlockRequest{
					LocationName: "loc",
					ResourceName: "coins",
				},
			},
			want:   "{\"coins\":2990,\"key\":true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.module.Unlock(tt.args.userId, tt.args.requestData); got != tt.want {
				t.Errorf("LocationsModule.Unlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

