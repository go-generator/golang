package model

type Status string

const (
	StatusNew       Status = "N"
	StatusSubmitted Status = "S"
	StatusApproved  Status = "A"
)

//func (s Status) String() string {
//	return []string{"N", "S", "A"}
//}