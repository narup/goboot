package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"math/rand"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"strings"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/service"
	"github.com/phil-inc/plib/core/util"
	"golang.org/x/crypto/bcrypt"
)

// SignUpUser creates and save new user data based on passed in values.
// It also hashes the raw password uses default account status.
func SignUpUser(ctx context.Context, newUser *data.User) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	exists, err := IsExistingUserInSession(session, newUser.PhoneNumber, newUser.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		log.Printf("[ERRPR] user with email already exists. Duplicate sign up attempt. Email: %s", newUser.Email)
		return nil, errors.New("User exists")
	}

	newUser.InitData()

	if newUser.AccountStatus == "" {
		newUser.AccountStatus = data.AccountPending
	}
	newUser.Password = encryptedPassword(newUser.Password)
	if newUser.Password == "" {
		return nil, errors.New("Invalid password")
	}
	if newUser.AccountNumber == "" {
		newUser.AccountNumber = GenerateAccountNumber()
	}

	role := new(data.UserRole)
	if err := session.Find(gmgo.Q{"name": "ROLE_USER"}, role); err != nil {
		return newUser, err
	}

	roleRef := &mgo.DBRef{Collection: role.CollectionName(), Id: role.ID}
	newUser.RoleRefs = []*mgo.DBRef{roleRef}

	t := time.Now().UTC()
	newUser.CreatedDate = &t
	newUser.UpdatedDate = &t
	err = session.Save(newUser)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

//UpdatePassword updates user password based on user id
func UpdatePassword(ctx context.Context, userID, password string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	savedUsr, err := FindByIDInSession(ctx, userID, session)
	if err != nil {
		return nil, err
	}

	encryptedPwd := encryptedPassword(password)
	savedUsr.Password = encryptedPwd

	UnlockAccount(ctx, savedUsr)

	return savedUsr, session.Update(gmgo.Q{"_id": savedUsr.ID}, savedUsr)
}

// SaveUserSource save user source
func SaveUserSource(ctx context.Context, userSource *data.UserSource) error {
	session := data.Session()
	defer session.Close()

	us := new(data.UserSource)
	us.InitData()
	us.UserID = userSource.UserID
	us.Source = userSource.Source
	us.Channel = userSource.Channel
	us.CampaignParameters = userSource.CampaignParameters

	return session.Save(us)
}

//UpdateUserSource update user source data
func UpdateUserSource(ctx context.Context, userSource *data.UserSource) error {
	session := data.Session()
	defer session.Close()

	return session.Update(gmgo.Q{"_id": userSource.ID}, userSource)
}

//DeleteUser deletes the user based on given user id
func DeleteUser(ctx context.Context, userID string) error {
	session := data.Session()
	defer session.Close()

	return session.Remove(gmgo.Q{"_id": data.ObjectID(userID)}, new(data.User))
}

//Update update the user
func Update(ctx context.Context, user *data.User) error {
	session := data.Session()
	defer session.Close()

	return session.Update(gmgo.Q{"_id": user.ID}, user)
}

//UpdateInSession updates the given order using passed in session
func UpdateInSession(u *data.User, session *gmgo.DbSession) error {
	u.UpdatedDate = util.NowUTC()
	return session.Update(gmgo.Q{"_id": u.ID}, u)
}

//UpdateEmail updates user email based on user id
func UpdateEmail(ctx context.Context, userID, email string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	savedUsr, err := FindByIDInSession(ctx, userID, session)
	if err != nil {
		return nil, err
	}

	//nothing to update
	if savedUsr.Email == email {
		return savedUsr, nil
	}

	//check if new email already exists for a different user
	exists, err := session.Exists(gmgo.Q{"email": email}, new(data.User))
	if err != nil {
		return nil, errors.New("Error updating user email")
	}
	if exists {
		return nil, errors.New("Email already exists")
	}

	savedUsr.Email = email
	savedUsr.AccountStatus = data.AccountPending

	return savedUsr, session.Update(gmgo.Q{"_id": savedUsr.ID}, savedUsr)
}

