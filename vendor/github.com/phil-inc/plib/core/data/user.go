package data

import (
	"time"

	"github.com/narup/gmgo"

	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	//AccountPending - account pending email verification
	AccountPending = "EmailVerificationPending"
	//AccountActive - email verified
	AccountActive = "Active"
	//AccountLocked - account locked and cannot access resources
	AccountLocked = "Locked"
	//AccountActive - account inactive
	AccountInActive = "InActive"
	//AccountOffline - account offline
	AccountOffline = "Offline"
	//CreatorTypeUser - type of creator
	CreatorTypeUser = "User"
	//CreatorTypeSystem - type of creator
	CreatorTypeSystem = "System"
	//CreatorTypeAgent - type of creator
	CreatorTypeAgent = "Agent"

	//FreeTrialNotEligible free trial not eligible
	FreeTrialNotEligible = "NOT_ELIGIBLE"
	//FreeTrialNotApplied - free trial not applied
	FreeTrialNotApplied = "NOT_APPLIED"
	//FreeTrialApplied - free trial applied
	FreeTrialApplied = "APPLIED"
	//FreeTrialPromoApplied - free trial promo applied
	FreeTrialPromoApplied = "PROMO_APPLIED"
)

const (
	//PartnerHealthEndeavor partner string for Health Endeavor
	PartnerHealthEndeavor = "health-endeavor"
	//PartnerLemonaid partner string for lemonaid health
	PartnerLemonaid = "lemonaid"
	//PartnerPatientBank partner string for Patient data bank
	PartnerPatientBank = "patient-bank"
	//PartnerAdBloom partner string for AdBloom, marketing channel
	PartnerAdBloom = "adbloom"
	//PartnerMDOffice string that represents users that come from md funnel
	PartnerMDOffice = "md-office"
	//PartnerPushHealth string that represents users that come from push health
	PartnerPushHealth = "push-health"
)

const (
	//NotificationPrefSMS represents the notification preference for the user as SMS
	NotificationPrefSMS = "SMS"
	//NotificationPrefEmail represents the notification preference for the user as EMAIL
	NotificationPrefEmail = "EMAIL"
	//NotificationPrefEmailAndSMS means all communication channel for user notification
	NotificationPrefEmailAndSMS = "BOTH"
	//NotificationPrefVoice represents voice for user notification
	NotificationPrefVoice = "VOICE"
)

const (
	//UserTypePhilAdmin represents user type Phil who have admin rights
	UserTypePhilAdmin = "phil.admin"
	//UserTypePhilCS represents user type Phil cs agents who have admin rights
	UserTypePhilCS = "phil.cs"
	//UserTypePhilUser represents user type phil user
	UserTypePhilUser = "phil.user"
	//UserTypePhilSales represents user who are sales rep
	UserTypePhilSales = "phil.sales"
	//UserTypePhilPPAdmin represents user type Phil who are on PP admin
	UserTypePhilPPAdmin = "pp.admin"
	//UserTypePhilPPUser represents user type Phil who are on PP user
	UserTypePhilPPUser = "pp.user"
	//UserTypeMDOAdmin represetns user who are MD office admins. Usually MAs
	UserTypeMDOAdmin = "mdo.admin"
	//UserTypeMDOUser represetns user who are MD office users who may have limited views
	UserTypeMDOUser = "mdo.user"
)

