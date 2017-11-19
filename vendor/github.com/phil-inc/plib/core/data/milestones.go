package data

import "time"

//represents different milestone status
const (
	//NotInitiated means milestone hasn't been initiated or processed yet.
	MSNotInitiated                                   = "NOT_INITIATED"
	MSCallCompleted                                  = "CALL_COMPLETED"
	MSCompleted                                      = "COMPLETED"
	MSTransferInitiated                              = "TRANSFER_INITIATED"
	MSTransferRequested                              = "TRANSFER_REQUEST_SENT"
	MSZeroRefills                                    = "ZERO_REFILLS"
	MSEScriptRequestSent                             = "ESCRIPT_REQUEST_SENT"
	MSTransferBack                                   = "TRANSFER_BACK"
	MSTransferOut                                    = "TRANSFER_OUT"
	MSVerificationInitiated                          = "VERIFICATION_INITIATED"
	MSTransferBackCompleted                          = "TRANSFER_BACK_COMPLETED"
	MSTransferOutCompleted                           = "TRANSFER_OUT_COMPLETED"
	MSRefillNotDue                                   = "REFILL_NOT_DUE"
	MSSendNewRxFax                                   = "SEND_NEW_RX_FAX"
	MSTransferControlled                             = "TRANSFER_CONTROLLED"
	MSRouteControlled                                = "ROUTE_CONTROLLED"
	MSReverseInsurance                               = "REVERSE_INSURANCE"
	MSInsuranceReversed                              = "INSURANCE_REVERSED"
	MSPendingApproval                                = "PENDING_APPROVAL"
	MSPaymentAuthorized                              = "PAYMENT_AUTHORIZED"
	MSFaxSent                                        = "FAX_SENT"
	MSPriorAuth                                      = "PRIOR_AUTH"
	MSPriorAuthProcessing                            = "PRIOR_AUTH_PROCESSING"
	MSShipped                                        = "SHIPPED"
	MSWaitingLabelScan                               = "WAITING_LABEL_SCAN"
	MSShippingLabelGenerated                         = "SHIPPING_LABEL_GENERATED"
	MSGenerateShippingLabel                          = "GENERATE_SHIPPING_LABEL"
	MSPendingActionTaken                             = "PENDING_ACTION_TAKEN"
	MSCashPriceCheck                                 = "CASH_PRICE_CHECK"
	MSCashPriceIdentified                            = "CASH_PRICE_IDENTIFIED"
	MSPendingApprovalUntilFreeTrial                  = "PENDING_APPROVAL_UNTIL_FREE_TRIAL"
	MSPendingApprovalUntilPriceInspection            = "PENDING_APPROVAL_UNTIL_PRICE_INSPECTION"
	MSPendingApprovalUntilPriceInspectionByPP        = "PENDING_APPROVAL_UNTIL_PRICE_INSPECTION_BY_PP"
	MSPendingApprovalUntilPriceInspectionNoFreeTrial = "PENDING_APPROVAL_UNTIL_PRICE_INSPECTION_NO_FREE_TRIAL"
	MSPendingApprovalUntilPriceInspectionFreeTrial   = "PENDING_APPROVAL_UNTIL_PRICE_INSPECTION_FREE_TRIAL"
	MSPendingApprovalAfterPriceInspection            = "PENDING_APPROVAL_AFTER_PRICE_INSPECTION"
	MSPendingApprovalAndSignup                       = "PENDING_APPROVAL_AND_SIGNUP"
	MSPendingPaymentOffline                          = "PENDING_APPROVAL_OFFLINE"
	MSPaymentCompleteInvoicePending                  = "PAYMENT_COMPLETE_INVOICE_PENDING"
	MSPaymentIncompleteInvoicePending                = "PAYMENT_INCOMPLETE_INVOICE_PENDING"
	MSPaymentIncompleteOffline                       = "PAYMENT_INCOMPLETE_OFFLINE"
	MSInsuranceOverride                              = "INSURANCE_OVERRIDE"
	MSInTransit                                      = "IN_TRANSIT"
	MSDelivered                                      = "DELIVERED"
)