// UpdateUserProfile updates the user data based on existing ObjectId
func UpdateUserProfile(ctx context.Context, user *data.User) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	savedUser := new(data.User)
	err := session.FindByID(user.StringID(), savedUser)
	if err != nil && err.Error() == data.ErrNotFound {
		return nil, errors.New("User not found")
	} else if err != nil {
		return nil, err
	}

	//check if user is trying to change email
	if user.Email != savedUser.Email {
		exists, err := session.Exists(gmgo.Q{"email": user.Email}, new(data.User))
		if err != nil {
			return nil, errors.New("Error updating user info")
		}
		if exists {
			return nil, errors.New("Email already exists")
		}
		savedUser.AccountStatus = data.AccountPending
	}

	savedUser.Email = user.Email
	savedUser.FullName = user.FullName
	savedUser.ZipCode = user.ZipCode

	return savedUser, session.Update(gmgo.Q{"_id": savedUser.ID}, savedUser)
}

// GenerateAccountNumber generates the random account number of format xxxxxxxx
func GenerateAccountNumber() string {
	session := data.Session()
	defer session.Close()

	accountNumber := generateAccountNumber()
	usr, _ := FindByAccountNumberInSession(accountNumber, session)
	for usr != nil {
		accountNumber = generateAccountNumber()
		usr, _ = FindByAccountNumberInSession(accountNumber, session)
	}
	return accountNumber
}

func generateAccountNumber() string {
	size := 8
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = digits[rand.Intn(len(digits))]
	}
	return string(buf)
}

// FindByAccountNumber finds the user by account number. References are not loaded.
func FindByAccountNumber(ctx context.Context, accountNumber string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	return FindByAccountNumberInSession(accountNumber, session)
}

