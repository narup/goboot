package data

import (
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Address address data representation
type Address struct {
	BaseData     `bson:",inline"`
	Street1      string     `json:"street1" bson:"street1,omitempty"`
	Street2      string     `json:"street2" bson:"street2,omitempty"`
	City         string     `json:"city" bson:"city,omitempty" pson:"city"`
	State        string     `json:"state" bson:"state,omitempty" pson:"state"`
	ZipCode      string     `json:"zipCode" bson:"zipCode,omitempty" pson:"zip_code"`
	ZipCodeAddon string     `json:"zipCodeAddon" bson:"zipCodeAddon,omitempty"`
	Country      string     `json:"country" bson:"country,omitempty"`
	PatientID    string     `json:"patientId" bson:"patientId,omitempty"`
	ManagerRef   *mgo.DBRef `json:"-" bson:"manager,omitempty" pson:"manager"`
	Manager      *User      `json:"manager,omitempty" bson:"-"`
}

// CollectionName function from gmgo.Document interface
func (addr Address) CollectionName() string {
	return "address"
}

//Zipcode zip code lookup table
type Zipcode struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Zip      string        `json:"zip" bson:"zip"`
	City     string        `json:"city" bson:"primary_city,omitempty"`
	State    string        `json:"state" bson:"state,omitempty"`
	County   string        `json:"county" bson:"county,omitempty"`
	Timezone string        `json:"timezone" bson:"timezone,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (z Zipcode) CollectionName() string {
	return "zipcodes"
}

// Equals check if two addresses equal each other
func (addr Address) Equals(otherAddr *Address) bool {
	if strings.ToLower(addr.Street1) != strings.ToLower(otherAddr.Street1) {
		return false
	}

	if strings.ToLower(addr.Street2) != strings.ToLower(otherAddr.Street2) {
		return false
	}

	if strings.ToLower(addr.City) != strings.ToLower(otherAddr.City) {
		return false
	}

	if strings.ToLower(addr.State) != strings.ToLower(otherAddr.State) {
		return false
	}

	if strings.ToLower(addr.Country) != strings.ToLower(otherAddr.Country) {
		return false
	}

	return true
}
