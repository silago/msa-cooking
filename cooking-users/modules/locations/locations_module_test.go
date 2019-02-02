package main

import (
	"github.com/pkg/errors"
	"github.com/silago/msa-cooking/cooking-users/modules/locations"
	"github.com/silago/msa-cooking/cooking-users/user_data_providers"
	"log"
	"strconv"
	"strings"
	"testing"
)
const TEST_USER_ID string="10"

type TestUserDataProvider struct {
	storage map[string]map[string]string
}

func (p *TestUserDataProvider) GetMany(userId string, names []string) map[string]string{
	result:=make(map[string]string)
	for k, v:=range p.storage[userId] {
		result[k] = "x"
		result[k] = v
	}
	return result
}

func (p *TestUserDataProvider) DecrementOne(userId string, resourceName string, value interface{}) (string, error) {
	currentValue,_:=strconv.Atoi(p.storage[userId][resourceName])
	var val int
	switch value.(type) {
	case int:
		val = value.(int)
		break
	case string:
		val, _ = strconv.Atoi(value.(string))
	}
	if currentValue<val {
		return string(currentValue), errors.New("not enough resources")
	} else {
		newVal:=strconv.Itoa(int(currentValue-val))
		e:=p.SetOne(userId,resourceName,newVal)
		return newVal, e
	}
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

func (p *TestUserDataProvider) IncrementMany(userId string, values map[string]interface{}) error {
	val:=	 0
	current:=0
	for k, v:= range values {
		current,_= strconv.Atoi(p.storage[userId][k])
		switch v.(type) {
		case int:
			val = v.(int)
		case string:
			val, _ = strconv.Atoi(v.(string))
		}
		p.storage[userId][k]=strconv.Itoa(val+current)
	}
	return nil



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

func (p *TestUserDataProvider) SetOne(userId string, key string, value interface{}) error {
	p.storage[userId][key] =value.(string)
	return nil
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
		"baker_level":"1",
	}
	provider.storage = map[string]map[string]string{
							TEST_USER_ID:userData,
	}
	return &provider
}

func TestUnlock(t *testing.T) {
	filePath := "assets/resources/locations.xml"
	config, err:= locations.LoadConfig(filePath)
	if err!=nil {
		t.Errorf("%s ", err.Error())
	}

	testProvider:=NewTestDataProvider()
	locationsModule:= locations.LocationsModule{config, testProvider}
	result:=locationsModule.Unlock(TEST_USER_ID,locations.LocationUnlockRequest{
		LocationName: "stars", ResourceName: "coins",
	})
	if !strings.Contains(result, "burgers_unlock") {
		t.Errorf("response {%s} does not contains {%s}" ,result, "burgers_unlock")
	} else if !strings.Contains(result, "coins") {
		t.Errorf("response {%s} does not contains {%s}" ,result, "coins")
	}
	if result:= locationsModule.Unlock(TEST_USER_ID,locations.LocationUnlockRequest{
		LocationName: "stars", ResourceName: "coins",
	}); !strings.Contains(result,"error") {
		t.Errorf("response {%s} does not contains error message" ,result)
	}
}

func TestUpgrade(t *testing.T) {
	filePath := "assets/resources/locations.xml"
	config, err:= locations.LoadConfig(filePath)
	if err!=nil {
		t.Errorf("%s ", err.Error())
	}

	testProvider:=NewTestDataProvider()
	locationsModule:= locations.LocationsModule{Config: config, Provider: testProvider}

	log.Printf("[test] not existing item")
	if result:= locationsModule.Upgrade(TEST_USER_ID, locations.ItemUpgradeRequest{
		LocationName: "stars", ResourceName: "coins",
	}); !strings.Contains(result,"not found") {
		t.Errorf("response {%s} does not contains error message" ,result)
	}

	if result:= locationsModule.Upgrade(TEST_USER_ID, locations.ItemUpgradeRequest{
		LocationName: "bakery", ResourceName: "coins", UpgradeName:"baker",
	}); !strings.Contains(result,"coins") {
		t.Errorf("response {%s} does not contains required params" ,result)
	} else if !strings.Contains(result, "baker_level") {
		t.Errorf("response {%s} does not contains required params" ,result)
	} else {
		log.Printf("[ok] upgrade bakery response: %s ", result)
	}

	if result:= locationsModule.Upgrade(TEST_USER_ID, locations.ItemUpgradeRequest{
		LocationName: "bakery", ResourceName: "coins", UpgradeName:"baker",
	}); !strings.Contains(result,"coins") {
		t.Errorf("response {%s} does not contains required params 1." ,result)
	} else if !strings.Contains(result, "baker_level") {
		t.Errorf("response {%s} does not contains required params 2." ,result)
	} else {
		log.Printf("[ok] upgrade bakery response: %s ", result)
	}




}