// FindByAccountNumberInSession finds the usr by account number in session. References are not loaded.
func FindByAccountNumberInSession(accountNumber string, session *gmgo.DbSession) (*data.User, error) {
	usr := new(data.User)
	err := session.Find(gmgo.Q{"accountNumber": accountNumber}, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// FindByMissingInsuranceToken returns the user by given missing insurance token
func FindByMissingInsuranceToken(ctx context.Context, token string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	missingInsuranceToken := new(data.MissingInsuranceToken)
	err := session.Find(gmgo.Q{"token": token}, missingInsuranceToken)
	if err != nil {
		return nil, err
	}

	return FindByIDInSession(ctx, missingInsuranceToken.UserID, session)
}

//FindUserForInvalidInsuranceToken finds user from token for invalid insurance
func FindUserForInvalidInsuranceToken(ctx context.Context, token string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	it := new(data.InvalidInsuranceToken)
	err := session.Find(gmgo.Q{"token": token}, it)
	if err != nil {
		return nil, err
	}

	return FindByIDInSession(ctx, it.UserID, session)
}

// FindByID returns the user by given id
func FindByID(ctx context.Context, id string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	usr, err := FindByIDInSession(ctx, id, session)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// FindByEmail returns user complete data for given email
func FindByEmail(ctx context.Context, email string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	usr, err := FindByEmailInSession(ctx, email, session)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// FindByPhoneNumber returns user complete data for given phone number
func FindByPhoneNumber(ctx context.Context, phoneNumber string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	usr, err := FindByPhoneNumberInSession(ctx, phoneNumber, session)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

//FindByIDInSession finds the user on the given db session for the given user id
func FindByIDInSession(ctx context.Context, userID string, session *gmgo.DbSession) (*data.User, error) {
	user := new(data.User)
	err := session.FindByID(userID, user)
	if err != nil {
		return nil, err
	}
	return user, LoadAllUserReferences(user, session)
}

//FindByEmailInSession finds user for the given email in the given session
func FindByEmailInSession(ctx context.Context, email string, session *gmgo.DbSession) (*data.User, error) {
	user := new(data.User)
	if err := session.Find(gmgo.Q{"email": email}, user); err != nil {
		return user, err
	}

	return user, LoadAllUserReferences(user, session)
}

//FindByPhoneNumberInSession finds user for the given phone number in the given session
func FindByPhoneNumberInSession(ctx context.Context, phoneNumber string, session *gmgo.DbSession) (*data.User, error) {
	user := new(data.User)
	if err := session.Find(gmgo.Q{"phoneNumber": phoneNumber}, user); err != nil {
		return user, err
	}

	return user, LoadAllUserReferences(user, session)
}

func FindByEmailOrPhoneInSession(ctx context.Context, session *gmgo.DbSession, email, phoneNumber string) (*data.User, error) {
	usr, err := FindByEmailInSession(ctx, email, session)
	if err == nil {
		return usr, err
	}

	usr, err = FindByPhoneNumberInSession(ctx, phoneNumber, session)
	return usr, err
}

//FindAllPatientsForManager find patients for the manager
func FindAllPatientsForManager(ctx context.Context, managerID string) ([]*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{"manager.$id": data.ObjectID(managerID)}, new(data.Patient))
	if err != nil {
		return nil, err
	}
	return result.([]*data.Patient), nil
}

//FindPatientByID finds patient by given ID
func FindPatientByID(ctx context.Context, patientID string) (*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	return findPatientByID(session, patientID)
}

//FindPatientByPhoneNumber finds patient by phone number
func FindPatientByPhoneNumber(ctx context.Context, phoneNumber string) (*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	if strings.HasPrefix(phoneNumber, ")") {
		phoneNumber = strings.Replace(phoneNumber, "(", "", 1)
		phoneNumber = strings.Replace(phoneNumber, ")", "", 1)
		phoneNumber = strings.Replace(phoneNumber, "-", "", 1)
	}
	if strings.HasPrefix(phoneNumber, "+1") {
		phoneNumber = strings.Replace(phoneNumber, "+1", "", 1)
	}
	patient := new(data.Patient)
	err := session.Find(gmgo.Q{"phoneNumber": phoneNumber}, patient)
	if err != nil {
		return nil, err
	}

	mgr, err := data.ManagerByRef(patient.ManagerRef, session)
	if err != nil {
		return nil, err
	}
	patient.Manager = mgr

	return patient, nil
}

// SavePatient saves new patient
func SavePatient(ctx context.Context, np *data.Patient) (*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	patient := new(data.Patient)
	patient.InitData()
	patient.PatientName = strings.Title(np.PatientName)
	patient.PhoneNumber = util.SanitizePhoneNumber(np.PhoneNumber)
	patient.DateOfBirth = np.DateOfBirth
	patient.Allergies = np.Allergies
	patient.ManagerRef = data.DBRef(np.Manager.CollectionName(), np.Manager.ID)
	patient.Manager = np.Manager
	return patient, session.Save(patient)
}

// UpdatePatient update patient
func UpdatePatient(ctx context.Context, up *data.Patient) (*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	savedPt := new(data.Patient)
	err := session.Find(gmgo.Q{"_id": up.ID}, savedPt)
	if err != nil {
		return nil, err
	}

	savedPt.PatientName = strings.Title(up.PatientName)
	savedPt.PhoneNumber = util.SanitizePhoneNumber(up.PhoneNumber)
	savedPt.Allergies = up.Allergies
	savedPt.DefaultAddressID = up.DefaultAddressID
	savedPt.DefaultPaymentID = up.DefaultPaymentID
	savedPt.DefaultInsuranceID = up.DefaultInsuranceID
	if savedPt.DateOfBirth != up.DateOfBirth {
		savedPt.DateOfBirth = up.DateOfBirth
	}
	return savedPt, session.Update(gmgo.Q{"_id": savedPt.ID}, savedPt)
}

// UpdatePatientDob update patient Date of birth
func UpdatePatientDob(ctx context.Context, patientID string, newDob *time.Time) (*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	savedPt := new(data.Patient)
	err := session.Find(gmgo.Q{"_id": data.ObjectID(patientID)}, savedPt)
	if err != nil {
		return nil, err
	}

	savedPt.DateOfBirth = newDob
	return savedPt, session.Update(gmgo.Q{"_id": savedPt.ID}, savedPt)
}

//DeletePatient deletes the patient based on given patient id
func DeletePatient(ctx context.Context, patientID string) error {
	session := data.Session()
	defer session.Close()

	return session.Remove(gmgo.Q{"_id": data.ObjectID(patientID)}, new(data.Patient))
}

// FindUserConfig find user configiguration
func FindUserConfig(ctx context.Context, userID string) (*data.UserConfig, error) {
	session := data.Session()
	defer session.Close()

	cfg := new(data.UserConfig)
	err := session.Find(gmgo.Q{"userId": userID}, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// FindUserSource finds user source data for given user id
func FindUserSource(ctx context.Context, userID string) (*data.UserSource, error) {
	session := data.Session()
	defer session.Close()

	source := new(data.UserSource)
	err := session.Find(gmgo.Q{"userId": userID}, source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

//FindAllUserSource returns all user sources
func FindAllUserSource(ctx context.Context) ([]*data.UserSource, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.UserSource))
	if err != nil {
		return nil, err
	}

	return results.([]*data.UserSource), nil
}

//FindPatients returns all the patients for a given manager ID
func FindPatients(ctx context.Context, managerID string) ([]*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"manager.$id": data.ObjectID(managerID)}, new(data.Patient))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Patient), nil
}

//FindAllPatients returns all patients
func FindAllPatients(ctx context.Context) ([]*data.Patient, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Patient))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Patient), nil
}

