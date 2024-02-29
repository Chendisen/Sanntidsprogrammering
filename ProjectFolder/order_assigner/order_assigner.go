package order_assigner


type ElevatorState struct {
    Behaviour   string
    Floor       int 
    Direction   string
    CabRequests []bool
}

type HRAInput struct {
    HallRequests    [][2]bool
    States          map[string]ElevatorState
}


func assign_orders(input HRAInput) 