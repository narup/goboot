package data

//Mention represents the FulfillmentTaskInbox of the user
type Mention struct {
	BaseData          `bson:",inline"`
	RxID              string `json:"rxId" bson:"rxId"`
	OrderNumber       string `json:"orderNumber" bson:"orderNumber"`
	PatientID         string `json:"patientId" bson:"patientId,omitempty"`
	PatientName       string `json:"patientName" bson:"patientName,omitempty"`
	Message           string `json:"message" bson:"message,omitempty"`
	CommentID         string `json:"commentId" bson:"commentId,omitempty"`
	Comment           string `json:"comment" bson:"comment,omitempty"`
	AgentID           string `json:"agentId" bson:"agentId,omitempty"`
	AgentName         string `json:"agentName" bson:"agentName,omitempty"`
	AssignerAgentID   string `json:"assignerAgentId" bson:"assignerAgentId,omitempty"`
	AssignerAgentName string `json:"assignerAgentName" bson:"assignerAgentName,omitempty"`
	MDOfficeID        string `json:"mdOfficeId" bson:"mdOfficeId,omitempty"`
	Resolved          bool   `json:"resolved" bson:"resolved"`
	Opened            bool   `json:"opened" bson:"opened"`
}

// CollectionName function to implement Document interface for patient
func (f Mention) CollectionName() string {
	return "fulfillmentTaskInbox"
}
