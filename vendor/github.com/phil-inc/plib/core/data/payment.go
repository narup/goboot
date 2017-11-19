package data

import mgo "gopkg.in/mgo.v2"

// Payment data representation. Note that it has 2 bson fields with prefix stripe
// but it could be BrainTree token and customer Id. It's determined by provider value
type Payment struct {
	BaseData           `bson:",inline"`
	CardName           string     `json:"cardName" bson:"cardName" pson:"card_name" binding:"required"`
	Last4Digit         string     `json:"last4Digit" bson:"last4Digit" binding:"required"`
	Cvc                string     `json:"cvc" bson:"cvc" binding:"required"`
	ExpiryMonth        int        `json:"expiryMonth" bson:"expiryMonth" pson:"expiry_month" binding:"required"`
	ExpiryYear         int        `json:"expiryYear" bson:"expiryYear" pson:"expiry_year" binding:"required"`
	PaymentToken       string     `json:"paymentToken" bson:"stripeToken,omitempty"`
	ProviderCustomerID string     `json:"providerCustomerID" bson:"stripeCustomerId,omitempty"`
	Provider           string     `json:"provider" bson:"provider,omitempty"`
	FsaHsaCard         bool       `json:"fsaHsaCard" bson:"fsaHsaCard,omitempty" pson:"fsa_hsa_card"`
	TestCard           bool       `json:"testCard" bson:"testCard,omitempty"`
	PatientID          string     `json:"patientId" bson:"patientId,omitempty"`
	ManagerRef         *mgo.DBRef `json:"-" bson:"manager,omitempty"`
	Manager            *User      `json:"manager,omitempty" bson:"-"`
	CardNumber         string     `json:"cardNumber,omitempty" bson:"-"`
}

// CollectionName function from gmgo.Document interface
func (p Payment) CollectionName() string {
	return "paymentInfo"
}

//PaymentData represents the encrypted payment data.
type PaymentData struct {
	BaseData  `bson:",inline"`
	RawData   string `bson:"rawData,omitempty"`
	PaymentID string `bson:"paymentId,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (pd PaymentData) CollectionName() string {
	return "paymentData"
}

//PaymentAuthorizationToken stores the payment auth token for copay approvals
type PaymentAuthorizationToken struct {
	BaseData `bson:",inline"`
	UserID   string     `json:"userId" bson:"userId"`
	OrderID  string     `json:"orderId" bson:"orderId"`
	Token    string     `json:"token" bson:"token"`
	OrderRef *mgo.DBRef `json:"-" bson:"order,omitempty"` //NOTE: Not used now
}

// CollectionName function from gmgo.Document interface
func (pat PaymentAuthorizationToken) CollectionName() string {
	return "paymentAuthorizationToken"
}
