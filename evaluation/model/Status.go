package model

type Status string

const (
	New       Status = "N"
	Submitted Status = "S"
	Approved  Status = "A"
)

//func (s Status) String() string {
//	return []string{"N", "S", "A"}
//}