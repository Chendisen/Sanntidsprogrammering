package world_view

type RequestType int 
const(
	SetBehaviour RequestType = iota
	SetFloor
	SetDirection
	SeenRequestAtFloor
	FinishedRequestAtFloor
	SetMyAvailabilityStatus
)

type UpdateRequest struct {
    Type  RequestType
    Value interface{}
}

func GenerateUpdateRequest(requestType RequestType, value interface{}) UpdateRequest {
    return UpdateRequest{
        Type:  requestType,
        Value: value,
    }
}