package types

type Hobby struct {
	ID         string `json:"_id,omitempty" bson:"objectId"`
	Content    string `json:"content"`
	ID_Created string `json:"id_created" bson:"objectId"`
}
