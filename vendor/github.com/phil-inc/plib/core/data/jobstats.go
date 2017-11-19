package data

const (
	JobStarted  string = "STARTED"
	JobComplete string = "COMPLETED"
	JobInfo     string = "INFO"
	JobError    string = "ERROR"
)

const (
	JobRefillAutoFax               string = "REFILL_AUTO_FAX"
	JobPushReminderPaymentApproval string = "PUSH_REMINDER_PAYMENT_APPROVAL"
	JobRefillRequested             string = "REFILL_REQUESTED"
	JobRefillDelayed               string = "REFILL_DELAYED"
	JobRefillNoResponse            string = "REFILL_NO_RESPONSE"
	JobTransferDelayed             string = "TRANSFER_DELAYED"
	JobRefillAuth                  string = "REFILL_AUTHORIZATION"
	JobRefillSkipNotification      string = "REFILL_SKIP_NOTIFICATION"
	JobCheckFaxStatus              string = "CHECK_FAX_STATUS"
	JobPhilInTownEmail             string = "PHIL_IN_TOWN_EMAIL"
	JobMissingInsuranceSMS         string = "MISSING_INSURANCE_SMS"
	JobShipNotifications           string = "SHIP_NOTIFICATIONS"
	JobLemonaidUserSignupReminder  string = "LEMONAID_SIGNUP_REMINDER"
)

//JobStats data representation for job stats
type JobStats struct {
	BaseData    `bson:",inline"`
	Name        string `json:"jobName" bson:"jobName" binding:"required"`
	Status      string `json:"status" bson:"status" binding:"required"`
	Description string `json:"description" bson:"description,omitempty"`
}

// CollectionName function to implement Document interface for promo
func (js JobStats) CollectionName() string {
	return "jobStats"
}

// FaxLog data representation for keeping fax log that's sent automatically.
type FaxLog struct {
	BaseData    `bson:",inline"`
	FaxQueueID  string `json:"faxQueueId" bson:"faxQueueId,omitempty"`
	OrderNumber string `json:"orderNumber" bson:"orderNumber,omitempty"`
}

//CollectionName implements Document interface
func (fl FaxLog) CollectionName() string {
	return "faxLog"
}

// AutoCallLog represents auto call
type AutoCallLog struct {
	RxID        string `json:"rxId" bson:"rxId"`
	OrderNumber string `json:"orderNumber" bson:"orderNumber"`
	PatientName string `json:"patientName" bson:"patientName"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	Label       string `json:"label" bson:"label"`
	SubLabel    string `json:"subLabel" bson:"subLabel"`
}
