package data

import (
	"time"
)

// AdminConfig application configuration data
type AdminConfig struct {
	BaseData                    `bson:",inline"`
	LiveAppVersion              string                 `json:"liveAppVersion" bson:"liveAppVersion"`
	LiveAppURL                  string                 `json:"liveAppUrl" bson:"liveAppUrl,omitempty"`
	ServerMaintainenace         string                 `json:"serverMaintainenace" bson:"serverMaintainenace"`
	TwoDayMailCharge            string                 `json:"twoDayMailCharge" bson:"twoDayMailCharge,omitempty"`
	GlobalPromo                 string                 `json:"globalPromo" bson:"globalPromo,omitempty"`
	PaymentProvider             string                 `json:"paymentProvider" bson:"paymentProvider"`
	SupportedStates             []string               `json:"supportedStates" bson:"supportedStates"`
	StatesWithPartnerPharmacies []string               `json:"statesWithPartnerPharmacies" bson:"statesWithPartnerPharmacies"`
	CustomerSupportConfig       *CustomerSupportConfig `json:"customerSupportConfig" bson:"customerSupportConfig"`
	OutboundNumberMap           map[string]string      `json:"outboundNumberMap" bson:"outboundNumberMap,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (cfg AdminConfig) CollectionName() string {
	return "adminConfig"
}

// AdminConfig application configuration data
type AppConfig struct {
	BaseData                `bson:",inline"`
	WhiteListedCards        []string `json:"whiteListedCards" bson:"whiteListedCards"`
	WhiteListedPhoneNumbers []string `json:"whiteListedPhoneNumbers" bson:"whiteListedPhoneNumbers"`
}

// CollectionName function from gmgo.Document interface
func (cfg AppConfig) CollectionName() string {
	return "appConfig"
}

//CustomerSupportConfig configuration for customer support hours.
type CustomerSupportConfig struct {
	MondayFrom    *time.Time `json:"mondayFrom" bson:"mondayFrom,omitempty"`
	MondayTo      *time.Time `json:"mondayTo" bson:"mondayTo,omitempty"`
	TuesdayFrom   *time.Time `json:"tuesdayFrom" bson:"tuesdayFrom,omitempty"`
	TuesdayTo     *time.Time `json:"tuesdayTo" bson:"tuesdayTo,omitempty"`
	WednesdayFrom *time.Time `json:"wednesdayFrom" bson:"wednesdayFrom,omitempty"`
	WednesdayTo   *time.Time `json:"wednesdayTo" bson:"wednesdayTo,omitempty"`
	ThursdayFrom  *time.Time `json:"thursdayFrom" bson:"thursdayFrom,omitempty"`
	ThursdayTo    *time.Time `json:"thursdayTo" bson:"thursdayTo,omitempty"`
	FridayFrom    *time.Time `json:"fridayFrom" bson:"fridayFrom,omitempty"`
	FridayTo      *time.Time `json:"fridayTo" bson:"fridayTo,omitempty"`
	SaturdayFrom  *time.Time `json:"saturdayFrom" bson:"saturdayFrom,omitempty"`
	SaturdayTo    *time.Time `json:"saturdayTo" bson:"saturdayTo,omitempty"`
	SundayFrom    *time.Time `json:"sundayFrom" bson:"sundayFrom,omitempty"`
	SundayTo      *time.Time `json:"sundayTo" bson:"sundayTo,omitempty"`
}

//UserConfig user configuration data
type UserConfig struct {
	BaseData  `bson:",inline"`
	UserID    string            `json:"userId" bson:"userId"`
	ConfigMap map[string]string `json:"configMap" bson:"configMap,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (uc UserConfig) CollectionName() string {
	return "userConfig"
}