//User represents phil user data
// Indexes:
// db.rexUser.createIndex({"email": "text", "fullName": "text"}) - to enable text search
type User struct {
	BaseData          `bson:",inline"`
	FullName          string                 `json:"fullName" bson:"fullName" pson:"full_name" binding:"required"`
	Email             string                 `json:"email" bson:"email" pson:"email" binding:"required"`
	AccountNumber     string                 `json:"accountNumber" bson:"accountNumber" pson:"accountNumber"`
	Password          string                 `json:"password" bson:"password" binding:"required"`
	PhoneNumber       string                 `json:"phoneNumber" bson:"phoneNumber,omitempty" pson:"phone_number"`
	FromNumber        string                 `json:"fromNumber" bson:"fromNumber,omitempty"`
	ZipCode           string                 `json:"zipCode" bson:"zipCode" pson:"zip_code" binding:"required"`
	City              string                 `json:"city" bson:"city,omitempty"`
	State             string                 `json:"state" bson:"state,omitempty" pson:"state"`
	AccountStatus     string                 `json:"accountStatus" bson:"accountStatus" pson:"account_status" binding:"required"`
	Source            string                 `json:"source" bson:"source,omitempty" pson:"source"`
	PromoCodeUsed     bool                   `json:"promoCodeUsed" bson:"promoCodeUsed"`
	FreeTrialState    string                 `json:"freeTrialState" bson:"freeTrialState,omitempty"`
	UserAttributes    *UserAttributes        `json:"userAttributes" bson:"userAttributes,omitempty"`
	NotificationPrefs *UserNotificationPrefs `json:"userNotificationPrefs" bson:"userNotificationPrefs,omitempty"`
	Promo             *mgo.DBRef             `json:"promoCode" bson:"promoCode,omitempty"`
	RoleRefs          []*mgo.DBRef           `json:"-" bson:"roles"`
	Roles             []UserRole             `json:"roles" bson:"-"`
}

// CollectionName Document interface implementation for struct User
func (u User) CollectionName() string {
	return "rexUser"
}

const funnelComplete = "funnel_complete"
const singleStepOnboarding = "single_step_onboarding"
const upsellCallCompleted = "upsell_call_completed"

const boolValueTrue = "true"
const boolValueFalse = "false"

//UserAttributes list of different user attributes
type UserAttributes struct {
	PhoneNumberVerified   bool              `json:"phoneNumberVerified" bson:"phoneNumberVerified,omitempty"`
	LoginAttemptsCount    int               `json:"loginAttemptsCount" bson:"loginAttemptsCount,omitempty"`
	LoginAttemptsExceeded *time.Time        `json:"loginAttemptsExceeded" bson:"loginAttemptsExceeded,omitempty"`
	LemonaidAttributes    map[string]string `json:"lemonaidAttributes" bson:"lemonaidAttributes,omitempty"`
	AdBloomAttributes     map[string]string `json:"adBloomAttributes" bson:"adBloomAttributes,omitempty"`
	MdFunnelAttributes    map[string]string `json:"mdFunnelAttributes" bson:"mdFunnelAttributes,omitempty"`
	ProcessingAttributes  map[string]string `json:"processingAttributes" bson:"processingAttributes,omitempty"`
}

//UserNotificationPrefs stores all the MD office user preferences. Each preference has 4 possible values
//NotificationPrefSMS, NotificationPrefEmail, NotificationPrefEmailAndSMS, and NotificationPrefVoice
type UserNotificationPrefs struct {
	Blocked           bool   `json:"blocked" bson:"blocked"`
	RxUpdatesPref     string `json:"rxUpdatesPref" bson:"rxUpdatesPref,omitempty"`
	AnnouncementsPref string `json:"announcementsPref" bson:"announcementsPref,omitempty"`
}

//CanSendSMS returns true if user notification preferences allows SMS
func (u User) CanSendSMS() bool {
	if u.NotificationPrefs == nil || u.NotificationPrefs.RxUpdatesPref == "" {
		return true
	}
	return u.NotificationPrefs.RxUpdatesPref == "SMS" || u.NotificationPrefs.RxUpdatesPref == "BOTH"
}

//CanSendEmail returns true if user notification preferences allows Email
func (u User) CanSendEmail() bool {
	if u.NotificationPrefs == nil || u.NotificationPrefs.RxUpdatesPref == "" {
		return true
	}
	return u.NotificationPrefs.RxUpdatesPref == "EMAIL" || u.NotificationPrefs.RxUpdatesPref == "BOTH"
}

