package locations

import (
	"encoding/json"
	"fmt"
	"github.com/silago/msa-cooking/cooking-users/user_data_providers"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type LocationsModule struct {
	 Config   Config
	 Provider user_data_providers.UserDataProvider
}

func (module *LocationsModule) Upgrade (userId string, requestData ItemUpgradeRequest)  string {
	if location := module.Config.LocationsConfig.GetLocationByName(requestData.LocationName); location ==nil {
		return LocationNotFound.AsJsonError()
	} else if upgrade:=location.GetUpgradeByName(requestData.UpgradeName);  upgrade==nil {
		return UpgradeNotFound.AsJsonError()
	} else {
		requirementsCache := make(map[string]int)
		var lastAvailableState *UpgradeState = nil
		for _, state := range upgrade.States {
			if _, ok := requirementsCache[state.Requirement.GetName()]; !ok {
				requirementsCache[state.Requirement.GetName()], _ = strconv.Atoi(module.Provider.GetOne(userId, state.Requirement.GetName()))
			}
			val := requirementsCache[state.Requirement.GetName()]
			if state.Requirement.Compare(val) == true {
				lastAvailableState = &state
			} else {
				break
			}
		}


		if lastAvailableState == nil {
			return ResponseMessage("no upgrades available").AsJsonError()
		}

		upgradeCurrencyValue := lastAvailableState.GetUpgradePriceByName(requestData.ResourceName)
		if upgradeCurrencyValue == "" {
			return ResponseMessage("resource not found").AsJsonError()
		}

		if newVal, err := module.Provider.DecrementOne(userId, requestData.ResourceName, upgradeCurrencyValue); err != nil {
			return ResponseMessage(err.Error() + fmt.Sprintf(" name: %s , current: %s, required: %s ", requestData.ResourceName, newVal, upgradeCurrencyValue)).AsJsonError()
		} else if err := module.Provider.IncrementMany(userId, map[string]interface{}{
			requestData.LocationName + "_upgrades": 1,
			requestData.UpgradeName + "_level":     1,
		}); err == nil {
			result:=make(map[string]interface{})
			data:= module.Provider.GetMany(userId,	[]string {requestData.LocationName + "_upgrades",	requestData.UpgradeName + "_level" } )
			for k, v:=range data {
				result[k]=v
			}

			result[requestData.ResourceName]= newVal;
			result[location.Key]= true
			return Sucess.MakeJsonResponse(result)
		} else {
			return ResponseMessage(err.Error()).AsJsonError()
		}
	}
}

func (module *LocationsModule) isRequirementSatisfied(userId string, requirement *Requirement) bool {
	requirementValue :=requirement.ToInt()
	if  requirement.Gt!="" {
		val, err:= strconv.Atoi(module.Provider.GetOne(userId,requirement.Gt))
		log.Printf("current %s", val )
		return err==nil && val > requirementValue
	} else if requirement.Eq!="" {
		val, err:= strconv.Atoi(module.Provider.GetOne(userId,requirement.Eq))
		return err==nil && val == requirementValue
	} else if requirement.Ge!="" {
		val, err:= strconv.Atoi(module.Provider.GetOne(userId,requirement.Ge))
		return err==nil && val >= requirementValue
	} else {
		return false
	}
}

func (module *LocationsModule) Unlock (userId string, requestData LocationUnlockRequest)  string {
	if location := module.Config.LocationsConfig.GetLocationByName(requestData.LocationName); location !=nil {
		if !module.isRequirementSatisfied(userId, &location.Requirement)  {
			log.Printf("requirement %s",location.Requirement)
			return NotEnoughResourceText.ToString()
		}

		upgradeCurrencyValue:=location.GetCurrencyByName(requestData.ResourceName);
		currentCurrencyValue, _:=  strconv.Atoi(module.Provider.GetOne(userId,requestData.ResourceName))
		if current, err:= module.Provider.DecrementOne(userId, requestData.ResourceName, upgradeCurrencyValue); err!=nil {
			log.Printf("current %s: %s, reqyired: %s",  requestData.ResourceName, current, upgradeCurrencyValue)
			return ResponseMessage(err.Error()).AsJsonError()
		} else {
			return Sucess.MakeJsonResponse(
				map[string]interface{}{
					requestData.ResourceName: currentCurrencyValue - upgradeCurrencyValue,
					location.Key:             true,
				})
		}
	}
	return LocationNotFound.AsJsonError()
}

func (module *LocationsModule) UpgradeHandler (userId string, request *http.Request)  string {
	var requestData ItemUpgradeRequest
	if body, err:= ioutil.ReadAll(request.Body); err==nil {
		_ = json.Unmarshal(body, &requestData)
	} else {
		return err.Error()
	}
	return module.Upgrade(userId, requestData)
}


func (module *LocationsModule) UnlockHandler (userId string, request *http.Request)  string {
		var requestData LocationUnlockRequest
		if body, err:= ioutil.ReadAll(request.Body); err==nil {
			_ = json.Unmarshal(body, &requestData)
		} else {
			return err.Error()
		}
		return module.Unlock(userId, requestData);
}



func NewLocationsModule(configPath string, provider user_data_providers.UserDataProvider) (*LocationsModule, error) {
	if config, err :=LoadConfig(configPath); err!=nil {
		panic(fmt.Sprintf("could not load Config {%s}", err.Error()))
	} else {
		module:=LocationsModule{config, provider}
		return &module, err;
	}
}



