package data

import mgo "gopkg.in/mgo.v2"

//Doctor data representation
type Doctor struct {
	BaseData    `bson:",inline"`
	Name        string `json:"name" bson:"name,omitempty"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber,omitempty"`
	FaxNumber   string `json:"faxNumber" bson:"faxNumber,omitempty"`
	Street1     string `json:"street1" bson:"street1,omitempty"`
	Street2     string `json:"street2" bson:"street2,omitempty"`
	City        string `json:"city" bson:"city,omitempty"`
	State       string `json:"state" bson:"state,omitempty"`
	ZipCode     string `json:"zipCode" bson:"zipCode,omitempty"`
	Country     string `json:"country" bson:"country,omitempty"`
	OrgName     string `json:"orgName" bson:"orgName,omitempty"`
	NPI         string `json:"npi" bson:"npi,omitempty"` //Indexed
	DEA         string `json:"dea" bson:"dea,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (d Doctor) CollectionName() string {
	return "doctor"
}

//MDOffice represents MD office data
type MDOffice struct {
	BaseData    `bson:",inline"`
	Name        string `json:"name" bson:"name,omitempty"`
	Street1     string `json:"street1" bson:"street1,omitempty"`
	Street2     string `json:"street2" bson:"street2,omitempty"`
	City        string `json:"city" bson:"city,omitempty"`
	State       string `json:"state" bson:"state,omitempty"`
	ZipCode     string `json:"zipCode" bson:"zipCode,omitempty"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber,omitempty"`
	FaxNumber   string `json:"faxNumber" bson:"faxNumber,omitempty"`
	WebsiteURL  string `json:"websiteUrl" bson:"websiteUrl,omitempty"`
	OfficeType  string `json:"officeType" bson:"officeType,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (mgo MDOffice) CollectionName() string {
	return "mdOffice"
}

//MDPartner data representation for MD partner
type MDPartner struct {
	BaseData          `bson:",inline"`
	MDOfficeID        string     `json:"mdOfficeId" bson:"mdOfficeId,omitempty"`
	ContractStatus    string     `json:"contractStatus" bson:"contractStatus,omitempty"`
	SalesContactName  string     `json:"salesContactName" bson:"salesContactName,omitempty"`
	SalesContactEmail string     `json:"salesContactEmail" bson:"salesContactEmail,omitempty"`
	PracticeName      string     `json:"practiceName" bson:"practiceName,omitempty"`
	Verified          bool       `json:"verified" bson:"verified"`
	DoctorRef         *mgo.DBRef `json:"-" bson:"doctor"`
	Doctor            *Doctor    `json:"doctor" bson:"-"`
}

//CollectionName returns collection name for the MDPartner data
func (mdp MDPartner) CollectionName() string {
	return "mDPartner"
}

// MDOfficeUser represents users from MD office who will be using MD dashboard to manage patients prescribed by a doctor in that office.
type MDOfficeUser struct {
	BaseData      `bson:",inline"`
	UserID        string         `json:"userId" bson:"userId,omitempty"`
	MDOfficeID    string         `json:"mdOfficeId" bson:"mdOfficeId,omitempty"`
	DoctorIDs     []string       `json:"doctorIds" bson:"doctorIds,omitempty"`
	MDOfficePrefs *MDOfficePrefs `json:"mdOfficePrefs" bson:"mdOfficePrefs,omitempty"`
	UserType      string         `json:"userType" bson:"userType,omitempty"`
}

//CollectionName returns collection name used in MongoDB. gmgo.Document interface implementation
func (mdOUsr MDOfficeUser) CollectionName() string {
	return "mdOfficeUser"
}

//MDOfficePrefs stores all the MD office user preferences. CommentPref/PriorOrRefillAuthPref has 4 possible values
//NotificationPrefSMS, NotificationPrefEmail, NotificationPrefEmailAndSMS, and NotificationPrefVoice
type MDOfficePrefs struct {
	ReportMail            bool   `json:"reportMail" bson:"reportMail,omitempty"`
	SendOnBehalf          bool   `json:"sendOnBehalf" bson:"sendOnBehalf,omitempty"`
	TempPasswordChanged   bool   `json:"tempPasswordChanged" bson:"tempPasswordChanged,omitempty"`
	AcceptedTermsOfUse    bool   `json:"acceptedTermsOfUse" bson:"acceptedTermsOfUse,omitempty"`
	CommentPref           string `json:"commentPref" bson:"commentPref,omitempty"`
	PriorOrRefillAuthPref string `json:"priorOrRefillAuthPref" bson:"priorOrRefillAuthPref,omitempty"`
}

//IsSMSNotificationForComment should send sms notification for comment
func (p MDOfficePrefs) IsSMSNotificationForComment() bool {
	return p.CommentPref == NotificationPrefEmailAndSMS || p.CommentPref == NotificationPrefSMS
}

//IsEmailNotificationForComment should send email notification for comment
func (p MDOfficePrefs) IsEmailNotificationForComment() bool {
	return p.CommentPref == NotificationPrefEmailAndSMS || p.CommentPref == NotificationPrefEmail
}

//IsSMSNotificationForPriorOrRefillAuth should send sms notification for prior or refill auth
func (p MDOfficePrefs) IsSMSNotificationForPriorOrRefillAuth() bool {
	return p.PriorOrRefillAuthPref == NotificationPrefEmailAndSMS || p.PriorOrRefillAuthPref == NotificationPrefSMS
}

//IsEmailNotificationForPriorOrRefillAuth should send sms notification for prior or refill auth
func (p MDOfficePrefs) IsEmailNotificationForPriorOrRefillAuth() bool {
	return p.PriorOrRefillAuthPref == NotificationPrefEmailAndSMS || p.PriorOrRefillAuthPref == NotificationPrefEmail
}
