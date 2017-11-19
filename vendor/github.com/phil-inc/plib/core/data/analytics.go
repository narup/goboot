package data

import "time"

//CallCustomerReason stores the call customer reason
type CallCustomerReason struct {
	BaseData    `bson:",inline"`
	RxID        string     `json:"rxId,omitempty" pson:"rx_id"`
	OrderID     string     `json:"orderId,omitempty"`
	OrderNumber string     `json:"orderNumber,omitempty" pson:"order_number"`
	Date        *time.Time `json:"date,omitempty" pson:"date"`
	Category    string     `json:"category,omitempty" pson:"category"`
	Subcategory string     `json:"subcategory,omitempty" pson:"subcategory"`
}

// CollectionName function from gmgo.Document interface
func (callCustomerReason CallCustomerReason) CollectionName() string {
	return "callCustomerReason"
}
