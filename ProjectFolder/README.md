Elevator Project Group 87
=========================

This is the elevator project for group 87. The project contains a handful of modules to simulate multiple elevators working together over the network. The network topology implemented is a peer to peer structure, but with a single master node to delegate all orders. The master is chosen by the simple decision rule that the node among the alive nodes with the highest ID will become the master. 

##  Contents

    - Modules description
    - Usage


##  Modules

The project consists of the following modules 

Communication: This module is responsible for establishing connection to the network, identifying which nodes exist on the network, and sending messages between the nodes on the network. All messaging is done via UDP broadcasting. Messages are sent periodically, at a rate of 10 messages per second. The transmit period may be modified by setting the constant PeriodInMilliseconds to the desired transmission period.

Driver: This module is the inteface that allows for control of the physical elevator model. 

Elevator: This module defines the elevator class, and the routines that go along with it. An elevator object is used to set and read the state of the physical elevator model. Paired with Driver the two modules are responsible for the low-level control of the physical elevator. Contains the submodules:

    Finite State Machine (fsm): This module is the state machine for the elevator. 

    Elevator Lights: Lights are not to be turned on before an order is confirmed across the network, to guarantee that the order will be serviced. This module is responsible for setting the right lights

    Requests: This module is responsible for sorting incoming requests that are delegated to each elevator, such that the elevator handles its assigned orders in an efficient manner.

    Stop Button: Can't be explained, must only be experienced.

Order Assigner: Responsible for assigning hall requests to the optimal elevator among the nodes connected to the network. This is implemented with the hand-out cost function hall_request_assigner. 

Process Pair: Implements process pair for the main goroutine. 

Resources: Contains 

Timer: Implements different timers used for elevator functionality. Contains doorOpen, watchdog, process pair and network timer

World View: Responsible for keeping track of the states of all elevators on the network, all hall requests seen from all parts of the network and all assigned hall requests. Implements logic for correct updates of requests by using cyclic counters and a sort of agreement for changing values. Also contains:

    Message handler: Defines the structure of a StandardMessage to be broadcasted on the network. Implements routines for converting a StandardMessage to json.

    Network Overview: Struct for storing nodes alive, and keeping track of which process is the master. The master is chosen based on an id of each process.


##  Running the project

Make sure to run the project using "go run main.go -id X", using a different integer id for each elevator. This is to make sure the elevators have different id's and agree on who must become master.  





     

