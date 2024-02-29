package messages

import(
	"Sanntid/world_view"
)

func SendWorldView(wv world_view.WorldView, wvTx chan<- world_view.WorldView){
	wvTx <- wv
}