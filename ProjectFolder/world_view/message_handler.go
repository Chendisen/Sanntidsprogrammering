package world_view

import (
	"encoding/json"
	"fmt"
	//"Sanntid/world_view"
)

type StandardMessage struct {
	IPAddress string               `json:"IPAddress"`
	WorldView WorldView `json:"worldView"`
	SendTime  string               `json:"sendTime"`
}


func GetSenderIP(message StandardMessage) string {
	return message.IPAddress
}

func GetWorldView(message StandardMessage) WorldView {
	return message.WorldView
}

func GetSendTime(message StandardMessage) string {
	return message.SendTime
}

func CreateStandardMessage(a_world_view WorldView, ip_address string, send_time string) StandardMessage {
	return StandardMessage{
		IPAddress: ip_address,
		WorldView: a_world_view,
		SendTime:  send_time,
	}
}

func PackMessage(message StandardMessage) []byte {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic(err)
	}
	return jsonBytes
}

func UnpackMessage(jsonBytes []byte) StandardMessage {
	var message StandardMessage
	err := json.Unmarshal(jsonBytes, &message)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		panic(err)
	}
	return message
}
