package data

import (
	"time"

	"fmt"

	"strings"

	"github.com/narup/gmgo"
	mgo "gopkg.in/mgo.v2"
)

const (
	//Transfer represents order origin type transfer
	Transfer = "Transfer"
	//TransferAll represents medication name for transfer all
	TransferAll = "Transfer All"
	//NewMedication represents medication name when it's unknown
	NewMedication = "New Medication"
	//PaperPrescription medication name for paper rx when name is unknown
	PaperPrescription = "Paper Prescription"
)
const (
	//PayThroughInsurance payment option select as pay with insurance
	PayThroughInsurance = "PayThroughInsurance"
	//PayDirectly payment option as cash
	PayDirectly = "PayDirectly"
)
const (
	//OriginTransfer represents transfer origin type
	OriginTransfer = "Transfer"
	//OriginFromDoctorDirect origin type for eScript Rx directly from the doctor
	OriginFromDoctorDirect = "FromDoctorDirect"
	//OriginFromDoctorPaper  origin type for paper prescription
	OriginFromDoctorPaper = "FromDoctorPaper"
	//RxSourceJumpStart source type JumpStart
	RxSourceJumpStart = "jump-start"
	//RxSourcePushHealth source type PushHealth
	RxSourcePushHealth = "push-health"
	//RxSourceMDOffice source lemonaid
	RxSourceMDOffice = "lemonaid"
)
const (
	//OriginSubType_MDOffice_MDDashboard order was creating from MD partner office and MD Dashboard
	OriginSubType_MDOffice_MDDashboard = "MDOffice_MDDashboard"
	//OriginSubType_MDOffice_CSDashboard_Phil order was sent from MD Partner office and CS Dashboard through Phil
	OriginSubType_MDOffice_CSDashboard_Phil = "MDOffice_CSDashboard_Phil"
	//OriginSubType_MDOffice_CSDashboard_PostMeds order was sent from MD Partner office and CS Dashboard through Postmeds
	OriginSubType_MDOffice_CSDashboard_PostMeds = "MDOffice_CSDashboard_PostMeds"
)

const (
	//RxAttributesLemonaidRefillExpired atrribute key to represent lemonaid Rx as refill expired.
	RxAttributesLemonaidRefillExpired = "lemonaid-refill-expired"
)

