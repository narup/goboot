package data

import (
	"time"

	"github.com/narup/gmgo"
	mgo "gopkg.in/mgo.v2"
)

const (
	//FirstFill repreresnts order fill type as first ever fill with Phil
	FirstFill = "FirstFill"
	//ReFill reprerents refill order
	ReFill = "ReFill"
)

const (
	// EmailTransferDelayed notification stat type for transfer delay email notification
	EmailTransferDelayed = "EMAIL_TRANSFER_DELAYED"
	//EmailRefillAuthDelayed notification type for email notification sent if refill auth is delayed
	EmailRefillAuthDelayed = "EMAIL_REFILL_AUTH_DELAYED"
	//EmailRefillAuthRequested notification type for email notification sent if refill is requested
	EmailRefillAuthRequested = "EMAIL_REFILL_AUTH_REQUESTED"
	//EmailRefillAuthNoResponse notification type for email notification sent if refill auth request has no response from MD
	EmailRefillAuthNoResponse = "EMAIL_REFILL_AUTH_NO_RESPONSE"
	//SMSMissingInsuranceSent notification type for sms notification sent if insurance is missing
	SMSMissingInsuranceSent = "SMS_MISSING_INSURANCE_SENT"
	//EmailSkipRefillNotificationEmail notification type for skip refill notification
	EmailSkipRefillNotificationEmail = "EMAIL_SKIP_REFILL_NOTIFICATION"
)

const (
	OneDayMail      = "OneDayMail"
	TwoDayMail      = "TwoDayMail"
	SameDayDelivery = "SameDayDelivery"
)

var ShippingOptions = map[string]string{TwoDayMail: "USPS Priority Mail", OneDayMail: "USPS Overnight Express"}

// OrderSleepState that defines possible sleep states for rx order.
type OrderSleepState string

const (
	//CallPatient - represents orders that involves calling patient
	CallPatient OrderSleepState = "CallPatient"
	//ScrubAndRoute - represents order during scrub and transfer process
	ScrubAndRoute OrderSleepState = "ScrubAndRoute"
	//ScrubbedAndRouted - sleep state after routing and waiting for OP to take action
	ScrubbedAndRouted OrderSleepState = "ScrubbedAndRouted"
	//CheckFax - sleep state after fax is sent to MD office
	CheckFax OrderSleepState = "CheckFax"
	//CallMD - sleep state after MD office is called directly
	CallMD OrderSleepState = "CallMD"
	//PaymentApproval - sleep state for payment approval
	PaymentApproval OrderSleepState = "PaymentApproval"
	//PaymentApprovalUntilFreeTrialDelivered - sleep state for payment approval until free trial delivered
	PaymentApprovalUntilFreeTrialDelivered OrderSleepState = "PaymentApprovalUntilFreeTrialDelivered"
	//ProcessOrder sleep state before order is processed after routing
	ProcessOrder OrderSleepState = "ProcessOrder"
	//TransferFaxFollowUp state once transfer request is sent
	TransferFaxFollowUp OrderSleepState = "TransferFaxFollowUp"
	//DelayedEmailSent state after email is sent for delayed order. Note: Not used
	DelayedEmailSent OrderSleepState = "DelayedEmailSent"
	//DelayedCalledCustomer state after patient is called for delayed order
	DelayedCalledCustomer OrderSleepState = "DelayedCallCustomer"
	//TimedSleep sleep state for specified amount of time
	TimedSleep OrderSleepState = "TimedSleep"
)

// MissingStatus represents missing status type
type MissingStatus string

const (
	//MissingInfo - represents when order is missing information
	MissingInfo MissingStatus = "MISSING_INFO"
	//Exception - represents when order has an exception
	Exception                 MissingStatus = "EXCEPTION"
	NotCovered                MissingStatus = "NOT_COVERED"
	MissingInsurance          MissingStatus = "MISSING_INSURANCE"
	InvalidInsurance          MissingStatus = "INVALID_INSURANCE"
	InsuranceExpired          MissingStatus = "INSURANCE_EXPIRED"
	NoPatient                 MissingStatus = "NO_PATIENT"
	RefillsDenied             MissingStatus = "REFILLS_DENIED"
	PriorAuthDenied           MissingStatus = "PRIOR_AUTH_DENIED"
	StepTherapyDenied         MissingStatus = "STEP_THERAPY_DENIED"
	NewRxDenied               MissingStatus = "NEW_RX_DENIED"
	PaymentError              MissingStatus = "PAYMENT_ERROR"
	CopayIncreased            MissingStatus = "COPAY_INCREASED"
	CopayDecreased            MissingStatus = "COPAY_DECREASED"
	OriginalCopayNotAvailable MissingStatus = "ORIGINAL_COPAY_NOTAVAILABLE"
	InvalidOpNoPatient        MissingStatus = "INVALID_OP_NO_PATIENT"
	InvalidOpNoTransferableRx MissingStatus = "INVALID_OP_NO_TRANSFERABLE_RX"
	OpRefusal                 MissingStatus = "OP_REFUSAL"
	C2Drugs                   MissingStatus = "C2_DRUGS"
	InsuranceOverride         MissingStatus = "INSURANCE_OVERRIDE"
	DeliveryFailed            MissingStatus = "DELIVERY_FAILED"
	DeliveryReturned          MissingStatus = "DELIVERY_RETURNED"
	StockException            MissingStatus = "STOCK_EXCEPTION"
)

