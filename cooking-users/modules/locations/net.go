package locations

import (
	"encoding/json"
	"fmt"
)

type ResponseMessage string
const  (
	EmptyConditionText    ResponseMessage = "Condition is empty"
	NotEnoughResourceText ResponseMessage = "Not Enough Resource"
	DecodeError           ResponseMessage = "Decode error"
	EncodeError           ResponseMessage = "Encode error"
	DatabaseError         ResponseMessage = "Database error"
	Sucess                ResponseMessage = "Success"
	LocationNotFound      ResponseMessage = "Location not found"
	UpgradeNotFound       ResponseMessage = "Upgrade not found"
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

type LocationUnlockRequest struct {
	LocationName string `json:"location"`
	UpgradeName string  `json:"item"`
	ResourceName string `json:"currency"`
}

type LocationUnlockResponse struct {
	success bool
	msg     string
}

type ItemUpgradeRequest struct {
	LocationName string `json:"location"`
	UpgradeName string  `json:"item"`
	ResourceName string `json:"currency"`
}

type ItemUpgradeResponse struct {

}
