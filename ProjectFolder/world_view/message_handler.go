package world_view

import (
	"encoding/json"
	"fmt"
)

type StandardMessage struct {
	IPAddress string                `json:"IPAddress"`
	WorldView WorldView 			`json:"worldView"`
	SendTime  string                `json:"sendTime"`
}


func (message StandardMessage) GetSenderIP() string {
	return message.IPAddress
}

func (message StandardMessage) GetWorldView() WorldView {
	return message.WorldView
}

func (message StandardMessage) GetSendTime() string {
	return message.SendTime
}

func CreateStandardMessage(worldView WorldView, myIP string, sendTime string) StandardMessage {
	return StandardMessage{
		IPAddress: myIP,
		WorldView: worldView,
		SendTime:  sendTime,
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
