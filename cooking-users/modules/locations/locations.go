package locations

import (
	"cooking-users/user_data_providers"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)


type ResponseMessage string
const  (
	EmptyConditionText    ResponseMessage = "Condition is empty"
	NotEnoughResourceText ResponseMessage = "Not Enough Resource"
	DecodeError ResponseMessage = "Decode error"
	EncodeError ResponseMessage = "Encode error"
	DatabaseError ResponseMessage = "Database error"
	Sucess ResponseMessage = "Success"
	LocationNotDound ResponseMessage = "Location not found"
)

func (s ResponseMessage) MakeJsonResponse(params interface{}) string {
	response, _:= json.Marshal(params)
	return string(response)
}

func (s ResponseMessage) ToString(params... interface{} ) string {
	return fmt.Sprintf(string(s), params...)
}

func (s ResponseMessage) AsJsonError(params... interface{} ) string {
	response, _:= json.Marshal(map[string]string{
		"error": fmt.Sprintf(string(s)),
	})
	return string(response)
}



func (location *Locations) GetCurrencyByName(name string)  int {
	switch name {
		case "coins":
			val, _ := strconv.Atoi(location.Coins)
			return val
		default:
			return -1
	}
}


type Config struct {
	XMLName xml.Name `xml:"config"`
	LocationsConfig   LocationsConfig   `xml:"locations"`
}


type LocationsConfig struct {
	XMLName xml.Name `xml:"locations"`
	Locations   []Locations   `xml:"location"`
}

type Locations struct {
	XMLName     xml.Name   `xml:"location"`
	Name        string     `xml:"name,attr"`
	Level       string     `xml:"level,attr"`
	Coins       string     `xml:"coins,attr"`
	Key         string     `xml:"key,attr"`
	Requirement Requirement `xml:"req"`
}

type Requirement struct {
	Gt     string     `xml:"gt,attr"`
	Value  string     `xml:"val,attr"`
}

func (r *Requirement) toInt() int {
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

type LocationUpgradeRequest struct {
	LocationName string `json:"location"`
	ResourceName string `json:"resource"`
}

type LocationUpgradeResponse struct {
	success bool
	msg     string
}

type LocationsModule struct {
	 Config   Config
	 Provider user_data_providers.UserDataProvider
}

func (module *LocationsModule) Unlock (userId string, requestData LocationUpgradeRequest)  string {
	config:=module.Config
	provider:=module.Provider
	for _, location := range config.LocationsConfig.Locations {
		if location.Name == requestData.LocationName {
			if location.Requirement.Gt=="" {
				return EmptyConditionText.ToString();//"Condition Gt is empty"
			}
			if val, e := strconv.Atoi(provider.GetOne(userId,location.Requirement.Gt)); e != nil {
				return DecodeError.ToString();
			} else if val < location.Requirement.toInt() {
				return NotEnoughResourceText.ToString()
			}
			upgradeCurrencyValue:=location.GetCurrencyByName(requestData.ResourceName);
			currentCurrencyValue, _:=  strconv.Atoi(provider.GetOne(userId,requestData.ResourceName))
			if upgradeCurrencyValue<0 || currentCurrencyValue < upgradeCurrencyValue {
				return NotEnoughResourceText.AsJsonError()
			}
			if err:=provider.SetMany(userId, map[string]interface{}{
				requestData.ResourceName: currentCurrencyValue - upgradeCurrencyValue,
				location.Key:             true,
			}); err!=nil {
				return DatabaseError.AsJsonError();
			} else {
				return Sucess.MakeJsonResponse(
					map[string]interface{}{
						requestData.ResourceName: currentCurrencyValue - upgradeCurrencyValue,
						location.Key:             true,
					}) //{
			}
		}
	}
	return LocationNotDound.AsJsonError()
}

func (module *LocationsModule) UnlockHandler (userId string, request *http.Request)  string {
		var requestData LocationUpgradeRequest
		if body, err:= ioutil.ReadAll(request.Body); err==nil {
			_ = json.Unmarshal(body, &requestData)
		} else {
			return err.Error()
		}
		return module.Unlock(userId, requestData);
}

func NewLocationUpdateHandler(configPath string, provider user_data_providers.UserDataProvider) (func (userId string, request *http.Request) string , error) {
	if config, err :=LoadConfig(configPath); err!=nil {
		panic(fmt.Sprintf("could not load Config {%s}", err.Error()))
	} else {
		module:=LocationsModule{config, provider}
		return module.UnlockHandler, nil;
	}
}