//FindAllInsurances returns all insurances
func FindAllInsurances(ctx context.Context) ([]*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Insurance))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Insurance), nil
}

//FindAllShipments returns all shipments
func FindAllShipments(ctx context.Context) ([]*data.Shipment, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Shipment))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Shipment), nil
}

//FindAllUsers returns all the managers
func FindAllUsers(ctx context.Context) ([]*data.User, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.User))
	if err != nil {
		return nil, err
	}
	return results.([]*data.User), nil
}

//SearchUsers returns all the managers
func SearchUsers(ctx context.Context, searchStr string) ([]*data.User, error) {
	session := data.Session()
	defer session.Close()

	searchRegStr := "^" + searchStr

	searchQuery := gmgo.Q{"$or": []interface{}{
		bson.M{"email": bson.M{"$regex": bson.RegEx{searchRegStr, "i"}}},
		bson.M{"fullName": bson.M{"$regex": bson.RegEx{searchRegStr, "i"}}},
		bson.M{"phoneNumber": bson.M{"$regex": bson.RegEx{searchRegStr, "i"}}},
	}}

	results, err := session.FindAll(searchQuery, new(data.User))
	if err != nil {
		return nil, err
	}
	return results.([]*data.User), nil
}

//FindUserForMissingInsuranceToken finds user token for missing insurance info
func FindUserForMissingInsuranceToken(ctx context.Context, token string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	it := new(data.MissingInsuranceToken)
	err := session.Find(gmgo.Q{"token": token}, it)
	if err != nil {
		return nil, err
	}

	return FindByIDInSession(ctx, it.UserID, session)
}

// FindPatientNotes returns all the patient notes based on patient ID
func FindPatientNotes(ctx context.Context, patientID string) ([]*data.PatientNote, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"patientId": patientID}, new(data.PatientNote))
	return results.([]*data.PatientNote), err
}

// SavePayment saves new payment
func SavePayment(ctx context.Context, np *data.Payment) (*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	payment := new(data.Payment)
	payment.InitData()
	payment.CardName = np.CardName
	payment.Last4Digit = np.Last4Digit
	payment.ExpiryMonth = np.ExpiryMonth
	payment.ExpiryYear = np.ExpiryYear
	payment.Cvc = np.Cvc
	payment.FsaHsaCard = np.FsaHsaCard
	payment.PaymentToken = np.PaymentToken
	payment.Provider = np.Provider
	payment.ProviderCustomerID = np.ProviderCustomerID

	payment.Manager = np.Manager
	payment.ManagerRef = data.DBRef(np.Manager.CollectionName(), np.Manager.ID)

	payment.CardNumber = ""
	return payment, session.Save(payment)
}

// UpdatePayment update payment with given payment data
func UpdatePayment(ctx context.Context, up *data.Payment) (*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	savedPayment := new(data.Payment)
	err := session.Find(gmgo.Q{"_id": up.ID}, savedPayment)
	if err != nil {
		return nil, err
	}
	savedPayment.Cvc = up.Cvc
	savedPayment.ExpiryMonth = up.ExpiryMonth
	savedPayment.ExpiryYear = up.ExpiryYear

	savedPayment.UpdatedDate = util.NowUTC()
	return savedPayment, session.Update(gmgo.Q{"_id": savedPayment.ID}, savedPayment)
}

// FindUserPayments returns all the saved payments for given user id
func FindUserPayments(ctx context.Context, userID string) ([]*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"manager.$id": data.ObjectID(userID)}, new(data.Payment))
	return results.([]*data.Payment), err
}

