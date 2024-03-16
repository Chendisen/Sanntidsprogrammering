package world_view

import (
	"encoding/json"
	"fmt"
	"time"
)

type StandardMessage struct {
	IPAddress string                `json:"IPAddress"`
	WorldView WorldView 			`json:"worldView"`
	SendTime  string                `json:"sendTime"`
}

const PeriodInMilliseconds int = 100


func (standardMessage StandardMessage) GetSenderIP() string {
	return standardMessage.IPAddress
}

func (standardMessage StandardMessage) GetWorldView() WorldView {
	return standardMessage.WorldView
}

func (standardMessage StandardMessage) GetSendTime() string {
	return standardMessage.SendTime
}

func CreateStandardMessage(worldView WorldView, myIP string, sendTime string) StandardMessage {
	return StandardMessage{
		IPAddress: myIP,
		WorldView: worldView,
		SendTime:  sendTime,
	}
}

func (standardMessage *StandardMessage) ContinuouslyUpdateTransmittedMessage(myView *WorldView, msgTx chan<- StandardMessage) {
	for {
		standardMessage.WorldView = *myView
		standardMessage.SendTime = time.Now().String()[11:19]
		msgTx <- *standardMessage
		time.Sleep(time.Duration(PeriodInMilliseconds) * time.Millisecond)
	}
}

func PackMessage(standardMessage StandardMessage) []byte {
	jsonBytes, err := json.Marshal(standardMessage)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic(err)
	}
	return jsonBytes
}

func UnpackMessage(jsonBytes []byte) StandardMessage {
	var standardMessage StandardMessage
	err := json.Unmarshal(jsonBytes, &standardMessage)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		panic(err)
	}
	return standardMessage
}