// MissingField represents different fields
type MissingField string

const (
	MissingInsuranceField  MissingField = "INSURANCE"
	MissingDateOfBirth     MissingField = "DOB"
	MissingDeliveryAddress MissingField = "ADDRESS"
	MissingPayment         MissingField = "PAYMENT"
	MissingDoctor          MissingField = "DOCTOR"
	MissingMRN             MissingField = "MRN"
	MissingContactPhone    MissingField = "CONTACT_PHONE"
	MissingOtherInfo       MissingField = "OTHER"
)

//InsuranceExceptionType represents different insurance exceptions
type InsuranceExceptionType string

const (
	InsExInvalidCardHolderId     InsuranceExceptionType = "INVALID_CARD_HOLDER_ID"
	InsExInvalidGroup            InsuranceExceptionType = "INVALID_GROUP"
	InsExInvalidDOB              InsuranceExceptionType = "INVALID_DOB"
	InsExInvalidGenderCode       InsuranceExceptionType = "INVALID_GENDER_CODE"
	InsExInvalidPersonCode       InsuranceExceptionType = "INVALID_PERSON_CODE"
	InsExPriorAuthRequired       InsuranceExceptionType = "PRIOR_AUTH_REQUIRED"
	InsExStepTherapyRequired     InsuranceExceptionType = "STEP_THERAPY_REQUIRED"
	InsExNonFormularyMed         InsuranceExceptionType = "NON_FORMULARY_MED"
	InsExPlanLimitationsExceeded InsuranceExceptionType = "PLAN_LIMITATIONS_EXCEEDED"
	InsExCoverageTerminated      InsuranceExceptionType = "COVERAGE_TERMINATED"
	InsExMedicarePartBRequired   InsuranceExceptionType = "MEDICARE_PART_B_REQUIRED"
	InsExSubmitToOtherProcessor  InsuranceExceptionType = "SUBMIT_TO_OTHER_PROCESSOR"
	InsExPharmacyNotContracted   InsuranceExceptionType = "PHARMACY_NOT_CONTRACTED"
	InsExNonPreferredPharmacy    InsuranceExceptionType = "NON_PREFERRED_PHARMACY"
	InsExMustUseMailOrder        InsuranceExceptionType = "MUST_USE_MAIL_ORDER"
	InsExMDNotInNetwork          InsuranceExceptionType = "MD_NOT_IN_NETWORK"
	InsExMissingInsurance        InsuranceExceptionType = "MISSING_INSURANCE"
	InsExInvalidInsurance        InsuranceExceptionType = "INVALID_INSURANCE"
)

//RefillDeniedType represents different refill denied exceptions
type RefillDeniedType string

const (
	RefillDeniedNeedOfficeVisit    RefillDeniedType = "NEED_OFFICE_VISIT"
	RefillDeniedNoLongerTakingMeds RefillDeniedType = "NO_LONGER_TAKING_MEDS"
	RefillDeniedPatientNotOnFile   RefillDeniedType = "PATIENT_NOT_ON_FILE"
	RefillDeniedNoResponseFromMD   RefillDeniedType = "NO_RESPONSE_FROM_MD"
	RefillDeniedOther              RefillDeniedType = "OTHER"
)