//DeletePayment deletes the patient based on given payment id
func DeletePayment(ctx context.Context, paymentID string) error {
	session := data.Session()
	defer session.Close()

	return session.Remove(gmgo.Q{"_id": data.ObjectID(paymentID)}, new(data.Payment))
}

// FindUserAddresses returns all the saved addresses for the given user ID
func FindUserAddresses(ctx context.Context, userID string) ([]*data.Address, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"manager.$id": data.ObjectID(userID)}, new(data.Address))
	return results.([]*data.Address), err
}

//FindAllAddresses return all addresses
func FindAllAddresses(ctx context.Context) ([]*data.Address, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Address))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Address), nil
}

// SaveAddress saves new address
func SaveAddress(ctx context.Context, newAddr *data.Address) (*data.Address, error) {
	session := data.Session()
	defer session.Close()

	addr := new(data.Address)
	addr.InitData()
	addr.Street1 = newAddr.Street1
	addr.Street2 = newAddr.Street2
	addr.City = newAddr.City
	addr.State = newAddr.State
	addr.ZipCode = newAddr.ZipCode
	addr.Country = newAddr.Country
	addr.Manager = newAddr.Manager
	if newAddr.Country == "" {
		addr.Country = "US"
	}
	addr.PatientID = newAddr.PatientID

	addr.ManagerRef = data.DBRef(addr.Manager.CollectionName(), addr.Manager.ID)

	return addr, session.Save(addr)
}

// UpdateAddress update address with given address data
func UpdateAddress(ctx context.Context, ua *data.Address) (*data.Address, error) {
	session := data.Session()
	defer session.Close()

	savedAddr := new(data.Address)
	err := session.Find(gmgo.Q{"_id": ua.ID}, savedAddr)
	if err != nil {
		return nil, err
	}
	savedAddr.Street1 = ua.Street1
	savedAddr.Street2 = ua.Street2
	savedAddr.City = ua.City
	savedAddr.State = ua.State
	savedAddr.ZipCode = ua.ZipCode
	savedAddr.Country = ua.Country
	if savedAddr.Country == "" {
		savedAddr.Country = "US"
	}
	savedAddr.PatientID = ua.PatientID

	savedAddr.UpdatedDate = util.NowUTC()
	return savedAddr, session.Update(gmgo.Q{"_id": savedAddr.ID}, savedAddr)
}

//DeleteAddress deletes the patient based on given address id
func DeleteAddress(ctx context.Context, addressID string) error {
	session := data.Session()
	defer session.Close()

	return session.Remove(gmgo.Q{"_id": data.ObjectID(addressID)}, new(data.Address))
}

// FindUserInsurances returns all the saved insurances for the given user ID
func FindUserInsurances(ctx context.Context, userID string) ([]*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"manager.$id": data.ObjectID(userID)}, new(data.Insurance))
	return results.([]*data.Insurance), err
}

// SaveInsurance saves new insurance
func SaveInsurance(ctx context.Context, ni *data.Insurance) (*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	insr := new(data.Insurance)
	insr.InitData()
	insr.ClientName = ni.ClientName
	insr.ProviderName = ni.ProviderName
	insr.GroupNumber = ni.GroupNumber
	insr.BinNumber = ni.BinNumber
	insr.PcnNumber = ni.PcnNumber
	insr.InsuranceID = ni.InsuranceID
	insr.CardImageID = ni.CardImageID
	insr.Primary = ni.Primary
	insr.IsGovtSponsored = ni.IsGovtSponsored
	insr.Manager = ni.Manager
	insr.PatientID = ni.PatientID

	insr.ManagerRef = data.DBRef(ni.Manager.CollectionName(), ni.Manager.ID)

	return insr, session.Save(insr)
}

// UpdateInsurance update insurance with given new insurance data
func UpdateInsurance(ctx context.Context, ui *data.Insurance) (*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	savedInsr := new(data.Insurance)
	err := session.Find(gmgo.Q{"_id": ui.ID}, savedInsr)
	if err != nil {
		return nil, err
	}

	savedInsr.ClientName = ui.ClientName
	savedInsr.ProviderName = ui.ProviderName
	savedInsr.GroupNumber = ui.GroupNumber
	savedInsr.BinNumber = ui.BinNumber
	savedInsr.PcnNumber = ui.PcnNumber
	savedInsr.InsuranceID = ui.InsuranceID
	savedInsr.CardImageID = ui.CardImageID
	savedInsr.Primary = ui.Primary
	savedInsr.IsGovtSponsored = ui.IsGovtSponsored
	savedInsr.PatientID = ui.PatientID

	savedInsr.UpdatedDate = util.NowUTC()
	return savedInsr, session.Update(gmgo.Q{"_id": ui.ID}, savedInsr)
}

