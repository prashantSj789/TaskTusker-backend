package main

import "time"

type lane uint32

const(
	ToDo lane=iota
	InProgress
	DevComplete
	Testing
	Aprooved
)

func (l lane) String() string{
	switch l{
	case ToDo:
		return "To Do"
	case InProgress:
		return "In Progress"
	case DevComplete:
		return "Dev Complete"
	case Testing:
		return "Testing"
	case Aprooved:
		return "Aprooved"
	default:
		return "Lane does not exis"						
	}
}
func GetLane(s string) lane{
	switch s{
	case "To Do":
		return ToDo
	case "In Progress":
		return InProgress
	case "Dev Complete":
		return DevComplete
	case "Testing":
		return Testing
	case "Aprooved":
		return Aprooved
    default:
		return 6
	}
}

type Card struct{
    ID int `json:"id"`
	Name string `json:"name"`
	Task string `json:"task"`
	Organization string `json:"organizations"`
	Issued_to string `json:"issued_to"`
	Issued_at time.Time `json:"issued_at"`
	Lane string `json:"lane"`
}
type CreateCardRequest struct{
	Name string `json:"name"`
	Task string `json:"task"`
	Organization string `json:"organization"`
	Issued_to string `json:"handler"`
}
type GetAllCardResponse struct{
	Cards []*Card `json:"cards"`
}
func NewCard(name,task,organization,handler string) (*Card){
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return &Card{
		Name: name,
		Task: task,
		Organization: organization,
		Issued_to: handler,
		Issued_at: time.Now().In(loc),
		Lane: lane(0).String(),
	}

}