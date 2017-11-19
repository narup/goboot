package data

import "time"

//RorEvent - Ror event
type RorEvent struct {
	BaseData    `bson:",inline"`
	ManagerID   string     `json:"managerId" bson:"managerId"`
	PatientID   string     `json:"patientId" bson:"patientId"`
	RxID        string     `json:"rxId" bson:"rxId"`
	OrderID     string     `json:"orderId" bson:"orderId"`
	OrderNumber string     `json:"orderNumber" bson:"orderNumber"`
	EventType   string     `json:"eventType" bson:"eventType"`
	TimeStamp   *time.Time `json:"timeStamp" bson:"timeStamp"`
}

// CollectionName function from gmgo.Document interface
func (p RorEvent) CollectionName() string {
	return "rorEvent"
}