//MilestoneNames list out the names of all the milestones
var MilestoneNames = []string{"transferMilestone", "orderVerification", "insuranceVerificationMilestone", "paymentMilestone", "deliveryMilestone", "stockVerificationMilestone", "refillAuthorization"}

//MilestoneStatusList list all the possible values for milestone status
var MilestoneStatusList = []string{MSNotInitiated, MSCompleted, MSCallCompleted, MSTransferInitiated, MSTransferBack, MSTransferRequested, MSZeroRefills, MSEScriptRequestSent, MSFaxSent, MSPendingApproval, MSPaymentAuthorized, MSPriorAuth, MSRefillNotDue, MSShipped, MSWaitingLabelScan, MSReverseInsurance, MSInsuranceReversed}

const (
	//Yes represents enum yes
	Yes string = "Yes"
	//No represents enum No
	No string = "No"
)

// DrugType represents drug type
type DrugType string

const (
	// Generic type
	Generic DrugType = "Generic"
	//BrandName brand type
	BrandName DrugType = "BrandName"
)

//PaymentProvider for payment processing
type PaymentProvider string

const (
	//Stripe processor
	Stripe PaymentProvider = "Stripe"
	//Braintree processor
	Braintree PaymentProvider = "Braintree"
)

// Milestone represents the common milestone data for each milestone.
type Milestone struct {
	Status            string                 `json:"status" bson:"status,omitempty" pson:"status"` //sparse indexed
	Comment           string                 `json:"comment" bson:"comment,omitempty"`
	Parameters        map[string]interface{} `json:"parameters" bson:"parameters,omitempty"`
	LastProcessedDate *time.Time             `json:"lastProcessedDate" bson:"lastProcessedDate,omitempty"`
}

// TransferMilestone represents the transfer milestone data
type TransferMilestone struct {
	Milestone                  `bson:",inline"`
	MedicationName             string     `json:"medicationName" bson:"medicationName,omitempty"`
	MedicationStrength         string     `json:"medicationStrength" bson:"medicationStrength,omitempty"`
	MedicationForm             string     `json:"medicationForm" bson:"medicationForm,omitempty"`
	Sig                        string     `json:"sig" bson:"sig,omitempty"`
	RefillAuthorizedBy         string     `json:"refillAuthorizedBy" bson:"refillAuthorizedBy,omitempty"`
	MedicationQuantity         int        `json:"medicationQuantity" bson:"medicationQuantity"`
	DaysOfSupply               int        `json:"daysOfSupply" bson:"daysOfSupply"`
	NumberOfRefillsTransferred int        `json:"numberOfRefillsTransferred" bson:"numberOfRefillsTransferred"`
	DoctorName                 string     `json:"doctorName" bson:"doctorName,omitempty"`
	DoctorPhoneNumber          string     `json:"doctorPhoneNumber" bson:"doctorPhoneNumber,omitempty"`
	DoctorFaxNumber            string     `json:"doctorFaxNumber" bson:"doctorFaxNumber,omitempty"`
	DoctorNpi                  string     `json:"doctorNpi" bson:"doctorNpi,omitempty"`
	DoctorCity                 string     `json:"doctorCity" bson:"doctorCity,omitempty"`
	DrugType                   *DrugType  `json:"drugType" bson:"drugType,omitempty"`
	LastFillDate               *time.Time `json:"lastFillDate" bson:"lastFillDate,omitempty"`
	CurrentCopay               string     `json:"currentCopay" bson:"currentCopay,omitempty"`
	FaxQueueID                 string     `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	CouponAvailable            string     `json:"couponAvailable" bson:"couponAvailable,omitempty"`
	CopayAvailable             string     `json:"copayAvailable" bson:"copayAvailable,omitempty"`
	CouponGroupNumber          string     `json:"couponGroupNumber" bson:"couponGroupNumber,omitempty"`
	CouponBinNumber            string     `json:"couponBinNumber" bson:"couponBinNumber,omitempty"`
	CouponIDNumber             string     `json:"couponIdNumber" bson:"couponIdNumber,omitempty"`
	CouponPcnNumber            string     `json:"couponPcnNumber" bson:"couponPcnNumber,omitempty"`
	RefillRequestedDate        *time.Time `json:"refillRequestedDate" bson:"refillRequestedDate,omitempty"`
	PendingPriorStatus         string     `json:"pendingPriorStatus" bson:"pendingPriorStatus,omitempty"`
	MDApprovedOverPhone        bool       `json:"mdApprovedOverPhone" bson:"mdApprovedOverPhone"`
	OPSentFax                  bool       `json:"opSentFax" bson:"opSentFax"`
}

// NewRxMilestone new rx order milestone
type NewRxMilestone struct {
	Milestone          `bson:",inline"`
	MedicationName     string    `json:"medicationName" bson:"medicationName,omitempty"`
	MedicationQuantity int       `json:"medicationQuantity" bson:"medicationQuantity"`
	DaysOfSupply       int       `json:"daysOfSupply" bson:"daysOfSupply"`
	RefillsRemaining   int       `json:"refillsRemaining" bson:"refillsRemaining"`
	DoctorName         string    `json:"doctorName" bson:"doctorName,omitempty"`
	DoctorPhoneNumber  string    `json:"doctorPhoneNumber" bson:"doctorPhoneNumber,omitempty"`
	DoctorFaxNumber    string    `json:"doctorFaxNumber" bson:"doctorFaxNumber,omitempty"`
	DoctorNpi          string    `json:"doctorNpi" bson:"doctorNpi,omitempty"`
	DoctorCity         string    `json:"doctorCity" bson:"doctorCity,omitempty"`
	DrugType           *DrugType `json:"drugType" bson:"drugType,omitempty"`
	FaxQueueID         string    `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	PendingPriorStatus string    `json:"pendingPriorStatus" bson:"pendingPriorStatus,omitempty"`
	NeedsMDApproval    bool      `json:"needsMDApproval" bson:"needsMDApproval"`
}