// Prescription data representation
type Prescription struct {
	BaseData                   `bson:",inline"`
	Name                       string            `json:"name" bson:"name,omitempty" pson:"name"`
	Number                     string            `json:"number" bson:"number,omitempty"`
	Form                       string            `json:"form" bson:"form,omitempty" pson:"form"`
	Strength                   string            `json:"strength" bson:"strength,omitempty" pson:"strength"`
	Sig                        string            `json:"sig" bson:"sig,omitempty"`
	OriginalPrescriptionID     string            `json:"originalPrescriptionId" bson:"originalPrescriptionId,omitempty"`
	OldRxNumbers               []string          `json:"oldRxNumbers" bson:"oldRxNumbers,omitempty"`
	MedicationQuantity         int               `json:"medicationQuantity" bson:"medicationQuantity,omitempty" pson:"medication_quantity"`
	DaysOfSupply               int               `json:"daysOfSupply" bson:"daysOfSupply" pson:"days_of_supply"`
	Avatar                     string            `json:"avatar" bson:"avatar,omitempty"`
	OriginType                 string            `json:"originType" bson:"originType,omitempty" pson:"origin_type"` //Possible values: Transfer, PaperPrescription, FromDoctorDirect
	OriginSubType              string            `json:"originSubType" bson:"originSubType,omitempty"`
	ControlledDrugsOption      string            `json:"controlledDrugsOption" bson:"controlledDrugsOption,omitempty"`
	MedicationPrice            string            `json:"medicationPrice" bson:"medicationPrice,omitempty" pson:"medication_price"`
	PreviousMedicationPrice    string            `json:"previousMedicationPrice" bson:"previousMedicationPrice,omitempty"`
	Archived                   bool              `json:"archived" bson:"archived" pson:"archived"`
	Suspended                  bool              `json:"suspended" bson:"suspended" pson:"suspended"`
	Pending                    bool              `json:"pending" bson:"pending" pson:"pending"`
	CopayPreApproved           bool              `json:"copayPreApproved" bson:"copayPreApproved"`
	CurrentOrderID             string            `json:"currentOrderId" bson:"currentOrderId" pson:"current_order_id"` //Sparse index
	Promo                      string            `json:"promo" bson:"promo,omitempty" pson:"promo"`
	GlobalPromo                string            `json:"globalPromo" bson:"globalPromo,omitempty" pson:"global_promo"`
	BatchOrderID               string            `json:"batchOrderID" bson:"batchOrderID,omitempty"`
	BundleID                   string            `json:"bundleID" bson:"bundleID,omitempty"`
	ArchivedReason             string            `json:"archivedReason" bson:"archivedReason,omitempty" pson:"archived_reason"`
	Source                     string            `json:"source" bson:"source,omitempty" pson:"source"`
	RefillAuthorizedBy         string            `json:"refillAuthorizedBy" bson:"refillAuthorizedBy,omitempty"`
	ScheduleDate               *time.Time        `json:"scheduleDate" bson:"scheduleDate,omitempty"`
	ArchivedDate               *time.Time        `json:"archivedDate" bson:"archivedDate,omitempty" pson:"archived_date"`
	LastFillDate               *time.Time        `json:"lastFillDate" bson:"lastFillDate,omitempty" pson:"last_fill_date"`
	ReFillDate                 *time.Time        `json:"refillDate" bson:"refillDate,omitempty" pson:"refill_date"`
	ShipDate                   *time.Time        `json:"shipDate" bson:"shipDate,omitempty"` //ship date is a display only date to communicate to the user. It's usually 3 days after refill date.
	PendingSinceDate           *time.Time        `json:"pendingSinceDate" bson:"pendingSinceDate,omitempty" pson:"pending_since_date"`
	CurrentReFillTriggeredDate *time.Time        `json:"currentRefillTriggeredDate" bson:"currentRefillTriggeredDate,omitempty"`
	RefillsRemaining           int16             `json:"refillsRemaining" bson:"refillsRemaining" pson:"refills_remaining"`
	Coupon                     *Coupon           `json:"coupon" bson:"coupon,omitempty"`
	ManagerRef                 *mgo.DBRef        `json:"-" bson:"manager,omitempty" pson:"manager"`
	PatientRef                 *mgo.DBRef        `json:"-" bson:"patient,omitempty" pson:"patient"`
	DoctorRef                  *mgo.DBRef        `json:"-" bson:"doctor,omitempty" pson:"doctor"`
	InsuranceRef               *mgo.DBRef        `json:"-" bson:"insurance,omitempty" pson:"insurance"`
	CurrentOrderRef            *mgo.DBRef        `json:"-" bson:"currentOrder,omitempty"`
	Manager                    *User             `json:"manager" bson:"-"`
	Patient                    *Patient          `json:"patient" bson:"-"`
	Doctor                     *Doctor           `json:"doctor" bson:"-"`
	Insurance                  *Insurance        `json:"insurance" bson:"-"`
	CurrentOrder               *Order            `json:"currentOrder" bson:"-"`
	EstimatedPrice             string            `json:"estimatedPrice" bson:"estimatedPrice"`
	Attributes                 map[string]string `json:"attributes" bson:"attributes,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (rx Prescription) CollectionName() string {
	return "prescription"
}

//IsNameTransferAll checks if rx name indicates Transfer all
func (rx Prescription) IsNameTransferAll() bool {
	return rx.Name == "Transfer all" || rx.Name == "Transfer All" || rx.Name == "TransferAll"
}

//HasValidInsurance check if Rx has valid insurance
func (rx Prescription) HasValidInsurance() bool {
	if rx.Insurance == nil {
		return false
	}
	if (rx.Insurance.InsuranceID != "" && rx.Insurance.BinNumber != "") || rx.Insurance.CardImageID != "" {
		return true
	}
	return false
}

//FullRxName returns complete Rx name
func (rx Prescription) FullRxName() string {
	if rx.Name == "" && rx.OriginType == OriginFromDoctorDirect {
		return "New Prescription"
	}
	if rx.Name == "" && rx.OriginType == OriginFromDoctorPaper {
		return "Paper Prescription"
	}

	rxName := rx.Name
	if rx.Form != "" {
		rxName = fmt.Sprintf("%s %s", rxName, rx.Form)
	}
	if rx.Strength != "" {
		rxName = fmt.Sprintf("%s %s", rxName, rx.Strength)
	}
	return strings.ToUpper(rxName)
}

//IsScheduledForLater checks if Rx is scheduled for later date
func (rx *Prescription) IsScheduledForLater() bool {
	if rx.ScheduleDate != nil {
		sd := *rx.ScheduleDate
		now := time.Now().UTC()
		todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
		if sd.After(todayEnd) {
			return true
		}
	}
	return false
}

//GetEarliestRefillDate returns the date at which Rx can be filled earliest
func (rx *Prescription) GetEarliestRefillDate() *time.Time {
	earliestRefillDate := rx.CurrentOrder.EarliestRefillDate()

	if rx.LastFillDate != nil && earliestRefillDate == nil {
		refillDate := rx.LastFillDate.Add(time.Duration(0.8 * float64(24) * float64(rx.DaysOfSupply)))
		return &refillDate
	}
	return earliestRefillDate
}

//AppliedFreeTrial applied free trial
func (rx Prescription) AppliedFreeTrial() bool {
	return rx.Manager.FreeTrialState == "APPLIED" && rx.CurrentOrder.PaymentMilestone.AppliedFreeTrialCharge
}

//InsuranceMissingOrException checks if the prescription order has missing
//insurance or has any insurance exception
func (rx Prescription) InsuranceMissingOrException() bool {
	if !rx.CurrentOrder.IsPayThroughInsurance() {
		return false
	}
	if rx.Insurance == nil {
		return true
	}
	if rx.CurrentOrder.MissingInfoState == nil {
		return false
	}
	o := rx.CurrentOrder
	missingStatus := *o.MissingInfoState.MissingStatus
	if missingStatus == InsuranceExpired ||
		missingStatus == InvalidInsurance ||
		missingStatus == NotCovered ||
		missingStatus == MissingInsurance {
		return true
	}
	return false
}

//IsInsuranceMissing returns true if prescription has no insurance
func (rx Prescription) IsInsuranceMissing() bool {
	if rx.CurrentOrder == nil {
		return false
	}
	return rx.CurrentOrder.IsPayThroughInsurance() && rx.InsuranceRef == nil
}

//IsRxFromPartner returns true if Rx is from the given partner
func (rx Prescription) IsRxFromPartner(partner string) bool {
	return rx.Attributes != nil && rx.Attributes["source"] == partner
}

//IsRxFromMDChannel returns true if Rx order is from MD channel
func (rx Prescription) IsRxFromMDChannel() bool {
	return rx.Attributes != nil && rx.Attributes["source"] == PartnerMDOffice
}

// Coupon data representation used by prescription for discounts
type Coupon struct {
	BinNumber   string `json:"binNumber" bson:"binNumber,omitempty"`
	CouponID    string `json:"couponId" bson:"couponId,omitempty"`
	PcnNumber   string `json:"pcnNumber" bson:"pcnNumber,omitempty"`
	GroupNumber string `json:"groupNumber" bson:"groupNumber,omitempty"`
	Active      bool   `json:"active" bson:"active"`
}

//RxComment the rx log message
type RxComment struct {
	BaseData  `bson:",inline"`
	RxID      string            `json:"rxId" bson:"rxId"`
	AgentID   string            `json:"agentId" bson:"agentId"`
	AgentName string            `json:"agentName" bson:"agentName"`
	AgentType string            `json:"agentType" bson:"agentType,omitempty"`
	Message   string            `json:"message" bson:"message,omitempty"`
	Important bool              `json:"important" bson:"important"`
	Type      string            `json:"type" bson:"type,omitempty"`
	PatientID string            `json:"patientId" bson:"patientId,omitempty"`
	RxList    map[string]string `json:"rxList" bson:"rxList,omitempty"`
}

//CollectionName function from gmgo.Document interface
func (m RxComment) CollectionName() string {
	return "rxComment"
}

// RefillReport data representation for refill reports
type RefillReport struct {
	BaseData              `bson:",inline"`
	NumberOfPrescriptions int16      `json:"numberOfPrescriptions" bson:"numberOfPrescriptions"`
	SkipRefillEmailCount  int16      `json:"skipRefillEmailCount" bson:"skipRefillEmailCount,omitempty"`
	RxIDs                 []string   `json:"rxIds" bson:"rxIds,omitempty"`
	OrderNumbers          []string   `json:"orderNumbers" bson:"orderNumbers,omitempty"`
	Type                  string     `json:"type" bson:"type"`
	ExecutionLocalDate    *time.Time `json:"executionLocalDate" bson:"executionLocalDate"`
}

//CollectionName function from gmgo.Document interface
func (rr RefillReport) CollectionName() string {
	return "refillReport"
}

//SkipRefillToken represents skip refill token storage
type SkipRefillToken struct {
	BaseData            `bson:",inline"`
	Token               string     `json:"token" bson:"token"`
	RxID                string     `json:"rxId" bson:"rxId"`
	State               string     `json:"state" bson:"state"`
	MedicationName      string     `json:"medicationName" bson:"medicationName"`
	UpcomingRefillDate  *time.Time `json:"upcomingRefillDate" bson:"upcomingRefillDate"`
	RefillDateIfSkipped *time.Time `json:"refillDateIfSkipped" bson:"refillDateIfSkipped"`
}

//CollectionName function from gmgo.Document interface
func (sr SkipRefillToken) CollectionName() string {
	return "skipRefillToken"
}

//SMSLog represents SMS log
type SMSLog struct {
	BaseData     `bson:",inline"`
	ManagerID    string            `json:"managerId" bson:"managerId,omitempty"`
	MessageSid   string            `json:"messageSid" bson:"messageSid,omitempty"`
	PhoneNumber  string            `json:"phoneNumber" bson:"phoneNumber,omitempty"`
	ReplyPending bool              `json:"replyPending" bson:"replyPending,omitempty"`
	Params       map[string]string `json:"params" bson:"params,omitempty"`
}

//CollectionName function from gmgo.Document interface
func (sl SMSLog) CollectionName() string {
	return "sMSLog"
}

// DoctorByRef returns doctor data based on reference
func DoctorByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Doctor, error) {
	d := new(Doctor)
	if err := session.FindByRef(dbRef, d); err != nil {
		return nil, err
	}

	return d, nil
}

// InsuranceByRef returns doctor data based on reference
func InsuranceByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Insurance, error) {
	i := new(Insurance)
	if err := session.FindByRef(dbRef, i); err != nil {
		return nil, err
	}

	return i, nil
}

//InitializePrescription initialize new prescriptions based on given Rx params
func InitializePrescription(name string, originType string, avatar string, scheduleDate *time.Time) *Prescription {

	prescription := new(Prescription)
	prescription.InitData()
	prescription.ScheduleDate = scheduleDate
	prescription.Name = name
	prescription.OriginType = originType
	prescription.Avatar = avatar
	prescription.MedicationPrice = "0"
	prescription.PreviousMedicationPrice = "-1"
	prescription.RefillsRemaining = -1
	return prescription
}

//AddPatient adds patient to Rx
func (rx *Prescription) AddPatient(patient *Patient) {
	if patient != nil {
		rx.Patient = patient
		rx.PatientRef = DBRef(rx.Patient.CollectionName(), rx.Patient.ID)
	}
}

//AddDoctor adds Doctor info to Rx
func (rx *Prescription) AddDoctor(doctor *Doctor) {
	if doctor != nil {
		rx.Doctor = doctor
		rx.DoctorRef = DBRef(rx.Doctor.CollectionName(), rx.Doctor.ID)
	}
}

//AddOrder assigns order as current order.
func (rx *Prescription) AddOrder(order *Order) {
	if order != nil {
		rx.CurrentOrder = order
		rx.CurrentOrderRef = DBRef(rx.CurrentOrder.CollectionName(), rx.CurrentOrder.ID)
		rx.CurrentOrder.RxID = rx.StringID()
		rx.CurrentOrderID = rx.CurrentOrder.StringID()
	}
}

//AddManager assigns manager to prescription
func (rx *Prescription) AddManager(user *User) {
	if user != nil {
		rx.Manager = user
		rx.ManagerRef = DBRef(rx.Manager.CollectionName(), rx.Manager.ID)
	}

}

func (rx *Prescription) SetPaymentOption(outOfPocket bool) {
	if rx.CurrentOrder != nil && outOfPocket {
		rx.CurrentOrder.PaymentOption = PayDirectly
	} else if rx.CurrentOrder != nil {
		rx.CurrentOrder.PaymentOption = PayThroughInsurance
	}
}

func (rx *Prescription) SetInsuranceType(govtSponsoredInsurance bool) {
	if rx.CurrentOrder != nil && rx.CurrentOrder.InsuranceMilestone != nil {
		if govtSponsoredInsurance {
			rx.CurrentOrder.InsuranceMilestone.IsFederalSponsoredInsurance = Yes
		} else {
			rx.CurrentOrder.InsuranceMilestone.IsFederalSponsoredInsurance = No
		}
	}
}

func (rx *Prescription) SetMissingInfo(status MissingStatus) {
	if rx.CurrentOrder != nil {
		rx.CurrentOrder.SetMissingInfoState(status)
	}
}

func (rx *Prescription) HasPaymentError() bool {
	if !rx.CurrentOrder.IsPaymentAuthorized() {
		return false
	}
	if rx.IsScheduledForLater() || rx.Archived || rx.Suspended {
		return false
	}
	return rx.CurrentOrder.HasPaymentError()
}

func (rx *Prescription) IsMissingPayment() bool {
	if !rx.CurrentOrder.IsPaymentAuthorized() {
		return false
	}
	if rx.IsScheduledForLater() || rx.Archived || rx.Suspended {
		return false
	}
	return rx.CurrentOrder.Payment == nil
}
