package types

// member hold information of a member
type Member struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name"`
}