//IsPharmacyUser returns true if the user has a ROLE_PHARMACIST
func (u User) IsPharmacyUser() bool {
	if u.IsAdmin() {
		return false
	}
	for _, role := range u.Roles {
		if role.Name == "ROLE_PHARMACIST" {
			return true
		}
	}
	return false
}

//IsAdmin returns true if user has role ROLE_ADMIN
func (u User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

//IsSysAdmin returns true if user has role ROLE_SYSADMIN
func (u User) IsSysAdmin() bool {
	for _, role := range u.Roles {
		if role.Name == "ROLE_SYSADMIN" {
			return true
		}
	}
	return false
}

//IsSuperUser returns true if user has role ROLE_SUPER_USER
func (u User) IsSuperUser() bool {
	for _, role := range u.Roles {
		if role.Name == "ROLE_SUPER_USER" {
			return true
		}
	}
	return false
}

//IsMDUser returns true if user has role ROLE_MD
func (u User) IsMDUser() bool {
	for _, role := range u.Roles {
		if role.Name == "ROLE_MD" {
			return true
		}
	}
	return false
}

//FreeTrialEligibile returns if user is eligible for free trial or not
func (u User) FreeTrialEligibile() bool {
	if u.PromoCodeUsed {
		return false
	}
	if u.FreeTrialState == FreeTrialNotApplied {
		return true
	}
	return false
}

//IsFunnelComplete returns true if user has completed the initial funnel
func (u User) IsFunnelComplete(partner string) bool {
	if u.UserAttributes == nil {
		return false
	}
	if partner == PartnerMDOffice && u.UserAttributes.MdFunnelAttributes != nil {
		return u.UserAttributes.MdFunnelAttributes[funnelComplete] == boolValueTrue
	}
	if partner == PartnerLemonaid && u.UserAttributes.LemonaidAttributes != nil {
		return u.UserAttributes.LemonaidAttributes[funnelComplete] == boolValueTrue
	}
	return false
}

//IsTempEmail checks if the user's email is a temporarily assigned email
func (u User) IsTempEmail() bool {
	return strings.HasSuffix(u.Email, "-temp@phil.us")
}

//IsAccountInactive Is the user Inactive
func (u User) IsAccountInactive() bool {
	return u.AccountStatus == AccountInActive
}

//IsAccountOffline Is the user Offline
func (u User) IsAccountOffline() bool {
	return u.AccountStatus == AccountOffline
}

//IsMDFunnelSingleStepOnboarding returns true if user is setup for 1-step onboarding funnel/process
//1-step onboarding involves account creation and payment approval at the same time.
func (u User) IsMDFunnelSingleStepOnboarding() bool {
	if u.UserAttributes == nil {
		return false
	}
	if u.UserAttributes.MdFunnelAttributes != nil {
		return u.UserAttributes.MdFunnelAttributes[singleStepOnboarding] == "true"
	}
	return false
}

//IsUpSellCallCompleted returns true if upsell call to the user has been made
func (u User) IsUpSellCallCompleted() bool {
	if u.UserAttributes == nil {
		return false
	}
	if u.UserAttributes.ProcessingAttributes == nil {
		return false
	}
	return u.UserAttributes.ProcessingAttributes[upsellCallCompleted] == boolValueTrue
}

//MarkUpsellCallCompleted saves the attribute for user with upsell call status
func (u *User) MarkUpsellCallCompleted() {
	if u.UserAttributes == nil {
		u.UserAttributes = new(UserAttributes)
	}
	if u.UserAttributes.ProcessingAttributes == nil {
		u.UserAttributes.ProcessingAttributes = make(map[string]string)
	}
	u.UserAttributes.ProcessingAttributes[upsellCallCompleted] = boolValueTrue
}

//UserRole user role data
type UserRole struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
}

// CollectionName function to implement Document interface
func (userRole UserRole) CollectionName() string {
	return "role"
}

