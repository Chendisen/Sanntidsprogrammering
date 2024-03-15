package resources

type UpdateRequestType int

const (
	SetBehaviour UpdateRequestType = iota
	SetFloor
	SetDirection
	SeenRequestAtFloor
	FinishedRequestAtFloor
	SetMyAvailabilityStatus
	SetAssignedOrders
)

type UpdateRequest struct {
	Type  UpdateRequestType
	Value interface{}
}

func GenerateUpdateRequest(requestType UpdateRequestType, value interface{}) UpdateRequest {
	return UpdateRequest{
		Type:  requestType,
		Value: value,
	}
}