// StockCheckMilestone represents stock verification step
type StockCheckMilestone struct {
	Milestone         `bson:",inline"`
	StockAvailability *time.Time `json:"stockAvailability" bson:"stockAvailability,omitempty"`
}

// InsuranceMilestone for insurance verification step
type InsuranceMilestone struct {
	Milestone                     `bson:",inline"`
	CanRefill                     string     `json:"canRefill" bson:"canRefill,omitempty"` //CanRefill can only have 'Yes' or 'No'
	NewRxNumber                   string     `json:"newRxNumber" bson:"newRxNumber,omitempty"`
	FaxQueueID                    string     `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	IsFederalSponsoredInsurance   string     `json:"isFederalSponsoredInsurance" bson:"isFederalSponsoredInsurance,omitempty"`
	NewCopay                      string     `json:"newCopay" bson:"newCopay,omitempty"`
	EarliestRefillDate            *time.Time `json:"earliestRefillDate" bson:"earliestRefillDate,omitempty" pson:"earliest_refill_date"`
	InsuranceRunDate              *time.Time `json:"insuranceRunDate" bson:"insuranceRunDate,omitempty"`
	RefillNotDueMarkedDueDate     *time.Time `json:"refillNotDueMarkedDueDate" bson:"refillNotDueMarkedDueDate,omitempty"`
	EarliestRefillPullForwardDate *time.Time `json:"earliestRefillPullForwardDate" bson:"earliestRefillPullForwardDate,omitempty"`
	RefillNotDueCount             int        `json:"refillNotDueCount" bson:"refillNotDueCount"`
	RefillWindow                  int        `json:"refillWindow" bson:"refillWindow"`
	DoctorName                    string     `json:"doctorName" bson:"doctorName,omitempty"`
	DoctorPhoneNumber             string     `json:"doctorPhoneNumber" bson:"doctorPhoneNumber,omitempty"`
	DoctorFaxNumber               string     `json:"doctorFaxNumber" bson:"doctorFaxNumber,omitempty"`
	SendCopayConfirmEmail         bool       `json:"sendCopayConfirmEmail" bson:"sendCopayConfirmEmail"`
	HighCopay                     bool       `json:"highCopay" bson:"highCopay"`
	HighCopayReason               string     `json:"highCopayReason" bson:"highCopayReason,omitempty"`
	CopayChanged                  bool       `json:"copayChanged" bson:"copayChanged"`
	CopayChangedReason            string     `json:"copayChangedReason" bson:"copayChangedReason,omitempty"`
	SwitchedFromCash              bool       `json:"switchedFromCash" bson:"switchedFromCash"`
	StepTherapyRequired           bool       `json:"stepTherapyRequired" bson:"stepTherapyRequired"`
	MDApprovedOverPhone           bool       `json:"mdApprovedOverPhone" bson:"mdApprovedOverPhone"`
	CallCustomerForRefillNotDue   bool       `json:"callCustomerForRefillNotDue" bson:"callCustomerForRefillNotDue"`
	DiagnosisCode                 string     `json:"diagnosisCode" bson:"diagnosisCode"`
	MedicationHistory             string     `json:"medicationHistory" bson:"medicationHistory,omitempty"`
	AssociatedLabTests            bool       `json:"associatedLabTests" bson:"associatedLabTests"`
	SendToPP                      bool       `json:"sendToPP" bson:"sendToPP"`
}

// PaymentMilestone represents payment processing step
type PaymentMilestone struct {
	Milestone                  `bson:",inline"`
	PaymentError               string           `json:"paymentError" bson:"paymentError,omitempty"`
	ErrorCode                  string           `json:"errorCode" bson:"errorCode,omitempty"`
	ChargeAttempt              int              `json:"chargeAttempt" bson:"chargeAttempt"`
	AppliedPromo               string           `json:"appliedPromo" bson:"appliedPromo,omitempty" pson:"applied_promo"`
	TransactionID              string           `json:"transactionId" bson:"transactionId,omitempty"`
	PaymentProviderUsed        *PaymentProvider `json:"paymentProviderUsed" bson:"paymentProviderUsed,omitempty"`
	ChargeType                 string           `json:"chargeType" bson:"chargeType,omitempty"`
	AppliedDummyCharge         bool             `json:"appliedDummyCharge" bson:"appliedDummyCharge" pson:"applied_dummy_charge"`
	AppliedFreeTrialCharge     bool             `json:"appliedFreeTrialCharge" bson:"appliedFreeTrialCharge" pson:"applied_free_trial_charge"`
	MarkedShipLater            bool             `json:"markedShipLater" bson:"markedShipLater"`
	MedicationPrice            string           `json:"medicationPrice" bson:"medicationPrice,omitempty" pson:"medication_price"`
	DeliveryCharge             string           `json:"deliveryCharge" bson:"deliveryCharge,omitempty" pson:"delivery_charge"`
	FinalCharge                string           `json:"finalCharge" bson:"finalCharge,omitempty" pson:"final_charge"`
	AppliedPromoDiscount       string           `json:"appliedPromoDiscount" bson:"appliedPromoDiscount,omitempty" pson:"applied_promo_discount"`
	PaymentProcessedDate       *time.Time       `json:"paymentProcessedDate" bson:"paymentProcessedDate,omitempty"`
	ErrorAlertSentDate         *time.Time       `json:"errorAlertSentDate" bson:"errorAlertSentDate,omitempty"`
	LastReminderSentDate       *time.Time       `json:"lastReminderSentDate" bson:"lastReminderSentDate,omitempty"`
	PaymentReminderCount       int              `json:"paymentReminderCount" bson:"paymentReminderCount"`
	HighPrice                  bool             `json:"highPrice" bson:"highPrice"`
	HighPriceReason            string           `json:"highPriceReason" bson:"highPriceReason,omitempty"`
	PriceChanged               bool             `json:"priceChanged" bson:"priceChanged"`
	PriceChangedReason         string           `json:"priceChangedReason" bson:"priceChangedReason,omitempty"`
	ShouldReverseCashClaim     bool             `json:"shouldReverseCashClaim" bson:"shouldReverseCashClaim,omitempty"`
	ShouldCallPatientForSignup bool             `json:"shouldCallPatientForSignup" bson:"shouldCallPatientForSignup,omitempty"`
}

// DeliveryMilestone represents shipping/delivery processing step
type DeliveryMilestone struct {
	Milestone                           `bson:",inline"`
	ShippingRateID                      string     `json:"shippingRateId" bson:"shippingRateId,omitempty"`
	TrackingNumber                      string     `json:"trackingNumber" bson:"trackingNumber,omitempty"`
	TrackingURL                         string     `json:"trackingUrl" bson:"trackingUrl,omitempty"`
	TrackingURLForMobile                string     `json:"trackingUrlForMobile" bson:"trackingUrlForMobile,omitempty"`
	ShippingLabelURL                    string     `json:"shippingLabelUrl" bson:"shippingLabelUrl,omitempty"`
	ShipmentTransactionID               string     `json:"shipmentTransactionId" bson:"shipmentTransactionId,omitempty"`
	ShippingServiceName                 string     `json:"shippingServiceName" bson:"shippingServiceName,omitempty" pson:"shipping_service_name"`
	GenerateLabelError                  string     `json:"generateLabelError" bson:"generateLabelError,omitempty"`
	ShippingLabelCost                   string     `json:"shippingLabelCost" bson:"shippingLabelCost,omitempty" pson:"shipping_label_cost"`
	EstimatedDeliveryDays               int        `json:"estimatedDeliveryDays" bson:"estimatedDeliveryDays,omitempty"`
	EstimatedDeliveryDate               *time.Time `json:"estimatedDeliveryDate" bson:"estimatedDeliveryDate,omitempty"`
	ShippedDate                         *time.Time `json:"shippedDate" bson:"shippedDate,omitempty" pson:"shipped_date"`
	DeliveredDate                       *time.Time `json:"deliveredDate" bson:"deliveredDate,omitempty"`
	DeliveryConfirmed                   bool       `json:"deliveryConfirmed" bson:"deliveryConfirmed"`
	SameDayDelivery                     bool       `json:"sameDayDelivery" bson:"sameDayDelivery,omitempty"`
	UserSelectedSameDayPickupTime       bool       `json:"userSelectedSameDayPickupTime" bson:"userSelectedSameDayPickupTime,omitempty"`
	PharmacyPrintedSameDayShippingLabel bool       `json:"pharmacyPrintedSameDayShippingLabel" bson:"pharmacyPrintedSameDayShippingLabel,omitempty"`
	CarrierType                         string     `json:"carrierType" bson:"carrierType,omitempty"`
	PhilShipmentID                      string     `json:"philShipmentId" bson:"philShipmentId,omitempty"`
	DeliveryStatus                      string     `json:"deliveryStatus" bson:"deliveryStatus,omitempty"`
}

// RefillMilestone represents refill authorization processing step
type RefillMilestone struct {
	Milestone           `bson:",inline"`
	RxName              string     `json:"rxName" bson:"rxName,omitempty"`
	Quantity            int        `json:"quantity" bson:"quantity"`
	DaysOfSupply        int        `json:"daysOfSupply" bson:"daysOfSupply"`
	AuthorizedRefills   int        `json:"authorizedRefills" bson:"authorizedRefills"`
	DoctorName          string     `json:"doctorName" bson:"doctorName,omitempty"`
	DoctorPhoneNumber   string     `json:"doctorPhoneNumber" bson:"doctorPhoneNumber,omitempty"`
	DoctorFaxNumber     string     `json:"doctorFaxNumber" bson:"doctorFaxNumber,omitempty"`
	DoctorNpi           string     `json:"doctorNpi" bson:"doctorNpi,omitempty"`
	DoctorCity          string     `json:"doctorCity" bson:"doctorCity,omitempty"`
	FaxQueueID          string     `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	RefillRequestedDate *time.Time `json:"refillRequestedDate" bson:"refillRequestedDate,omitempty"`
	MDApprovedOverPhone bool       `json:"mdApprovedOverPhone" bson:"mdApprovedOverPhone"`
	Sig                 string     `json:"sig" bson:"sig,omitempty"`
	RefillAuthorizedBy  string     `json:"refillAuthorizedBy" bson:"refillAuthorizedBy,omitempty"`
}