// UserSource represents user source and campaign data
type UserSource struct {
	BaseData           `bson:",inline"`
	UserID             string            `json:"userId" bson:"userId" pson:"userId" binding:"required"`
	Source             string            `json:"source" bson:"source,omitempty" pson:"source"`
	Channel            string            `json:"channel" bson:"channel,omitempty" pson:"channel"`
	CampaignParameters map[string]string `json:"campaignParameters" bson:"campaignParameters,omitempty" pson:"campaign_parameters"`
	LoggedInFromApp    bool              `json:"loggedInFromApp" bson:"loggedInFromApp" pson:"logged_in_from_app"`
}

// CollectionName function to implement Document interface
func (us UserSource) CollectionName() string {
	return "userSource"
}

// UserToken - stores user refresh token info
type UserToken struct {
	BaseData        `bson:",inline"`
	UserID          string     `json:"userId" bson:"userId" binding:"required"`
	RefreshToken    string     `json:"refreshToken" bson:"refreshToken" binding:"required"`
	BlackListed     bool       `json:"blackListed" bson:"blackListed"`
	TokenExpiration *time.Time `json:"-" bson:"tokenExpiration,omitempty"`
}

// CollectionName function to implement Document interface
func (ut UserToken) CollectionName() string {
	return "userToken"
}

//SessionToken temporary token data used for session migration between different
//Phil apps. Token more than 5 mins old are considered inactive.
type SessionToken struct {
	BaseData `bson:",inline"`
	UserID   string `json:"userId" bson:"userId" binding:"required"`
	Token    string `json:"token" bson:"token" binding:"required"` //Indexed
}

// CollectionName function to implement Document interface
func (st SessionToken) CollectionName() string {
	return "sessionToken"
}

//IdentityToken anonymous user token to identify user based on ID
type IdentityToken struct {
	BaseData `bson:",inline"`
	UserID   string `json:"userId" bson:"userId" binding:"required"`
	Token    string `json:"token" bson:"token" binding:"required"` //Indexed
	Status   string `json:"status" bson:"status" binding:"required"`
}

// CollectionName function to implement Document interface
func (st IdentityToken) CollectionName() string {
	return "identityToken"
}

// Patient represents patient data
type Patient struct {
	BaseData            `bson:",inline"`
	PatientName         string     `json:"patientName" bson:"patientName" pson:"patient_name" binding:"required"` //Indexed
	PhoneNumber         string     `json:"phoneNumber" bson:"phoneNumber" pson:"phone_number" binding:"required"`
	IsPrimaryCardHolder string     `json:"isPrimaryCardHolder" bson:"isPrimaryCardHolder,omitempty"`
	Gender              string     `json:"gender" bson:"gender,omitempty"`
	Allergies           []string   `json:"allergies" bson:"allergies,omitempty" pson:"allergies"`
	DateOfBirth         *time.Time `json:"dateOfBirth" bson:"dateOfBirth" pson:"date_of_birth" binding:"required"` //Indexed
	ManagerRef          *mgo.DBRef `json:"-" bson:"manager,omitempty" pson:"manager"`
	PaymentOption       string     `json:"paymentOption" bson:"paymentOption,omitempty"`
	Manager             *User      `json:"manager,omitempty" bson:"-"`
	DefaultAddressID    string     `json:"defaultAddressId, omitempty" bson:"defaultAddressId"`
	DefaultInsuranceID  string     `json:"defaultInsuranceId,omitempty" bson:"defaultInsuranceId"`
	DefaultPaymentID    string     `json:"defaultPaymentId,omitempty" bson:"defaultPaymentId"`
}

// CollectionName function to implement Document interface for patient
func (p Patient) CollectionName() string {
	return "patient"
}

// SetDefaultAddress sets the default address of patient
func (p *Patient) SetDefaultAddress(addressID string) {
	p.DefaultAddressID = addressID
}

// SetDefaultInsurance sets the default address of patient
func (p *Patient) SetDefaultInsurance(insuranceID string) {
	p.DefaultInsuranceID = insuranceID
}

// SetDefaultPaymentInfo sets the default payment info of patient
func (p *Patient) SetDefaultPaymentInfo(paymentInfoID string) {
	p.DefaultPaymentID = paymentInfoID
}

