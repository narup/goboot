package data

import "gopkg.in/mgo.v2"

// Insurance data representation
type Insurance struct {
	BaseData        `bson:",inline"`
	ClientName      string     `json:"clientName" bson:"clientName,omitempty" pson:"client_name"`
	ProviderName    string     `json:"providerName" bson:"providerName,omitempty" pson:"provider_name"`
	GroupNumber     string     `json:"groupNumber" bson:"groupNumber,omitempty" pson:"group_number"`
	BinNumber       string     `json:"binNumber" bson:"binNumber,omitempty" pson:"bin_number"`
	PcnNumber       string     `json:"pcnNumber" bson:"pcnNumber,omitempty" pson:"pcn_number"`
	InsuranceID     string     `json:"insuranceId" bson:"insuranceId,omitempty" pson:"insurance_id"`
	CardImageID     string     `json:"cardImageId" bson:"cardImageId,omitempty" pson:"card_image_id"`
	PatientID       string     `json:"patientId" bson:"patientId,omitempty"`
	PhoneNumber     string     `json:"phoneNumber" bson:"phoneNumber,omitempty"`
	IsGovtSponsored string     `json:"isGovtSponsored" bson:"isFederalSponsoredInsurance,omitempty" pson:"is_federal_sponsored_insurance"`
	Primary         bool       `json:"primary" bson:"primary" pson:"primary"`
	ManagerRef      *mgo.DBRef `json:"-" bson:"manager,omitempty" pson:"manager"`
	Manager         *User      `json:"manager,omitempty" bson:"-"`
	Verified        bool       `json:"verified" bson:"verified,omitempty" pson:"verified"`
}

// CollectionName function from gmgo.Document interface
func (i Insurance) CollectionName() string {
	return "insurance"
}

//MissingInsuranceToken stores the missing insurance token
type MissingInsuranceToken struct {
	BaseData `bson:",inline"`
	UserID   string `json:"userId" bson:"userId"`
	Token    string `json:"token" bson:"token"`
}

// CollectionName function from gmgo.Document interface
func (missingInsurance MissingInsuranceToken) CollectionName() string {
	return "missingInsuranceToken"
}

//InvalidInsuranceToken stores the insurance exception token
type InvalidInsuranceToken struct {
	BaseData `bson:",inline"`
	UserID   string `json:"userId" bson:"userId"`
	Token    string `json:"token" bson:"token"`
}

// CollectionName function from gmgo.Document interface
func (insuranceException InvalidInsuranceToken) CollectionName() string {
	return "invalidInsuranceToken"
}
