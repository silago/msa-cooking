package main

import (
	"cooking-users/modules/locations"
	"cooking-users/user_data_providers"
	"strconv"
	"strings"
	"testing"
)
const TEST_USER_ID string="10"

type TestUserDataProvider struct {
	storage map[string]map[string]string
}

func (TestUserDataProvider) DeleteAll(userId string) error {
	panic("implement me")
}

func (TestUserDataProvider) IncrementOne(userId string, resourceName string, value interface{}) (string, error) {
	panic("implement me")
}

func (p TestUserDataProvider) GetOne(userId string, name string) string {
	return p.storage[userId][name];
}

func (TestUserDataProvider) IncrementMany(userId string, values map[string]interface{}) error {
	panic("implement me")
}

func (p TestUserDataProvider) SetMany(userId string, values map[string]interface{}) error {
	for k, v:= range values {

		switch v.(type) {
		case int:
			p.storage[userId][k]=strconv.Itoa(v.(int));
		case string:
			p.storage[userId][k]=v.(string)
		}
	}
	return nil
}

func (p TestUserDataProvider) SetOne(userId string, key string, value interface{}) error {
	panic("implement me")
}

func (p TestUserDataProvider) GetAll(userId string /*, values map[string]float64 */) map[string]string {
	panic("implement me")
}

func (TestUserDataProvider) GetOrCreate(userId string, values map[string]interface{}) (user_data_providers.User, error) {
	panic("implement me")
}

func NewTestDataProvider() user_data_providers.UserDataProvider {
	provider:= TestUserDataProvider{};
	userData:=map[string]string{
		"coins":"3000",
	}
	provider.storage = map[string]map[string]string{
							TEST_USER_ID:userData,
	}
	return provider
}

func TestUnlock(t *testing.T) {
	filePath := "assets/resources/locations.xml"
	config, err:= locations.LoadConfig(filePath)
	if err!=nil {
		t.Errorf("%s ", err.Error())
	}

	testProvider:=NewTestDataProvider();
	locationsModule:= locations.LocationsModule{config, testProvider}
	result:=locationsModule.Unlock(TEST_USER_ID,locations.LocationUpgradeRequest{
		"stars","coins",
	})
	if !strings.Contains(result, "burgers_unlock") {
		t.Errorf("response {%s} does not contains {%s}" ,result, "burgers_unlock")
	} else if !strings.Contains(result, "coins") {
		t.Errorf("response {%s} does not contains {%s}" ,result, "coins")
	}
	// check not enought respurces
	if result = locationsModule.Unlock(TEST_USER_ID,locations.LocationUpgradeRequest{
		"stars","coins",
	}); !strings.Contains(result,locations.NotEnoughResourceText.AsJsonError()) {
		t.Errorf("response {%s} does not contains error message " ,result)
	}
}

