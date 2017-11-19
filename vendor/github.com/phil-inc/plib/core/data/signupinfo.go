package data

import "gopkg.in/mgo.v2"

//SignupInfo represents the sign up info of the user
type SignupInfo struct {
	BaseData      `bson:",inline"`
	PatientRef    *mgo.DBRef `json:"-" bson:"patient,omitempty"`
	Patient       *Patient   `json:"patient" bson:"-"`
	SentText      bool       `json:"sentText" bson:"sentText"`
	LeftVoicemail bool       `json:"leftVoicemail" bson:"leftVoicemail"`
	NotInterested bool       `json:"notInterested" bson:"notInterested"`
	AddedRx       bool       `json:"addedRx" bson:"addedRx"`
	HelpEmailSent bool       `json:"helpEmailSent" bson:"helpEmailSent"`
	FollowUpTime  int        `json:"followUpTime" bson:"followUpTime"`
	AgentName     string     `json:"agentName" bson:"agentName,omitempty"`
	Comment       string     `json:"comment" bson:"comment,omitempty"`
}

// CollectionName function to implement Document interface for patient
func (s SignupInfo) CollectionName() string {
	return "signupInfo"
}
