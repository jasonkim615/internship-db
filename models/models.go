package models

//Company is the object with name attribute, and Application is a link to apply
type Company struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Application string `json:"application,omitempty" bson:"application,omitempty"`
}
