package message_handler

import (
	"encoding/json"
	"fmt"
	"Sanntid/world_view"
)

type StandardMessage struct {
	IPAddress string               `json:"IPAddress"`
	WorldView world_view.WorldView `json:"worldView"`
	SendTime  string               `json:"sendTime"`
}

// TODO: Should probably have functions that take in struct from wordlview
// 			and make the whole message, such that it can be called from main
// 			and then be given the correct variables.

func GetSenderIP(message StandardMessage) string {
	return message.IPAddress
}

func GetWorldView(message StandardMessage) world_view.WorldView {
	return message.WorldView
}

func GetSendTime(message StandardMessage) string {
	return message.SendTime
}

func CreateStandardMessage(a_world_view world_view.WorldView, ip_address string, send_time string) StandardMessage {
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
