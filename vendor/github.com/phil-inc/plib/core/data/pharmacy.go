package data

import "gopkg.in/mgo.v2"

const (
	//MailOrder keeps track of when rx comes from mail order pharmacy
	MailOrder = "MailOrder:"
)

// Pharmacy data representation.
// Index used is a compound index
// (name = "name-phone-index", def = "{'name': 1, 'phoneNumber': 1}")
type Pharmacy struct {
	BaseData                 `bson:",inline"`
	Name                     string `json:"name" bson:"name" pson:"name" binding:"required"`
	PhoneNumber              string `json:"phoneNumber" bson:"phoneNumber,omitempty" pson:"phone_number"`
	FaxNumber                string `json:"faxNumber" bson:"faxNumber,omitempty" pson:"fax_number"`
	StoreNumber              string `json:"storeNumber" bson:"storeNumber,omitempty"`
	Street1                  string `json:"street1" bson:"street1,omitempty"`
	City                     string `json:"city" bson:"city,omitempty" pson:"city"`
	State                    string `json:"state" bson:"state,omitempty" pson:"state"`
	ZipCode                  string `json:"zipCode" bson:"zipCode,omitempty" pson:"zip_code"`
	Email                    string `json:"email" bson:"email,omitempty"`
	Country                  string `json:"country" bson:",omitempty"`
	ContactName              string `json:"contactName" bson:"contactName,omitempty" pson:"contact_name"`
	ShippingLabelDisplayName string `json:"shippingLabelDisplayName" bson:"shippingLabelDisplayName,omitempty" pson:"shipping_label_display_name"`
	TimeZone                 string `json:"timeZone" bson:"timeZone,omitempty" pson:"time_zone"`
	NPI                      string `json:"npi" bson:"npi,omitempty" pson:"npi"`
	DEA                      string `json:"dea" bson:"dea,omitempty" pson:"dea"`
	PartnershipStatus        string `json:"partnershipStatus" bson:"partnershipStatus,omitempty"`
	WeekDaysSupportHours     string `json:"weekDaysSupportHours" bson:"weekDaysSupportHours,omitempty" pson:"weekdays_support_hours"`
	WeekendSupportHours      string `json:"weekendSupportHours" bson:"weekendSupportHours,omitempty" pson:"weekend_support_hours"`
	Partner                  bool   `json:"partner" bson:"partner,omitempty" pson:"partner"`
	InactivePartner          bool   `json:"inactivePartner" bson:"inactivePartner,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (ph Pharmacy) CollectionName() string {
	return "pharmacy"
}

//IsInTrial returns true if partnership status is in trial phase
func (ph Pharmacy) IsInTrial() bool {
	return ph.PartnershipStatus == "Trial"
}

//PharmacyUser data fro pharmacy user
type PharmacyUser struct {
	BaseData    `bson:",inline"`
	UserRef     *mgo.DBRef `json:"-" bson:"rexUser,omitempty"`
	User        *User      `json:"rexUser" bson:"-"`
	PharmacyRef *mgo.DBRef `json:"-" bson:"pharmacy,omitempty"`
	Pharmacy    *Pharmacy  `json:"pharmacy" bson:"-"`
}

// CollectionName function from gmgo.Document interface
func (pu PharmacyUser) CollectionName() string {
	return "pharmacyUser"
}

// InitializeOriginPharmacy initialize origin pharmacy
func InitializeOriginPharmacy(name string, mailOrder bool, street string) *Pharmacy {
	pharmacy := new(Pharmacy)
	pharmacy.InitData()
	pharmacy.Name = name
	pharmacy.ContactName = name
	if mailOrder {
		pharmacy.Name = MailOrder + pharmacy.Name
	}
	pharmacy.Street1 = street
	return pharmacy
}

// InitializePartnerPharmacy initializes partner pharmacy
func InitializePartnerPharmacy(pharmacy *Pharmacy) *Pharmacy {
	pharmacy.InitData()
	pharmacy.ContactName = pharmacy.Name
	return pharmacy
}