//DeleteInsurance deletes the insurance based on given address id
func DeleteInsurance(ctx context.Context, addressID string) error {
	session := data.Session()
	defer session.Close()

	return session.Remove(gmgo.Q{"_id": data.ObjectID(addressID)}, new(data.Address))
}

// LoadAllUserReferences loads all user references
func LoadAllUserReferences(user *data.User, session *gmgo.DbSession) error {
	roles := make([]data.UserRole, len(user.RoleRefs))
	for i, roleRef := range user.RoleRefs {
		role := new(data.UserRole)
		if err := session.FindByRef(roleRef, role); err != nil {
			return util.HandleRefLoadError("Error loading user role", []error{err})
		}

		roles[i] = role.CopyValue()
	}

	//set roles
	user.Roles = roles
	return nil
}

//MentionOnRxOrder mention user for any Rx order
func MentionOnRxOrder(ctx context.Context, assigneeID, assignerID, rxID, orderNumber, message string) error {
	session := data.Session()
	defer session.Close()

	assignee, err := FindByIDInSession(ctx, assigneeID, session)
	if err != nil {
		return err
	}

	assignerName := "System Phil"
	if assignerID != "system" {
		assigner, err := FindByIDInSession(ctx, assignerID, session)
		if err != nil {
			return err
		}
		assignerName = assigner.FullName
	}

	m := new(data.Mention)
	m.InitData()
	m.AgentID = assignee.StringID()
	m.AgentName = assignee.FullName
	m.RxID = rxID
	m.OrderNumber = orderNumber
	m.AssignerAgentID = assignerID
	m.AssignerAgentName = assignerName
	m.Message = message
	m.Comment = message

	err = session.Save(m)
	if err != nil {
		return err
	}

	return nil
}

//MentionByID finds user mention by ID
func MentionByID(ctx context.Context, mentionID string) (*data.Mention, error) {
	session := data.Session()
	defer session.Close()

	m := new(data.Mention)
	return m, session.FindByID(mentionID, m)
}

// ClearMentionsByID clears mentions from a list of IDs
func ClearMentionsByIDs(ctx context.Context, mentionIDs []string) error {
	session := data.Session()
	defer session.Close()

	for _, mentionID := range mentionIDs {
		objectID := data.ObjectID(mentionID)
		session.UpdateFieldValue(gmgo.Q{"_id": objectID}, new(data.Mention).CollectionName(), "resolved", true)
	}

	return nil
}

//GenerateEmailVerificationLink generates email verification link
func GenerateEmailVerificationLink(email string) string {
	token := util.RandomToken(6)
	err := SaveEmailVerificationToken(context.Background(), email, token)
	if err != nil {
		log.Printf("[ERROR] Failed to save email verification token %s", token)
		return ""
	}
	return fmt.Sprintf("%s/email/verify?token=%s", util.Config("email.landingPageUrl"), token)
}

//GeneratePasswordResetLink generates password reset link
func GeneratePasswordResetLink(usr *data.User) string {
	token := util.RandomToken(6)
	err := SavePasswordResetToken(context.Background(), usr, token)
	if err != nil {
		log.Printf("[ERROR] Failed to save password reset token %s", token)
		return ""
	}
	return fmt.Sprintf("%s/password/reset?token=%s", util.Config("email.landingPageUrl"), token)
}

//GenerateMdDashPasswordResetLink generates password reset link for the MD Dashboard
func GenerateMdDashPasswordResetLink(usr *data.User) string {
	token := util.RandomToken(6)
	err := SavePasswordResetToken(context.Background(), usr, token)
	if err != nil {
		log.Printf("[ERROR] Failed to save password reset token %s", token)
		return ""
	}
	return fmt.Sprintf("%s/password-reset/%s", util.Config("dashboard.md.url"), token)
}