// PatientNote represents patient note data
type PatientNote struct {
	BaseData  `bson:",inline"`
	PatientID string `json:"patientId" bson:"patientId" binding:"required"`
	AgentID   string `json:"agentId" bson:"agentId" binding:"required"`
	AgentName string `json:"agentName" bson:"agentName" binding:"required"`
	Message   string `json:"message" bson:"message,omitempty"`
	Type      string `json:"type" bson:"type,omitempty"`
	Important bool   `json:"important" bson:"important"`
}

// CollectionName function to implement Document interface for promo
func (pn PatientNote) CollectionName() string {
	return "patientNote"
}

// Promo data
type Promo struct {
	BaseData        `bson:",inline"`
	Code            string `json:"code" bson:"code" binding:"required"`
	Credit          string `json:"credit" bson:"credit,omitempty"`
	CopayCredit     string `json:"copayCredit" bson:"copayCredit,omitempty"`
	CashCredit      string `json:"cashCredit" bson:"cashCredit,omitempty"`
	ShippingCredit  string `json:"shippingCredit" bson:"shippingCredit,omitempty"`
	Description     string `json:"description" bson:"description,omitempty"`
	Campaign        string `json:"campaign" bson:"campaign,omitempty"`
	Active          bool   `json:"active" bson:"active"`
	MultipleSupport bool   `json:"multipleSupport" bson:"multipleSupport"`
	Recurring       bool   `json:"recurring" bson:"recurring"`
}

// CollectionName function to implement Document interface for promo
func (p Promo) CollectionName() string {
	return "promoCode"
}

// EmailVerificationToken email verification token
type EmailVerificationToken struct {
	BaseData `bson:",inline"`
	Email    string `json:"email" bson:"email"`
	Token    string `json:"token" bson:"token"`
}

// CollectionName function to implement Document interface for promo
func (e EmailVerificationToken) CollectionName() string {
	return "emailVerificationToken"
}

// PasswordResetToken stores password reset token
type PasswordResetToken struct {
	BaseData `bson:",inline"`
	Token    string     `json:"token" bson:"token"`
	User     string     `json:"user" bson:"-"`
	UserRef  *mgo.DBRef `json:"-" bson:"user,omitempty"`
}

// CollectionName function to implement Document interface for promo
func (pr PasswordResetToken) CollectionName() string {
	return "passwordResetToken"
}

// ZipUnsupportedUser record for zip unsupported user
type ZipUnsupportedUser struct {
	BaseData  `bson:",inline"`
	Email     string `json:"email" bson:"email"`
	FullName  string `json:"fullName" bson:"fullName"`
	Zip       string `json:"zip" bson:"zip"`
	EmailSent bool   `json:"emailSent" bson:"emailSent"`
}

// CollectionName function to implement Document interface for promo
func (z ZipUnsupportedUser) CollectionName() string {
	return "zipUnsupportedUser"
}

// CopyValue copies user role value
func (userRole *UserRole) CopyValue() UserRole {
	copyRole := UserRole{}
	copyRole.ID = userRole.ID
	copyRole.Name = userRole.Name
	copyRole.Description = userRole.Description

	return copyRole
}

//TwoStepToken represents intermediate 2-factor authentication token.
type TwoStepToken struct {
	BaseData     `bson:",inline"`
	Code         string `json:"code" bson:"code"`
	UserID       string `json:"userId" bson:"userId"`
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

//CollectionName returns collection name for 2-factor auth token collection
func (tst TwoStepToken) CollectionName() string {
	return "twoStepToken"
}

// ManagerByRef returns user data based on database reference
func ManagerByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*User, error) {
	user := new(User)
	if err := session.FindByRef(dbRef, user); err != nil {
		return nil, err
	}

	return user, nil
}

// PatientByRef returns patient data based on given patient reference
func PatientByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Patient, error) {
	pt := new(Patient)
	if err := session.FindByRef(dbRef, pt); err != nil {
		return nil, err
	}

	return pt, nil
}