//Order represents the single prescription order can be of fillType 'FirstFill' or 'Refill'
type Order struct {
	BaseData                                `bson:",inline"`
	OrderNumber                             string               `json:"orderNumber" bson:"orderNumber" pson:"order_number" binding:"required"` //Indexed
	RxID                                    string               `json:"rxId" bson:"rxId" pson:"rx_id" binding:"required"`                      //Indexed sparse
	FillType                                string               `json:"fillType" bson:"fillType" pson:"fill_type" binding:"required"`
	PaymentOption                           string               `json:"paymentOption" bson:"paymentOption,omitempty" pson:"payment_option"`
	PaperPrescriptionImageID                string               `json:"paperPrescriptionImageId" bson:"paperPrescriptionImageId,omitempty"`
	DrugType                                string               `json:"drugType" bson:"drugType,omitempty"`
	DeliveryOption                          string               `json:"deliveryOption" bson:"deliveryOption,omitempty"`
	TransferAll                             bool                 `json:"transferAll" bson:"transferAll"`
	RerunInsurance                          bool                 `json:"rerunInsurance" bson:"rerunInsurance"`
	Archived                                bool                 `json:"archived" bson:"archived"`
	HighAlert                               bool                 `json:"highAlert" bson:"highAlert"`
	ExpeditedShipping                       bool                 `json:"expeditedShipping" bson:"expeditedShipping"`
	ShouldContactCustomerForPaymentApproval bool                 `json:"shouldContactCustomerForPaymentApproval" bson:"shouldContactCustomerForPaymentApproval"`
	SingleStepMDFunnelSignup                bool                 `json:"singleStepMDFunnelSignup" bson:"singleStepMDFunnelSignup"`
	CallCustomerForUpsell                   bool                 `json:"callCustomerForUpsell" bson:"callCustomerForUpsell"`
	BundleFill                              bool                 `json:"bundleFill" bson:"bundleFill"`
	FaxQueueID                              string               `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	SleepStatusList                         []*OrderSleepStatus  `json:"sleepStatusList" bson:"sleepStatusList,omitempty"`
	FaxSentCount                            int                  `json:"faxSentCount" bson:"faxSentCount"`
	CalledMDCount                           int                  `json:"calledMDCount" bson:"calledMDCount"`
	EScriptSentCount                        int                  `json:"eScriptSentCount" bson:"eScriptSentCount"`
	MissingInfoState                        *MissingInfoState    `json:"missingInfoState" bson:"missingInfoState,omitempty" pson:"missing_info_state"`
	TransferMilestone                       *TransferMilestone   `json:"transferMilestone" bson:"transferMilestone,omitempty" pson:"transfer_milestone"`
	NewRxMilestone                          *NewRxMilestone      `json:"newRxMilestone" bson:"orderVerification,omitempty" pson:"new_rx_milestone"` //NOTE: bson field is different than JSON!
	InsuranceMilestone                      *InsuranceMilestone  `json:"insuranceMilestone" bson:"insuranceVerificationMilestone,omitempty" pson:"insurance_milestone"`
	StockCheckMilestone                     *StockCheckMilestone `json:"stockCheckMilestone" bson:"stockVerificationMilestone,omitempty"`
	PaymentMilestone                        *PaymentMilestone    `json:"paymentMilestone" bson:"paymentMilestone,omitempty" pson:"payment_milestone"`
	DeliveryMilestone                       *DeliveryMilestone   `json:"deliveryMilestone" bson:"deliveryMilestone,omitempty" pson:"delivery_milestone"`
	RefillMilestone                         *RefillMilestone     `json:"refillMilestone" bson:"refillAuthorization,omitempty" pson:"refill_milestone"`
	NotificationStatsList                   []*NotificationStats `json:"notificationStatsList" bson:"notificationStatsList,omitempty"`
	OrderPlacedDate                         *time.Time           `json:"orderPlacedDate" bson:"orderPlacedDate,omitempty"`
	OriginPharmacyRef                       *mgo.DBRef           `json:"-" bson:"originPharmacy,omitempty" pson:"origin_pharmacy"`
	PartnerPharmacyRef                      *mgo.DBRef           `json:"-" bson:"partnerPharmacy,omitempty" pson:"partner_pharmacy"`
	AddressRef                              *mgo.DBRef           `json:"-" bson:"deliveryAddress,omitempty" pson:"delivery_address"`
	PaymentRef                              *mgo.DBRef           `json:"-" bson:"paymentInfo,omitempty" pson:"payment_info"`
	OriginPharmacy                          *Pharmacy            `json:"originPharmacy" bson:"-"`
	PartnerPharmacy                         *Pharmacy            `json:"partnerPharmacy" bson:"-"`
	Address                                 *Address             `json:"address" bson:"-"`
	Payment                                 *Payment             `json:"payment" bson:"-"`
}

// CollectionName function from gmgo.Document interface
func (o Order) CollectionName() string {
	return "order"
}

//FaxContent data representation for fax content
type FaxContent struct {
	BaseData    `bson:",inline"`
	OrderNumber string `json:"orderNumber" bson:"orderNumber,omitempty"`
	Content     string `json:"content" bson:"content,omitempty"`
	FaxType     string `json:"faxType" bson:"faxType,omitempty"`
}

//CollectionName function from gmgo.Document interface to return collection name in MongoDB
func (fc FaxContent) CollectionName() string {
	return "faxContent"
}

// IsFirstFill check if order is first fill
func (o Order) IsFirstFill() bool {
	return o.FillType == "FirstFill"
}

// IsReFill check if order is first fill
func (o Order) IsReFill() bool {
	return o.FillType == "ReFill"
}

//IsProcessingTransfer returns true if transfer is in process
func (o *Order) IsProcessingTransfer() bool {
	if o.TransferStatus() != "" && o.TransferStatus() != MSCompleted && !o.IsTransferBack() {
		return true
	}
	return false
}

//IsProcessingNewPrescription returns true if new prescription processing is in process
func (o *Order) IsProcessingNewPrescription() bool {
	if o.NewPrescriptionStatus() != "" && o.NewPrescriptionStatus() != MSCompleted {
		return true
	}
	return false
}

//IsProcessingInsurance returns true if rx is in run insurance step
func (o *Order) IsProcessingInsurance() bool {
	if o.IsPayOutOfPocket() {
		return false
	}
	if o.InsuranceStatus() != "" && o.InsuranceStatus() != MSCompleted {
		return true
	}
	return false
}

//IsProcessingInsuranceComplete returns true if rx is in run insurance is complete
func (o *Order) IsProcessingInsuranceComplete() bool {
	if o.IsPayOutOfPocket() {
		return false
	}
	if o.InsuranceStatus() != "" && o.InsuranceStatus() == MSCompleted {
		return true
	}
	return false
}

//IsProcessingCash returns true if order's cash price has not been determined
func (o *Order) IsProcessingCash() bool {
	if !o.IsPayOutOfPocket() {
		return false
	}
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSNotInitiated {
		return true
	}
	return false
}

//IsNotDue returns true if rx is marked as not due for now
func (o *Order) IsNotDue() bool {
	if o.IsPayOutOfPocket() {
		return false
	}
	if o.InsuranceStatus() != "" && o.InsuranceStatus() == MSRefillNotDue {
		return true
	}
	return false
}

// IsProcessingPayment returns true if rx is in payment processing
func (o *Order) IsProcessingPayment() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() != MSCompleted {
		return true
	}
	return false
}

// IsPendingPaymentApproval returns true if order is in payment approval state
func (o *Order) IsPendingPaymentApproval() bool {
	if o.PaymentStatus() != "" &&
		(o.PaymentStatus() == MSPendingApproval ||
			(o.PaymentStatus() == MSPendingApprovalAfterPriceInspection &&
				!o.ShouldContactCustomerForPaymentApproval)) {
		return true
	}
	return false
}

// IsPendingPaymentApproval returns true if order is in payment approval state
func (o *Order) IsPendingPaymentApprovalAndSignup() bool {
	return o.PaymentStatus() == MSPendingApprovalAndSignup
}

//IsCashPriceIdentified returns true if order's cash price has been determined
func (o *Order) IsCashPriceIdentified() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSCashPriceIdentified {
		return true
	}
	return false
}

//IsPendingApprovalUntilFreeTrial returns true if order is in free trial and is pending approval
func (o *Order) IsPendingApprovalUntilFreeTrial() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSPendingApprovalUntilFreeTrial {
		return true
	}
	return false
}

//IsPendingApprovalUntilPriceInspectionFreeTrial returns true if order is in free trial and is pending approval
func (o *Order) IsPendingApprovalUntilPriceInspectionFreeTrial() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSPendingApprovalUntilPriceInspectionFreeTrial {
		return true
	}
	return false
}

//IsPendingApprovalUntilPriceInspectionNoFreeTrial returns true if order is in free trial and is pending approval
func (o *Order) IsPendingApprovalUntilPriceInspectionNoFreeTrial() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSPendingApprovalUntilPriceInspectionNoFreeTrial {
		return true
	}
	return false
}

// IsPaymentAuthorized returns true if order is in payment is authorized
func (o *Order) IsPaymentAuthorized() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSPaymentAuthorized {
		return true
	}
	return false
}

//IsPaymentComplete returns true if payment is completed, meaning we already charged the customer
func (o *Order) IsPaymentComplete() bool {
	if o.PaymentStatus() != "" && o.PaymentStatus() == MSCompleted {
		return true
	}
	return false
}

// IsProcessingStockCheck returns true if it's processing stock
func (o *Order) IsProcessingStockCheck() bool {
	if o.StockCheckStatus() != "" && o.StockCheckStatus() != MSCompleted {
		return true
	}
	return false
}

// IsProcessingDelivery returns true if it's processing delivery
func (o *Order) IsProcessingDelivery() bool {
	if o.DeliveryStatus() != "" && o.DeliveryStatus() != MSShipped {
		return true
	}
	return false
}

//IsShipped returns true if order is shipped
func (o *Order) IsShipped() bool {
	if o.DeliveryStatus() != "" && (o.DeliveryStatus() == MSShipped || o.DeliveryStatus() == MSWaitingLabelScan) {
		return true
	}
	return false
}

//IsShippingLabelGenerated returns true if shipping label has been generated
func (o *Order) IsShippingLabelGenerated() bool {
	if o.IsShipped() {
		return true
	}
	return o.DeliveryStatus() == MSShippingLabelGenerated
}

//IsDelivered returns true if delivery is confirmed
func (o *Order) IsDelivered() bool {
	return o.DeliveryMilestone.DeliveryConfirmed
}

//HasDeliveryFailed returns true if delivery failed
func (o *Order) HasDeliveryFailed() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == DeliveryFailed
}

//HasDeliveryReturned returns true if delivery was returned
func (o *Order) HasDeliveryReturned() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == DeliveryReturned
}

//IsProcessingRefillAuth returns true if refill auth is being processes
func (o *Order) IsProcessingRefillAuth() bool {
	if o.RefillAuthStatus() != "" && o.RefillAuthStatus() != MSCompleted && o.RefillAuthStatus() != MSNotInitiated {
		return true
	}

	return false
}

//IsRefillsDenied returns true if order is marked as refill denied by MD
func (o *Order) IsRefillsDenied() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == RefillsDenied
}

//IsPriorAuthDenied returns true if order is marked as prior auth is denied by MD
func (o *Order) IsPriorAuthDenied() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == PriorAuthDenied
}

//IsStepTherapyDenied returns true if order is marked as step therapy denied by MD
func (o *Order) IsStepTherapyDenied() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == StepTherapyDenied
}

//RefillDeniedType returns type of refill denied
func (o *Order) RefillDeniedType() *RefillDeniedType {
	if !o.IsRefillsDenied() {
		return nil
	}
	return o.MissingInfoState.RefillDeniedType
}

// HasPaymentError check if the order has payment error
func (o *Order) HasPaymentError() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == PaymentError
}

// TransferStatus returns transfer status of the order
func (o *Order) TransferStatus() string {
	if o.TransferMilestone == nil {
		return ""
	}
	return o.TransferMilestone.Status
}

//IsTransferBack returns if transfer back is initiated
func (o *Order) IsTransferBack() bool {
	return o.TransferStatus() == MSTransferBack || o.TransferStatus() == MSTransferBackCompleted
}

//IsTransferBackComplete returns true if transfer back process is completed
func (o *Order) IsTransferBackComplete() bool {
	return o.TransferStatus() == MSTransferBackCompleted
}

//IsTransferOut returns if transfer out is initiated
func (o *Order) IsTransferOut() bool {
	return o.NewPrescriptionStatus() == MSTransferOut || o.NewPrescriptionStatus() == MSTransferOutCompleted
}

//IsTransferBackComplete returns true if transfer back process is completed
func (o *Order) IsTransferOutComplete() bool {
	return o.NewPrescriptionStatus() == MSTransferOutCompleted
}

//IsTransferRequestSent returns if transfer request was sent
func (o *Order) IsTransferRequestSent() bool {
	return o.TransferStatus() == MSTransferRequested || o.TransferStatus() == MSFaxSent
}

//IsTransferInitiated returns if transfer is initiated
func (o *Order) IsTransferInitiated() bool {
	return o.TransferStatus() == MSTransferInitiated
}

//IsOrderVerificationInitiated returns if new rx processing is initiated
func (o *Order) IsOrderVerificationInitiated() bool {
	return o.NewPrescriptionStatus() == MSVerificationInitiated
}

//NewPrescriptionStatus returns new rx processing status
func (o *Order) NewPrescriptionStatus() string {
	if o.NewRxMilestone == nil {
		return ""
	}
	return o.NewRxMilestone.Status
}

// PaymentStatus returns payment milestone current status
func (o *Order) PaymentStatus() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	return o.PaymentMilestone.Status
}

// InsuranceStatus insurance processing milestone current status
func (o *Order) InsuranceStatus() string {
	if o.InsuranceMilestone == nil {
		return ""
	}
	return o.InsuranceMilestone.Status
}

//StockCheckStatus returns refill auth status
func (o *Order) StockCheckStatus() string {
	if o.StockCheckMilestone == nil {
		return ""
	}
	return o.StockCheckMilestone.Status
}

// DeliveryStatus returns delivery current status
func (o *Order) DeliveryStatus() string {
	if o.DeliveryMilestone == nil {
		return ""
	}
	return o.DeliveryMilestone.Status
}

// MissingStatus returns missing info status
func (o *Order) MissingStatus() *MissingStatus {
	mis := o.MissingInfoState
	if mis == nil {
		return nil
	}
	return mis.MissingStatus
}

//RefillAuthStatus returns refill auth status
func (o *Order) RefillAuthStatus() string {
	if o.RefillMilestone == nil {
		return ""
	}
	return o.RefillMilestone.Status
}

//MedicationPrice - gets medication price
func (o *Order) MedicationPrice() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	return o.PaymentMilestone.MedicationPrice
}

// CarrierType - get carrier type
func (o Order) CarrierType() string {
	if o.DeliveryMilestone == nil {
		return ""
	}
	return o.DeliveryMilestone.CarrierType
}

//DeliveryCharge - get delivery charge
func (o *Order) DeliveryCharge() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	return o.PaymentMilestone.DeliveryCharge
}

//FinalCharge - get final charge
func (o *Order) FinalCharge() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	return o.PaymentMilestone.FinalCharge
}

//PromoDiscount - get discount
func (o *Order) PromoDiscount() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	return o.PaymentMilestone.AppliedPromoDiscount
}

//IsPromoDiscountApplied returns true if promo discount is applied
func (o *Order) IsPromoDiscountApplied() bool {
	return o.PaymentMilestone.AppliedPromoDiscount != "" &&
		o.PaymentMilestone.AppliedPromoDiscount != "0.00" &&
		o.PaymentMilestone.AppliedPromoDiscount != "0.0" &&
		o.PaymentMilestone.AppliedPromoDiscount != "0"
}

//EarliestRefillDate - get earliest refill date
func (o *Order) EarliestRefillDate() *time.Time {
	if o.InsuranceMilestone == nil {
		return nil
	}

	return o.InsuranceMilestone.EarliestRefillDate
}

// IsPayOutOfPocket returns true if patient want to pay out of pocket for the med
func (o *Order) IsPayOutOfPocket() bool {
	return o.PaymentOption == "PayDirectly"
}

//IsPayThroughInsurance returns true if patient is paying with insurance
func (o *Order) IsPayThroughInsurance() bool {
	return o.PaymentOption == "PayThroughInsurance"
}

//PaymentAuthToken returns payment authorization token
func (o *Order) PaymentAuthToken() string {
	var token string
	if o.IsPayThroughInsurance() {
		token = o.paymentTokenForPayWithInsurance()
		if token == "" {
			token = o.paymentTokenForPayCash()
		}
	} else {
		token = o.paymentTokenForPayCash()
		if token == "" {
			token = o.paymentTokenForPayWithInsurance()
		}
	}
	return token
}

func (o *Order) paymentTokenForPayWithInsurance() string {
	if o.InsuranceMilestone == nil {
		return ""
	}
	if token, ok := o.InsuranceMilestone.Parameters["paymentAuthorizationToken"]; ok {
		return token.(string)
	}

	return ""
}

func (o *Order) paymentTokenForPayCash() string {
	if o.PaymentMilestone == nil {
		return ""
	}
	if token, ok := o.PaymentMilestone.Parameters["paymentAuthorizationToken"]; ok {
		return token.(string)
	}
	return ""
}

//UpdateNotificatonStats update the notification stats based on type
func (o *Order) UpdateNotificatonStats(notificationType string) {
	ns := new(NotificationStats)
	ns.Type = notificationType

	t := time.Now()
	ns.SentDate = &t

	notificationStatsList := o.NotificationStatsList
	if notificationStatsList == nil {
		notificationStatsList = make([]*NotificationStats, 0)
		notificationStatsList = append(notificationStatsList, ns)
	} else {
		notificationStatsList = append(notificationStatsList, ns)
	}

	o.NotificationStatsList = notificationStatsList
}

// UpdateSleepState updates the order sleep state
func (o *Order) UpdateSleepState(sleepState OrderSleepState) {
	t := time.Now().UTC()
	o.UpdateSleepStateWithStartDate(sleepState, t)
}

//UpdateTimedSleepState add timed sleep state to order with wake up time in future
func (o *Order) UpdateTimedSleepState(wakeupAt *time.Time) {
	t := time.Now().UTC()
	o.UpdateSleepStateWithStartAndStopDate(TimedSleep, &t, wakeupAt)
}

// UpdateSleepStateWithStartDate updates the order sleep state with start date
func (o *Order) UpdateSleepStateWithStartDate(sleepState OrderSleepState, sleepStartDate time.Time) {
	var sleepStopDate *time.Time
	if sleepState == TimedSleep {
		t := time.Now().UTC()
		nextDay := t.Add(24 * time.Hour)
		sleepStopDate = &nextDay
	}
	o.UpdateSleepStateWithStartAndStopDate(sleepState, &sleepStartDate, sleepStopDate)
}

//UpdateSleepStateWithStartAndStopDate add sleep state to order with start and stop date
func (o *Order) UpdateSleepStateWithStartAndStopDate(sleepState OrderSleepState, sleepStartDate *time.Time, sleepStopDate *time.Time) {
	if o.SleepStatusList == nil {
		o.SleepStatusList = make([]*OrderSleepStatus, 0)

		status := newSleepStateStatus(sleepState, sleepStartDate, sleepStopDate)
		o.SleepStatusList = append(o.SleepStatusList, status)
	} else {
		found := false
		for _, status := range o.SleepStatusList {
			if *status.State == sleepState {
				found = true
				if sleepState != TimedSleep {
					status.ResetNeeded = true
				}

				status.SleepStartDate = sleepStartDate
				status.SleepStopDate = sleepStopDate
			}
		}
		if !found {
			status := newSleepStateStatus(sleepState, sleepStartDate, sleepStopDate)
			o.SleepStatusList = append(o.SleepStatusList, status)
		}
	}
}

// GetSleepStatus get sleep status based on state
func (o *Order) GetSleepStatus(sleepState OrderSleepState) *OrderSleepStatus {
	if o.SleepStatusList == nil {
		return nil
	}
	for _, status := range o.SleepStatusList {
		if *status.State == sleepState {
			return status
		}
	}
	return nil
}

// RemoveSleepStatus remove sleep status
func (o *Order) RemoveSleepStatus(sleepState OrderSleepState) {
	if o.SleepStatusList == nil {
		return
	}

	newList := make([]*OrderSleepStatus, 0)
	for _, status := range o.SleepStatusList {
		if *status.State != sleepState {
			newList = append(newList, status)
		}
	}
	o.SleepStatusList = newList
}

//MoveToScrubAndRoute moves the order to scrub ready status
func (o *Order) MoveToScrubAndRoute() {
	t := time.Now().UTC()

	if o.TransferMilestone != nil && o.TransferStatus() == MSNotInitiated {
		o.TransferMilestone.Status = MSCallCompleted
		o.TransferMilestone.LastProcessedDate = &t
	}
	if o.NewRxMilestone != nil && o.NewPrescriptionStatus() == MSNotInitiated {
		o.NewRxMilestone.Status = MSCallCompleted
		o.NewRxMilestone.LastProcessedDate = &t
	}
}

//PendingActionTaken mark the order as pending action taken
func (o *Order) PendingActionTaken() {
	t := time.Now().UTC()

	if o.TransferMilestone != nil {
		o.TransferMilestone.Status = MSPendingActionTaken
		o.TransferMilestone.LastProcessedDate = &t
	}
	if o.NewRxMilestone != nil {
		o.NewRxMilestone.Status = MSPendingActionTaken
		o.NewRxMilestone.LastProcessedDate = &t
	}

	o.MissingInfoState = nil
	o.RemoveSleepStatus(CallPatient)
}

//SavePendingPriorStatus save the transfer/newRx status prior to pending
func (o *Order) SavePendingPriorStatus() {
	if o.TransferMilestone != nil && o.TransferMilestone.Status != MSPendingActionTaken {
		o.TransferMilestone.PendingPriorStatus = o.TransferMilestone.Status
	}
	if o.NewRxMilestone != nil && o.NewRxMilestone.Status != MSPendingActionTaken {
		o.NewRxMilestone.PendingPriorStatus = o.NewRxMilestone.Status
	}
}

//InsuranceOverride mark order as insurance override. Used when user requests
//ship sooner
func (o *Order) InsuranceOverride() {
	if o.InsuranceMilestone != nil {
		o.InsuranceMilestone.Status = MSInsuranceOverride
		o.InsuranceMilestone.SendToPP = false
	}
}

//IsInvoicePending return if invoice is still to be sent for order
func (o *Order) IsInvoicePending() bool {
	if o.PaymentMilestone != nil {
		status := o.PaymentMilestone.Status
		return status == MSPaymentCompleteInvoicePending || status == MSPaymentIncompleteInvoicePending
	}
	return false
}

//IsPaymentIncompleteOffline return true if offline payment is still incomplete
func (o *Order) IsPaymentIncompleteOffline() bool {
	if o.PaymentMilestone != nil && o.PaymentStatus() == MSPaymentIncompleteOffline {
		return true
	}
	return false
}

//IsInTransit returns true if the order delivery status is in transit. Delivery
//status is based on shippo callback
func (o *Order) IsInTransit() bool {
	if o.DeliveryMilestone != nil {
		return o.DeliveryMilestone.DeliveryStatus == MSInTransit
	}
	return false
}

// OrderSleepStatus order sleep status
type OrderSleepStatus struct {
	State           *OrderSleepState `json:"state" bson:"state"`
	SleepStartDate  *time.Time       `json:"sleepStartDate" bson:"sleepStartDate,omitempty"`
	SleepStopDate   *time.Time       `json:"sleepStopDate" bson:"sleepStopDate,omitempty"`
	ResetNeeded     bool             `json:"resetNeeded" bson:"resetNeeded"`
	SleepResetCount int              `json:"sleepResetCount" bson:"sleepResetCount"`
}

// MissingInfoState represents missing info state with missing info data
type MissingInfoState struct {
	MissingFields          []*MissingField         `json:"missingFields" bson:"missingFields,omitempty"`
	MissingStatus          *MissingStatus          `json:"missingStatus" bson:"missingStatus,omitempty"`
	Comment                string                  `json:"comment" bson:"comment,omitempty"`
	InsuranceExceptionType *InsuranceExceptionType `json:"insuranceExceptionType" bson:"insuranceExceptionType,omitempty"`
	RefillDeniedType       *RefillDeniedType       `json:"refillDeniedType" bson:"refillDeniedType,omitempty"`
}

// NotificationStats keeps track of order notification stats
type NotificationStats struct {
	Type     string     `json:"type" bson:"type,omitempty"`
	SentDate *time.Time `json:"sentDate" bson:"sentDate,omitempty"`
}

// OrderByRef returns order data based on given reference
func OrderByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Order, error) {
	order := new(Order)
	if err := session.FindByRef(dbRef, order); err != nil {
		return nil, err
	}
	return order, nil
}

// AddressByRef returns address data based on given reference
func AddressByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Address, error) {
	addr := new(Address)
	if err := session.FindByRef(dbRef, addr); err != nil {
		return addr, err
	}

	return addr, nil
}

// PaymentByRef returns payment data based on given patient reference
func PaymentByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Payment, error) {
	p := new(Payment)
	if err := session.FindByRef(dbRef, p); err != nil {
		return nil, err
	}

	return p, nil
}

// PharmacyByRef returns payment data based on given patient reference
func PharmacyByRef(dbRef *mgo.DBRef, session *gmgo.DbSession) (*Pharmacy, error) {
	p := new(Pharmacy)
	if err := session.FindByRef(dbRef, p); err != nil {
		return nil, err
	}

	return p, nil
}

//IsPendingTransfer - Rx is new transfer - passed welcome call and has not been transferred yet
func (o Order) IsPendingTransfer() bool {
	return o.TransferStatus() == MSNotInitiated || o.TransferStatus() == MSCallCompleted
}

//IsTransferInProcess - returns true if transfer is in progress and not completed yet
func (o Order) IsTransferInProcess() bool {
	return o.TransferStatus() != MSCompleted && o.TransferStatus() != MSCallCompleted && o.TransferStatus() != MSNotInitiated
}

//RequireDoctorApproval - checks whether the order requires MD approval
func (o Order) RequireDoctorApproval() bool {
	return o.TransferStatus() == MSZeroRefills || o.TransferStatus() == MSEScriptRequestSent || o.TransferStatus() == MSFaxSent
}

//RequireDoctorApprovalForNewPrescription checks whether the new prescription is waiting for MD approval
func (o Order) RequireDoctorApprovalForNewPrescription() bool {
	return o.NewPrescriptionStatus() == MSFaxSent || o.NewPrescriptionStatus() == MSEScriptRequestSent
}

//IsInvalidInsurance check if insurance is invalid
func (o Order) IsInvalidInsurance() bool {
	if o.MissingInfoState != nil && o.MissingInfoState.MissingStatus != nil {
		status := o.MissingInfoState.MissingStatus
		return *status == InvalidInsurance || *status == InsuranceExpired
	}
	return false
}

//IsStockException returns true if there is a stock exception
func (o *Order) IsStockException() bool {
	return o.MissingStatus() != nil && *o.MissingStatus() == StockException
}

// IsInReverseInsuranceProcess check if insurance is in process of being reversed
func (o Order) IsInReverseInsuranceProcess() bool {
	return o.InsuranceMilestone != nil && o.InsuranceMilestone.Status == MSReverseInsurance
}

// IsInsuranceReversed check if insurance has already been reversed
func (o Order) IsInsuranceReversed() bool {
	return o.InsuranceMilestone != nil && o.InsuranceMilestone.Status == MSInsuranceReversed
}

//RequirePriorAuth checks if order requires prior auth
func (o Order) RequirePriorAuth() bool {
	return o.InsuranceStatus() == MSPriorAuth || o.InsuranceStatus() == MSFaxSent
}

//InitializeOrder initialize new order data based on given number and origin type
func InitializeOrder(orderNumber string, originType string) *Order {
	order := new(Order)
	order.InitData()
	order.OrderNumber = orderNumber
	if originType == "Transfer" {
		order.InitializeTransferMilestone()
		order.TransferAll = true
	} else {
		order.InitializeNewRxMilestone()
	}
	order.FillType = FirstFill
	order.InitializeInsuranceMilestone()
	order.InitializeDeliveryMilestone()
	order.InitializePaymentMilestone()
	order.InitializeStockCheckMilestone()
	return order
}

//InitializeTransferMilestone initialize transfer milestone
func (o *Order) InitializeTransferMilestone() {
	tm := new(TransferMilestone)
	tm.Status = MSCallCompleted
	tm.DaysOfSupply = 30
	o.TransferMilestone = tm
}

//InitializeNewRxMilestone setup new prescription milestone
func (o *Order) InitializeNewRxMilestone() {
	nm := new(NewRxMilestone)
	nm.Status = MSCallCompleted
	nm.MedicationQuantity = 30
	o.NewRxMilestone = nm
}

//InitializeInsuranceMilestone setup insurance milestone
func (o *Order) InitializeInsuranceMilestone() {
	im := new(InsuranceMilestone)
	im.Status = MSNotInitiated
	im.SendCopayConfirmEmail = true
	o.InsuranceMilestone = im
}

//InitializePaymentMilestone setup payment milestone
func (o *Order) InitializePaymentMilestone() {
	pm := new(PaymentMilestone)
	pm.Status = MSNotInitiated
	pm.AppliedDummyCharge = false
	o.PaymentMilestone = pm
}

//InitializeStockCheckMilestone setup stock check milestone
func (o *Order) InitializeStockCheckMilestone() {
	sm := new(StockCheckMilestone)
	sm.Status = MSNotInitiated
	o.StockCheckMilestone = sm
}

//InitializeDeliveryMilestone setup delivery milestone
func (o *Order) InitializeDeliveryMilestone() {
	dm := new(DeliveryMilestone)
	dm.Status = MSNotInitiated
	dm.DeliveryConfirmed = false
	o.DeliveryMilestone = dm
}

//AddOriginPharmacy add OP info to order
func (o *Order) AddOriginPharmacy(pharmacy *Pharmacy) {
	if pharmacy != nil {
		o.OriginPharmacy = pharmacy
		o.OriginPharmacyRef = DBRef(o.OriginPharmacy.CollectionName(), o.OriginPharmacy.ID)
	}
}

//AddPartnerPharmacy add PP info to order
func (o *Order) AddPartnerPharmacy(pharmacy *Pharmacy) {
	if pharmacy != nil {
		o.PartnerPharmacy = pharmacy
		o.PartnerPharmacyRef = DBRef(o.PartnerPharmacy.CollectionName(), o.PartnerPharmacy.ID)
	}
}

//AddDeliveryAddress add delivery address to order
func (o *Order) AddDeliveryAddress(address *Address) {
	if address != nil {
		o.Address = address
		o.AddressRef = DBRef(o.Address.CollectionName(), o.Address.ID)
	}
}

//AddPayment add payment info to order
func (o *Order) AddPayment(payment *Payment) {
	if payment != nil {
		o.Payment = payment
		o.PaymentRef = DBRef(o.Payment.CollectionName(), o.Payment.ID)
	}
}

//SetMissingInfoState mark order has missing info or error state
func (o *Order) SetMissingInfoState(status MissingStatus) {
	missingInfo := new(MissingInfoState)
	missingInfo.MissingStatus = &status
	o.MissingInfoState = missingInfo
}

func newSleepStateStatus(sleepState OrderSleepState, sleepStartDate *time.Time, sleepStopDate *time.Time) *OrderSleepStatus {
	status := new(OrderSleepStatus)
	status.State = &sleepState
	if sleepState != TimedSleep {
		status.ResetNeeded = true
	}

	status.SleepStartDate = sleepStartDate
	status.SleepStopDate = sleepStopDate

	return status
}

//InitializeMilestones - initializes processing milestones
func (o *Order) InitializeMilestones() {
	o.InitializeInsuranceMilestone()
	o.InitializeDeliveryMilestone()
	o.InitializePaymentMilestone()
	o.InitializeStockCheckMilestone()
}
