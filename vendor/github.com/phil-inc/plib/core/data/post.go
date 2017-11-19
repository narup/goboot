package data

import (
	"time"
)

const (
	TypeTicket       = "Ticket"
	TypeEvent        = "Event"
	TypeInfo         = "Info"
	TypeAnnouncement = "Announcement"
	TypeInternal     = "Internal"
	TypeIgnore       = "Ignore"

	SourceEmail     = "Email"
	SourceSMS       = "SMS"
	SourceApp       = "App"
	SourceDashboard = "Dashboard"

	StateRead     = "Read"
	StateNew      = "New"
	StateResolved = "Resolved"
	StateFollowup = "Followup"
)

// Post represents the feed post data
// Indexes:
// 1. db.post.createIndex({"rxInfos.rxId": 1, "rxInfos.orderNumber" : 1})
// 2. db.post.createIndex({"creator.id": 1, "creator.id" : 1})
// 3. db.post.createIndex({"user.id": 1, "user.id" : 1})
// 4. db.post.createIndex({"patientName": 1})
type Post struct {
	BaseData       `bson:",inline"`
	User           *UserInfo     `json:"user" bson:"user"`
	Creator        *UserInfo     `json:"creator" bson:"creator"`
	Updater        *UserInfo     `json:"updater" bson:"updater,omitempty"`
	PatientID      string        `json:"patientId" bson:"patientId"`
	PatientName    string        `json:"patientName" bson:"patientName"`
	PhoneNumber    string        `json:"phoneNumber" bson:"phoneNumber"`
	Content        string        `json:"content" bson:"content,omitempty"`
	RawContent     string        `json:"rawContent" bson:"rawContent,omitempty"`
	Type           string        `json:"type" bson:"type"`
	TrackingNumber string        `json:"trackingNumber" bson:"trackingNumber,omitempty"`
	HighPriority   int           `json:"highPriority" bson:"highPriority,omitempty"`
	Source         string        `json:"source" bson:"source,omitempty"`
	ActionID       string        `json:"actionId" bson:"actionId,omitempty"`
	LastAction     string        `json:"lastAction" bson:"lastAction,omitempty"`
	LastUserAction string        `json:"lastUserAction" bson:"lastUserAction,omitempty"`
	CreatorType    string        `json:"creatorType" bson:"creatorType,omitempty"`
	Tags           []string      `json:"tags" bson:"tags,omitempty"`
	Deleted        bool          `json:"deleted" bson:"deleted,omitempty"`
	DeletedDate    *time.Time    `json:"deletedDate" bson:"deletedDate,omitempty"`
	UserStatus     *Status       `json:"userStatus" bson:"userStatus"`
	SupportStatus  *Status       `json:"supportStatus" bson:"supportStatus"`
	Category       *Category     `json:"category" bson:"category,omitempty"`
	RxInfos        []*RxInfo     `json:"rxInfos" bson:"rxInfos,omitempty"`
	Replies        []*Reply      `json:"replies" bson:"replies,omitempty"`
	Notes          []*Note       `json:"notes" bson:"notes,omitempty"`
	Attachments    []*Attachment `json:"attachments" bson:"attachments,omitempty"`
	Ticket         *Ticket       `json:"ticket" bson:"ticket,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (p Post) CollectionName() string {
	return "post"
}

// AddPatientInfo add patient info to the post
func (p *Post) AddPatientInfo(patient *Patient) {
	p.PatientID = patient.StringID()
	p.PatientName = patient.PatientName
}

// AddRxInfo adds rx info to the post
func (p *Post) AddRxInfo(rx *Prescription) {
	rxInfo := new(RxInfo)
	rxInfo.RxID = rx.StringID()
	rxInfo.MedName = rx.Name
	rxInfo.FullMedName = rx.FullRxName()
	rxInfo.OrderNumber = rx.CurrentOrder.OrderNumber

	p.AddPatientInfo(rx.Patient)
	if p.RxInfos == nil {
		p.RxInfos = make([]*RxInfo, 0)
	}
	p.RxInfos = append(p.RxInfos, rxInfo)
}

// AssignCreator assigns creator info to post
func (p *Post) AssignCreator(creator *User) {
	p.Creator = postUserInfo(creator)
}

//AssignUser assigns user to the post
func (p *Post) AssignUser(user *User) {
	p.User = postUserInfo(user)
}

//AssignUpdater assign updater
func (p *Post) AssignUpdater(user *User) {
	p.Updater = postUserInfo(user)
}

func postUserInfo(user *User) *UserInfo {
	userInfo := new(UserInfo)
	userInfo.ID = user.StringID()
	userInfo.Name = user.FullName
	userInfo.Email = user.Email

	return userInfo
}

//Category defines post category with main and secondary type.
//Description gives detail on the category if it's relevant. Label is used for filtering
type Category struct {
	Main        string `json:"main" bson:"main,omitempty"`
	Secondary   string `json:"secondary" bson:"secondary,omitempty"`
	Description string `json:"description" bson:"description,omitempty"`
}

//Status represents the post status
type Status struct {
	State           string     `json:"state" bson:"state"` //New, Read, Resolved, Followup
	OpenedDate      *time.Time `json:"openedDate" bson:"openedDate,omitempty"`
	ResolvedDate    *time.Time `json:"resolvedDate" bson:"resolvedDate,omitempty"`
	FollowupDate    *time.Time `json:"followupDate" bson:"followupDate,omitempty"`
	PendingCustomer bool       `json:"pendingCustomer" bson:"pendingCustomer,omitempty"`
}

// UserInfo data representation for user info
type UserInfo struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Email   string `json:"email" bson:"email"`
	ImageID string `json:"imageId" bson:"imageId,omitempty"`
}

// RxInfo rx info for the post
type RxInfo struct {
	MedName     string `json:"medName" bson:"medName,omitempty"`
	FullMedName string `json:"fullMedName" bson:"fullMedName,omitempty"`
	RxID        string `json:"rxId" bson:"rxId"`
	OrderNumber string `json:"orderNumber" bson:"orderNumber"`
}

//Ticket data representation for ticket
type Ticket struct {
	ID          string                 `json:"id" bson:"id"`
	Source      string                 `json:"source" bson:"source"` //HelpScout, HelpShift etc. incase we change the ticketing platform
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata,omitempty"`
	CreatedDate *time.Time             `json:"createdDate" bson:"createdDate"`
}

