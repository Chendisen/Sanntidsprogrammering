package world_view

import (
	"Sanntid/communication/localip"
	"Sanntid/communication/peers"
	"fmt"
)

type NetworkOverview struct {
	MyIP       string
	NodesAlive []string
	Master     string
}

func MakeNetworkOverview() NetworkOverview {
	myIP, _ := localip.LocalIP()
	//myIP := fmt.Sprintf("%d", os.Getpid())
	nodesAlive := make([]string, 1)
	nodesAlive[0] = myIP
	return NetworkOverview{MyIP: myIP, NodesAlive: nodesAlive, Master: myIP}
}

func (networkOverview NetworkOverview) GetMyIP() string {
	return networkOverview.MyIP
}

func (networkOverview NetworkOverview) AmIMaster() bool {
	if networkOverview.Master == networkOverview.MyIP {
		return true
	} else {
		return false
	}
}

func (networkOverview NetworkOverview) AmIAlive(myIP string) bool {
	var amIAlive bool = false
	for _, aliveNode := range networkOverview.NodesAlive {
		amIAlive = amIAlive || myIP == aliveNode
	}
	return amIAlive
}

func (networkOverview *NetworkOverview) ShouldUpdateList(p peers.PeerUpdate) bool {
	if len(p.Lost) != 0 {
		return true
	} else if len(p.New) != 0 {
		return true
	} else {
		return false
	}
}

func (networkOverview *NetworkOverview) ShouldUpdateMaster(p peers.PeerUpdate) (bool, string) {
	var shouldUpdate bool = false
	var newMaster string = ""
	if len(p.Lost) != 0 {
		for _, lostNode := range p.Lost {
			if lostNode == networkOverview.Master {
				shouldUpdate = true
				for _, candidate := range p.Peers {
					if candidate > newMaster {
						newMaster = candidate
					}
				}
				return shouldUpdate, newMaster
			}
		}
	} else if p.New > networkOverview.Master {
		newMaster = p.New
		shouldUpdate = true
		return shouldUpdate, newMaster
	}
	return shouldUpdate, newMaster
}

func (networkOverview *NetworkOverview) UpdateMaster(newMaster string) {
	networkOverview.Master = newMaster
}

func (networkOverview NetworkOverview) NetworkLost(p peers.PeerUpdate) bool {
	var networkGoing bool = false
	for _, aliveNode := range p.Peers {
		networkGoing = networkGoing || networkOverview.MyIP == aliveNode
	}
	return !networkGoing
}

func (networkOverview *NetworkOverview) UpdateNetworkOverview(p peers.PeerUpdate) {
	networkOverview.NodesAlive = p.Peers
	shouldUpdateMaster, newMaster := networkOverview.ShouldUpdateMaster(p)

	if shouldUpdateMaster {
		networkOverview.UpdateMaster(newMaster)
	}
}

func (networkOverview NetworkOverview) Print() {
	fmt.Printf("Current alive nodes: \n")
	for _, IP := range networkOverview.NodesAlive {
		fmt.Printf("A node	%s\n", IP)
	}
	fmt.Println("")
}