// SaveEmailVerificationToken save email verification token
func SaveEmailVerificationToken(ctx context.Context, email, token string) error {
	et := data.EmailVerificationToken{Email: email, Token: token}

	session := data.Session()
	defer session.Close()

	return session.Save(et)
}

// SavePasswordResetToken saves password reset token
func SavePasswordResetToken(ctx context.Context, usr *data.User, token string) error {
	session := data.Session()
	defer session.Close()

	rt := new(data.PasswordResetToken)
	rt.Token = token
	rt.UserRef = data.DBRef(usr.CollectionName(), usr.ID)

	return session.Save(rt)
}

//CreateReferralPromoCode create a promo code and saves it
func CreateReferralPromoCode(user *data.User) (*data.Promo, error) {
	session := data.Session()
	code := "REFER_PHIL"
	names := strings.Split(user.FullName, " ")
	if len(names) == 2 {
		code = names[0][:1] + names[1]
	} else if len(names) > 2 {
		code = names[2][:1] + names[0]
	}

	promo := &data.Promo{
		Code:            strings.ToUpper(code),
		CashCredit:      "10.0",
		CopayCredit:     "10.0",
		ShippingCredit:  "2.0",
		Campaign:        "referral_email",
		MultipleSupport: true,
		Active:          true,
		Recurring:       true,
	}

	err := session.Save(promo)
	return promo, err
}

//FindUsersByQueryInSession - find all users based on query in a given gmgo session
func FindUsersByQueryInSession(rq gmgo.Q, session *gmgo.DbSession) ([]*data.User, error) {
	results, err := session.FindAll(rq, new(data.User))
	if err != nil {
		return nil, err
	}
	return results.([]*data.User), nil
}

// FindAllMDOfficeUsers returns the all the md office users
func FindAllMDOfficeUsers() ([]*data.MDOfficeUser, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.MDOfficeUser))
	if err != nil {
		return nil, err
	}
	return results.([]*data.MDOfficeUser), nil
}

// FindMDOfficeUser returns the user data of the user from an md office that corresponds to the particular userID.
func FindMDOfficeUser(userID string) (*data.MDOfficeUser, error) {
	session := data.Session()
	defer session.Close()

	mdOUsr := new(data.MDOfficeUser)
	err := session.Find(gmgo.Q{"userId": userID}, mdOUsr)
	if err != nil {
		return nil, err
	}

	return mdOUsr, nil
}

// UpdateMDOfficeUser saves the md office user
func UpdateMDOfficeUser(mdOfficeUser *data.MDOfficeUser) error {
	session := data.Session()
	defer session.Close()

	return session.Update(gmgo.Q{"_id": mdOfficeUser.ID}, mdOfficeUser)
}

////// private functions/////
func encryptedPassword(rawPassword string) string {
	pb := []byte(rawPassword)
	// Hashing the password with the default cost of 10
	hashedPassword, hashError := bcrypt.GenerateFromPassword(pb, bcrypt.DefaultCost)
	if hashError != nil {
		log.Printf("[ERROR] Password encryption failed. Err: %s", hashError)
		return ""
	}
	return string(hashedPassword)
}

func findPatientByID(session *gmgo.DbSession, patientID string) (*data.Patient, error) {
	patient := new(data.Patient)
	err := session.Find(gmgo.Q{"_id": data.ObjectID(patientID)}, patient)
	if err != nil {
		return nil, err
	}

	//load manager
	mgr, err := data.ManagerByRef(patient.ManagerRef, session)
	if err != nil {
		return nil, err
	}

	patient.Manager = mgr

	return patient, nil
}

func IsExistingUser(phoneNumber, email string) (bool, error) {
	session := data.Session()
	defer session.Close()

	return IsExistingUserInSession(session, phoneNumber, email)
}

func IsExistingUserInSession(session *gmgo.DbSession, phoneNumber, email string) (bool, error) {
	existingUser := new(data.User)
	if phoneNumber != "" {
		isWhiteListedPhoneNumber, _ := base.IsWhiteListedPhoneNumber(context.Background(), phoneNumber)
		if isWhiteListedPhoneNumber {
			return false, nil
		}
		exists, err := session.Exists(gmgo.Q{"phoneNumber": phoneNumber}, existingUser)
		if exists || err != nil {
			return exists, err
		}
	}
	return session.Exists(gmgo.Q{"email": email}, existingUser)
}