// Reply data reprsentation for replies
type Reply struct {
	ReplyID        string        `json:"replyId" bson:"replyId"`
	Creator        *UserInfo     `json:"creator" bson:"creator"`
	Content        string        `json:"content" bson:"content,omitempty"`
	RawContent     string        `json:"rawContent" bson:"rawContent,omitempty"`
	Edited         bool          `json:"edited" bson:"edited"`
	Deleted        bool          `json:"deleted" bson:"deleted"`
	Attachments    []*Attachment `json:"attachments" bson:"attachments,omitempty"`
	ReadDate       *time.Time    `json:"readDate" bson:"readDate,omitempty"`
	CreatedDate    *time.Time    `json:"createdDate" bson:"createdDate"`
	UpdatedDate    *time.Time    `json:"updatedDate" bson:"updatedDate"`
	ActionID       string        `json:"actionId" bson:"actionId,omitempty"`
	LastUserAction string        `json:"lastUserAction" bson:"lastUserAction,omitempty"`
	CreatorType    string        `json:"creatorType" bson:"creatorType,omitempty"`
}

//AssignCreator assigns creator to the post
func (r *Reply) AssignCreator(user *User) {
	r.Creator = postUserInfo(user)
}

//Note represents feed post notes
type Note struct {
	NoteID      string        `json:"noteId" bson:"noteId"`
	Creator     *UserInfo     `json:"creator" bson:"creator"`
	Content     string        `json:"content" bson:"content,omitempty"`
	Edited      bool          `json:"edited" bson:"edited"`
	Deleted     bool          `json:"deleted" bson:"deleted"`
	Attachments []*Attachment `json:"attachments" bson:"attachments,omitempty"`
	CreatedDate *time.Time    `json:"createdDate" bson:"createdDate"`
	UpdatedDate *time.Time    `json:"updatedDate" bson:"updatedDate"`
}

//AssignCreator assigns creator to the post
func (n *Note) AssignCreator(user *User) {
	n.Creator = postUserInfo(user)
}

//Attachment data representation for post or reply attachment
type Attachment struct {
	AttachmentID  string     `json:"attachmentId" bson:"attachmentId,omitempty"`
	AttachmentURL string     `json:"attachmentUrl" bson:"attachmentUrl,omitempty"`
	Name          string     `json:"name" bson:"name,omitempty"`
	Comment       string     `json:"comment" bson:"comment,omitempty"`
	CreatedDate   *time.Time `json:"createdDate" bson:"createdDate"`
	UpdatedDate   *time.Time `json:"updatedDate" bson:"updatedDate"`
}

//SystemUserInfo - captures user info for a system
var SystemUserInfo = UserInfo{ID: "system", Name: "Phil System", Email: "hello@phil.us", ImageID: ""}

//HelpscoutUserInfo - captures user info for helpscout
var HelpscoutUserInfo = UserInfo{ID: "helpscout", Name: "Help Scout", Email: "hello@phil.us", ImageID: ""}
