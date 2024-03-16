Elevator Project Group 87
=========================

This is the elevator project for group 87. The project contains a handful of modules to simulate multiple elevators working together over the network.

##  Contents

    - Modules description
    - Usage


##  Modules

The project consists of the following modules 

Communication module: This module is responsible for establishing connection to the network, identifying which nodes exist on the network, and sending messages between the nodes on the network. All messaging is done via UDP broadcasting. Messages are sent periodically, at a rate of 10 messages per second. The transmit period may be modified by setting the constant PeriodInMilliseconds to the desired transmission period.

Driver: This module is the inteface that allows for control of the physical elevator model. 

Elevator: This module defines the elevator class, and the routines that go along with it. An elevator object is used to set and read the state of the physical elevator model. Paired with Driver the two modules are responsible for the low-level control of the physical elevator. 

Finite State Machine (fsm): This module is the state machine for the elevator. 

Message handler: 