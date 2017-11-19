package prescription

import (
	"context"
	"errors"

	"log"

	"time"

	"fmt"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/service/address"
	"github.com/phil-inc/plib/core/service/doctor"
	"github.com/phil-inc/plib/core/service/insurance"
	"github.com/phil-inc/plib/core/service/order"
	"github.com/phil-inc/plib/core/service/user"
	"github.com/phil-inc/plib/core/util"
	"gopkg.in/mgo.v2/bson"
)

// Transfer rx origin type transfer
var Transfer = "Transfer"

// EScript reprents paper prescription
var EScript = "FromDoctorDirect"

// RxLoadConfig defines configuration to load prescription data
type RxLoadConfig struct {
	LoadPrescriptionRefs bool
	LoadOrderRefs        bool
	FilterArchived       bool
	FilterPaused         bool
	FilterPending        bool
	LoadPatient          bool
	LoadManager          bool
	LoadDoctor           bool
	LoadInsurance        bool
	LoadOrder            bool
	OrderLoadConfig      *order.OrderLoadConfig
}

// Init - initizalizes rxloadconfig
func (c *RxLoadConfig) Init() {
	c.LoadDoctor = true
	c.LoadPatient = true
	c.LoadManager = true
	c.LoadInsurance = true
	c.LoadOrder = true

	orderConfig := new(order.OrderLoadConfig)
	orderConfig.LoadAddress = true
	orderConfig.LoadPayment = true
	orderConfig.LoadOriginPharmacy = true
	orderConfig.LoadPartnerPharmacy = true

	c.OrderLoadConfig = orderConfig
}

// DefaultRxLoadConfig returns default rx load config
func DefaultRxLoadConfig() *RxLoadConfig {
	cfg := new(RxLoadConfig)

	cfg.Init()
	cfg.FilterArchived = true
	cfg.FilterPaused = true
	cfg.FilterPending = true
	cfg.LoadOrderRefs = true
	cfg.LoadPrescriptionRefs = true

	return cfg
}

// NoFilterLoadAllRefsConfig loads all prescription references
func NoFilterLoadAllRefsConfig() *RxLoadConfig {
	cfg := new(RxLoadConfig)

	cfg.Init()
	cfg.FilterArchived = false
	cfg.FilterPaused = false
	cfg.FilterPending = false
	cfg.LoadOrderRefs = true
	cfg.LoadPrescriptionRefs = true

	return cfg
}

// RxWithoutReferencesConfig - This should be used only to load prescription order data for read only operations.
func RxWithoutReferencesConfig() *RxLoadConfig {
	cfg := new(RxLoadConfig)

	cfg.Init()
	cfg.FilterArchived = true
	cfg.FilterPaused = true
	cfg.FilterPending = true
	cfg.LoadPrescriptionRefs = true
	cfg.LoadOrderRefs = false

	return cfg
}

// SaveInSession - saves the given prescription using a db session. Returns error if there's any while saving.
func SaveInSession(rx *data.Prescription, session *gmgo.DbSession) error {
	t := time.Now().UTC()

	if rx.ID.Hex() == "" {
		rx.CreatedDate = &t
	}
	rx.UpdatedDate = &t
	return session.Save(rx)
}

// UpdateInSession updates the given prescriptionn using passed in session
func UpdateInSession(p *data.Prescription, session *gmgo.DbSession) error {
	p.UpdatedDate = util.NowUTC()
	return session.Update(gmgo.Q{"_id": p.ID}, p)
}

// FindAllPrescriptions returns all prescriptions
func FindAllPrescriptions(ctx context.Context) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Prescription))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Prescription), nil
}

// FindAllForManager returns all the prescriptions for a given manager
// it also loads all the references for each prescription data
func FindAllForManager(ctx context.Context, manager *data.User) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	return FindAllForManagerInSession(manager.StringID(), session)
}

// FindAllForManagerWithID returns all the prescriptions for a given manager id
// it also loads all the references for each prescription data
func FindAllForManagerWithID(ctx context.Context, managerID string) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	return FindAllForManagerInSession(managerID, session)
}

// FindAllForManagerInSession find all the prescriptions for a given manager in db session
func FindAllForManagerInSession(managerID string, session *gmgo.DbSession) ([]*data.Prescription, error) {
	q := gmgo.Q{
		"manager.$id": data.ObjectID(managerID),
	}
	return FindPrescriptionsWithQuery(q, session)
}

// FindAllForPatient returns all the prescriptions for a given patient id
func FindAllForPatient(ctx context.Context, patientID string) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	q := gmgo.Q{
		"patient.$id": data.ObjectID(patientID),
	}

	return FindPrescriptionsWithQuery(q, session)
}

// FindByID finds prescription data and all its references by id. It returns
// complete data with all the DBRefs loaded and populated.
func FindByID(ctx context.Context, rxID string) (*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	return FindByIDInSession(rxID, session)
}

//FindByIDInSession finds prescription based on given rx ID and db session
func FindByIDInSession(rxID string, session *gmgo.DbSession) (*data.Prescription, error) {
	rx := new(data.Prescription)
	if err := session.FindByID(rxID, rx); err != nil {
		return rx, err
	}

	//load references based on DBRef
	return rx, LoadRxReferences(rx, session, true)
}

// FindPrescriptionsWithIDs - find the prescriptions for the given list of prescription IDs
func FindPrescriptionsWithIDs(rxIDs []bson.ObjectId, session *gmgo.DbSession) ([]*data.Prescription, error) {
	rq := gmgo.Q{
		"_id": gmgo.Q{
			"$in": rxIDs,
		},
	}
	return FindPrescriptionsWithQuery(rq, session)
}

// FindPrescriptionsWithQuery find all the prescriptions based on given query and session
func FindPrescriptionsWithQuery(rq gmgo.Q, session *gmgo.DbSession) ([]*data.Prescription, error) {
	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxReferences(rx, session, false)
		if err != nil {
			log.Printf("Error loading rx references for rx %s, %s\n", rx.StringID(), err)
		}
		if rx.CurrentOrder != nil {
			err = order.LoadReferences(rx.CurrentOrder, rx.Manager, session)
			if err != nil {
				log.Printf("Error loading order references for rx %s, %s\n", rx.StringID(), err)
			}
		} else {
			fmt.Printf("Missing order reference for rx %s\n", rx.StringID())
		}
	}
	return rxs, nil
}

//FindByOrderNumber find prescription by current order number
func FindByOrderNumber(orderNumber string) (*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	return FindByOrderNumberInSession(orderNumber, session)
}

//FindByOrderNumberInSession find prescription by current order number in session
func FindByOrderNumberInSession(orderNumber string, session *gmgo.DbSession) (*data.Prescription, error) {
	rxs, err := FindPrescriptionsWithOrderQueryInSession(gmgo.Q{"orderNumber": orderNumber}, nil, session)
	if err != nil {
		return nil, err
	}

	if len(rxs) > 0 {
		return rxs[0], nil
	}
	return nil, errors.New("not found")
}

// FindPrescriptionsWithOrderQuery - finds all the prescriptions based on the current order query
func FindPrescriptionsWithOrderQuery(q gmgo.Q, config *RxLoadConfig) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	return FindPrescriptionsWithOrderQueryInSession(q, config, session)
}

// FindPrescriptionsWithOrderQueryInSession finds all the prescriptions based on the current order query and config in given session
func FindPrescriptionsWithOrderQueryInSession(q gmgo.Q, config *RxLoadConfig, session *gmgo.DbSession) ([]*data.Prescription, error) {
	if config == nil {
		config = DefaultRxLoadConfig()
	}

	result, err := session.FindAllWithFields(q, []string{"_id"}, new(data.Order))
	if err != nil {
		return util.EmptyRxList(), err
	}

	orders := result.([]*data.Order)
	orderIds := make([]string, 0)
	for _, o := range orders {
		orderIds = append(orderIds, o.StringID())
	}

	rq := rxQuery(config, orderIds)
	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		if config.LoadPrescriptionRefs {
			err := LoadRxReferences(rx, session, false)
			if err != nil {
				log.Printf("Error loading rx references for rx %s, %s\n", rx.StringID(), err)
			}
			if rx.CurrentOrder != nil && config.LoadOrderRefs {
				err := order.LoadReferences(rx.CurrentOrder, rx.Manager, session)
				if err != nil {
					log.Printf("Error loading order references for rx %s, %s\n", rx.StringID(), err)
				}
			}
		}
	}
	return rxs, nil
}

// FindPrescriptionsWithOrderQueryAndLoadConfig - finds all the prescriptions based on the current order query and load config
func FindPrescriptionsWithOrderQueryAndLoadConfig(q gmgo.Q, config *RxLoadConfig) ([]*data.Prescription, error) {
	if config == nil {
		config = DefaultRxLoadConfig()
	}
	session := data.Session()
	defer session.Close()

	result, err := session.FindAllWithFields(q, []string{"_id"}, new(data.Order))
	if err != nil {
		return util.EmptyRxList(), err
	}

	orders := result.([]*data.Order)
	orderIds := make([]string, 0)
	for _, o := range orders {
		orderIds = append(orderIds, o.StringID())
	}

	rq := rxQuery(config, orderIds)
	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxReferencesWithConfig(rx, session, config)
		if err != nil {
			log.Printf("Error loading rx references for rx %s, %s\n", rx.StringID(), err)
		}
	}
	return rxs, nil
}

// FindPrescriptionsWithQueryAndLoadConfig find all the prescriptions based on given query and load config
func FindPrescriptionsWithQueryAndLoadConfig(rq gmgo.Q, config *RxLoadConfig, session *gmgo.DbSession) ([]*data.Prescription, error) {
	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxReferencesWithConfig(rx, session, config)
		if err != nil {
			log.Printf("Error loading rx references for rx %s, %s\n", rx.StringID(), err)
		}
	}
	return rxs, nil
}

// FindPrescriptionsByDoctorID finds all of the prescriptions associated with that doctor id.
func FindPrescriptionsByDoctorID(doctorID string) ([]*data.Prescription, error) {
	session := data.Session()
	defer session.Close()

	rxResults, err := session.FindAll(gmgo.Q{"doctor.$id": data.ObjectID(doctorID)}, new(data.Prescription))
	if err != nil {
		return nil, err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxReferences(rx, session, false)
		if err != nil {
			log.Printf("Error loading rx references for rx %s, %s\n", rx.StringID(), err)
		}
	}
	return rxs, nil
}

// FindStrippedDownPrescriptionsQuery find all the prescriptions with orders based on given query and session
func FindStrippedDownPrescriptionsQuery(rq gmgo.Q, session *gmgo.DbSession) ([]*data.Prescription, error) {
	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxOrder(rx, session)
		if err != nil {
			log.Printf("Error loading order for rx %s, %s\n", rx.StringID(), err)
		}
	}
	return rxs, nil
}

// FindStrippedDownPrescriptionsWithOrderQuery - finds all the prescriptions with orders based on the current order query
func FindStrippedDownPrescriptionsWithOrderQuery(q gmgo.Q, session *gmgo.DbSession) ([]*data.Prescription, error) {

	result, err := session.FindAllWithFields(q, []string{"_id"}, new(data.Order))
	if err != nil {
		return util.EmptyRxList(), err
	}

	orders := result.([]*data.Order)
	orderIds := make([]string, 0)
	for _, o := range orders {
		orderIds = append(orderIds, o.StringID())
	}

	rq := gmgo.Q{
		"archived": false,
		"currentOrderId": gmgo.Q{
			"$in": orderIds,
		},
	}

	rxResults, err := session.FindAll(rq, new(data.Prescription))
	if err != nil {
		return util.EmptyRxList(), err
	}

	rxs := rxResults.([]*data.Prescription)
	for _, rx := range rxs {
		err := LoadRxOrder(rx, session)
		if err != nil {
			log.Printf("Error loading order for rx %s, %s\n", rx.StringID(), err)
		}
	}
	return rxs, nil
}

// TriggerReverseInsuranceOrCash - trigger reverse insurance for the prescription
func TriggerReverseInsuranceOrCash(rx *data.Prescription) error {
	session := data.Session()
	defer session.Close()

	o := rx.CurrentOrder
	if o.PaymentOption == data.PayThroughInsurance &&
		o.InsuranceMilestone != nil && o.InsuranceMilestone.Status == data.MSCompleted {

			o.InsuranceMilestone.Status = data.MSReverseInsurance
	}
	if o.PaymentOption == data.PayDirectly &&
		o.PaymentMilestone != nil && o.PaymentMilestone.Status != data.MSNotInitiated {

			o.PaymentMilestone.ShouldReverseCashClaim = true
	}

	order.UpdateInSession(o, session)
	return nil
}

// LoadRxOrder loads order for a prescription data
func LoadRxOrder(rx *data.Prescription, session *gmgo.DbSession) error {
	errList := make([]error, 0)

	//get order data by reference
	o, err := data.OrderByRef(rx.CurrentOrderRef, session)
	if err != nil && err.Error() != data.ErrNotFound {
		errList = append(errList, err)
	} else {
		rx.CurrentOrder = o
	}

	return util.HandleRefLoadError("ERRORS", errList)
}

// LoadRxReferences loads all the references for a prescription data and its references
// returning complete data.
func LoadRxReferences(rx *data.Prescription, session *gmgo.DbSession, loadOrderRefs bool) error {
	errList := make([]error, 0)

	//fetch manager
	mgr, err := data.ManagerByRef(rx.ManagerRef, session)
	if err != nil && err.Error() != data.ErrNotFound {
		errList = append(errList, err)
	} else {
		err = user.LoadAllUserReferences(mgr, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			rx.Manager = mgr
		}
	}

	//fetch patient
	if rx.PatientRef != nil {
		pt, err := data.PatientByRef(rx.PatientRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			rx.Patient = pt
		}
	}
	//get doctor by reference
	if rx.DoctorRef != nil {
		d, err := data.DoctorByRef(rx.DoctorRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			rx.Doctor = d
		}
	}
	//get insurance by refernce
	if rx.InsuranceRef != nil {
		i, err := data.InsuranceByRef(rx.InsuranceRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			rx.Insurance = i
		}
	}
	//get order data by reference
	o, err := data.OrderByRef(rx.CurrentOrderRef, session)
	if err != nil && err.Error() != data.ErrNotFound {
		errList = append(errList, err)
	} else {
		rx.CurrentOrder = o
	}

	if loadOrderRefs {
		err := order.LoadReferences(rx.CurrentOrder, rx.Manager, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		}
	}
	if rx.Patient != nil {
		rx.Patient.Manager = rx.Manager
	}
	if rx.Insurance != nil {
		rx.Insurance.Manager = rx.Manager
	}

	return util.HandleRefLoadError("ERRORS", errList)
}

// LoadRxReferencesWithConfig loads all the references for a prescription data and its references
// returning complete data.
func LoadRxReferencesWithConfig(rx *data.Prescription, session *gmgo.DbSession, config *RxLoadConfig) error {
	errList := make([]error, 0)

	if config.LoadManager && rx.Manager == nil {
		//fetch manager
		mgr, err := data.ManagerByRef(rx.ManagerRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			err = user.LoadAllUserReferences(mgr, session)
			if err != nil && err.Error() != data.ErrNotFound {
				errList = append(errList, err)
			} else {
				rx.Manager = mgr
			}
		}
	}

	if config.LoadPatient && rx.Patient == nil {
		//fetch patient
		if rx.PatientRef != nil {
			pt, err := data.PatientByRef(rx.PatientRef, session)
			if err != nil && err.Error() != data.ErrNotFound {
				errList = append(errList, err)
			} else {
				rx.Patient = pt
			}
		}
	}
	if config.LoadDoctor && rx.Doctor == nil {
		//get doctor by reference
		if rx.DoctorRef != nil {
			d, err := data.DoctorByRef(rx.DoctorRef, session)
			if err != nil && err.Error() != data.ErrNotFound {
				errList = append(errList, err)
			} else {
				rx.Doctor = d
			}
		}
	}

	if config.LoadInsurance && rx.Insurance == nil {
		//get insurance by reference
		if rx.InsuranceRef != nil {
			i, err := data.InsuranceByRef(rx.InsuranceRef, session)
			if err != nil && err.Error() != data.ErrNotFound {
				errList = append(errList, err)
			} else {
				rx.Insurance = i
			}
		}
	}

	if config.LoadOrder && rx.CurrentOrder == nil {
		//get order data by reference
		o, err := data.OrderByRef(rx.CurrentOrderRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			rx.CurrentOrder = o
		}
	}

	if rx.CurrentOrder != nil && config.OrderLoadConfig != nil {
		err := order.LoadReferencesWithConfig(rx.CurrentOrder, rx.Manager, config.OrderLoadConfig, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		}
	}

	if rx.Patient != nil {
		rx.Patient.Manager = rx.Manager
	}

	if rx.Insurance != nil {
		rx.Insurance.Manager = rx.Manager
	}

	return util.HandleRefLoadError("ERRORS", errList)
}

// Update update prescription
func Update(ctx context.Context, rx *data.Prescription) error {
	session := data.Session()
	defer session.Close()

	rx.UpdatedDate = util.NowUTC()
	return UpdateInSession(rx, session)
}

// AddInsurance assigns data.Insurance to prescription
func AddInsurance(ctx context.Context, rx *data.Prescription, i *data.Insurance) error {
	session := data.Session()
	defer session.Close()

	//Update insurance reference
	rx.InsuranceRef = data.DBRef(i.CollectionName(), i.ID)
	rx.Insurance = i
	rx.UpdatedDate = util.NowUTC()

	o := rx.CurrentOrder

	// remove missing status
	if o.MissingInfoState != nil {
		missingStatus := *o.MissingInfoState.MissingStatus
		//remove any missing insurance in missing info state
		if missingStatus == data.MissingInsurance || missingStatus == data.InsuranceExpired || missingStatus == data.InvalidInsurance {
			o.MissingInfoState = nil
			o.RemoveSleepStatus(data.CallPatient)
		}
	}

	//clear any payment milestone actions taken
	if o.PaymentMilestone != nil && o.PaymentMilestone.Status == data.MSCashPriceCheck {
		o.PaymentMilestone.Status = data.MSNotInitiated
	}

	o.PaymentOption = data.PayThroughInsurance

	// handle first fill
	if o.IsFirstFill() {
		o.MoveToScrubAndRoute()
	}

	// clear pending status
	if rx.Pending {
		rx.Pending = false
		rx.PendingSinceDate = nil
	}

	rx.CurrentOrder = o
	err := UpdateInSession(rx, session)
	if err != nil {
		log.Printf("[ERROR] Updating updatedDate for rx with id %s", rx.StringID())
		return err
	}

	order.Update(o)
	return nil
}

// AddInsuranceToPrescription assigns model.PrescriptionInsurance to prescription
func AddInsuranceToPrescription(ctx context.Context, insuranceID, rxID string) (*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	//Find Rx
	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return nil, err
	}

	//Find Insurance
	i, err := insurance.FindInsurance(ctx, insuranceID)
	if err != nil {
		return nil, err
	}

	err = AddInsurance(ctx, rx, i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

// AddPaymentToPrescription assigns payment to prescription
func AddPaymentToPrescription(ctx context.Context, paymentID, rxID string) (*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	//Find Rx
	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return nil, err
	}

	//Find Payment
	p, err := FindPaymentInfo(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	//Find Order
	o, err := order.FindByID(ctx, rx.CurrentOrderID)
	if err != nil {
		return nil, err
	}

	//Update payment reference
	o.PaymentRef = data.DBRef(p.CollectionName(), p.ID)
	if o.MissingInfoState != nil && *o.MissingInfoState.MissingStatus == data.PaymentError {
		o.MissingInfoState = nil
		o.RemoveSleepStatus(data.CallPatient)
	}
	err = order.UpdateInSession(o, session)
	if err != nil {
		return nil, err
	}

	rx.UpdatedDate = util.NowUTC()

	if rx.Pending {
		rx.Pending = true
		rx.PendingSinceDate = nil
	}

	err = UpdateInSession(rx, session)
	if err != nil {
		log.Printf("[ERROR] updating updatedDate for rx with id %s", rx.StringID())
	}

	return p, nil
}

// AddAddressToPrescription assigns address to a prescription
func AddAddressToPrescription(ctx context.Context, addressID, rxID string) (*data.Address, error) {
	session := data.Session()
	defer session.Close()

	//Find Rx
	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return nil, err
	}

	//Find Address
	addr, err := address.FindAddress(ctx, addressID)
	if err != nil {
		return nil, err
	}

	//Find Order
	o, err := order.FindByID(ctx, rx.CurrentOrderID)
	if err != nil {
		return nil, err
	}
	o.Address = addr
	o.AddressRef = data.DBRef(addr.CollectionName(), addr.ID)

	o.UpdatedDate = util.NowUTC()
	err = session.Update(gmgo.Q{"_id": o.ID}, o)
	if err != nil {
		return nil, err
	}

	rx.UpdatedDate = util.NowUTC()
	err = session.Update(gmgo.Q{"_id": rx.ID}, rx)
	if err != nil {
		log.Printf("[ERROR] updating updatedDate for rx with id %s", rx.StringID())
	}

	return addr, nil
}

// AddDoctorToPrescription assigns doctor to a prescription
func AddDoctorToPrescription(ctx context.Context, doctorID, rxID string) (*data.Doctor, error) {
	session := data.Session()
	defer session.Close()

	//Find Rx
	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return nil, err
	}

	//Find Doctor
	dctr, err := doctor.FindByIDInSession(ctx, doctorID, session)
	if err != nil {
		return nil, err
	}

	rx.AddDoctor(dctr)

	rx.UpdatedDate = util.NowUTC()
	err = UpdateInSession(rx, session)
	if err != nil {
		log.Printf("[ERROR] updating updatedDate for rx with id %s", rx.StringID())
	}

	return dctr, nil
}

// MarkPrescriptionAsPending - set prescription to pending status
func MarkPrescriptionAsPending(ctx context.Context, rxID string) error {
	session := data.Session()
	defer session.Close()

	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return err
	}

	rx.UpdatedDate = util.NowUTC()
	rx.Pending = true
	rx.PendingSinceDate = util.NowUTC()

	o := rx.CurrentOrder
	order.UpdateInSession(o, session)

	return UpdateInSession(rx, session)
}

// SkipRefill skip refill for the given prescription id
func SkipRefill(ctx context.Context, rxID string) error {
	session := data.Session()
	defer session.Close()

	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return err
	}
	if rx.ReFillDate == nil {
		return errors.New("Prescription already in processing state")
	}

	nextRefillDate := rx.ReFillDate.Add(time.Hour * 24 * time.Duration(rx.DaysOfSupply))
	return UpdateRefillDateForPrescription(ctx, rx, nextRefillDate, session)
}

// UpdateRefillDate updates the refill date for given prescription
func UpdateRefillDate(ctx context.Context, userID, rxID string, refillDate time.Time) error {
	session := data.Session()
	defer session.Close()

	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return err
	}
	return UpdateRefillDateForPrescription(ctx, rx, refillDate, session)
}

// UpdateRefillDateForPrescription updates the refill date for given prescription in given session
func UpdateRefillDateForPrescription(ctx context.Context, rx *data.Prescription, nextRefillDate time.Time, session *gmgo.DbSession) error {

	if nextRefillDate.Weekday() == time.Sunday || nextRefillDate.Weekday() == time.Saturday {
		nextRefillDate = nextRefillDate.Add(time.Hour * 24 * time.Duration(-2))
	}

	rx.ReFillDate = &nextRefillDate
	rx.ScheduleDate = &nextRefillDate
	rx.UpdatedDate = util.NowUTC()

	err := UpdateInSession(rx, session)
	if err != nil {
		return err
	}

	o := rx.CurrentOrder
	if o.DeliveryMilestone.Parameters == nil {
		o.DeliveryMilestone.Parameters = make(map[string]interface{})
	}

	o.DeliveryMilestone.Parameters["nextFillDate"] = nextRefillDate.Unix()
	err = order.UpdateInSession(o, session)
	if err != nil {
		return err
	}

	return nil
}

// UpdateFollowUpDate updates the follow up date for a prescription
func UpdateFollowUpDate(ctx context.Context, userID, rxID string, followUpDate *time.Time) error {
	session := data.Session()
	defer session.Close()

	rx, err := FindByIDInSession(rxID, session)
	if err != nil {
		return err
	}
	rx.ScheduleDate = followUpDate
	err = UpdateInSession(rx, session)
	if err != nil {
		return err
	}
	return nil
}

// TriggerRefillAuthorization triggers refill authorization for given prescription
func TriggerRefillAuthorization(ctx context.Context, rx *data.Prescription) error {
	session := data.Session()
	defer session.Close()

	o := rx.CurrentOrder

	refillAuth := new(data.RefillMilestone)
	refillAuth.Status = data.MSNotInitiated
	refillAuth.RxName = rx.Name
	refillAuth.Quantity = rx.MedicationQuantity
	refillAuth.DaysOfSupply = rx.DaysOfSupply
	refillAuth.LastProcessedDate = util.NowUTC()

	o.RefillMilestone = refillAuth
	rx.ScheduleDate = util.NowUTC()

	err := UpdateInSession(rx, session)
	if err != nil {
		log.Printf("[ERROR] Error saving prescription data %s\n", err)
		return err
	}

	err = order.UpdateInSession(o, session)
	if err != nil {
		log.Printf("[ERROR] Error updating order data %s\n", err)
		return err
	}

	return nil
}

// FindAllPaymentInfo returns all payment info
func FindAllPaymentInfo(ctx context.Context) ([]*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Payment))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Payment), nil
}

// FindPaymentInfo
func FindPaymentInfo(ctx context.Context, paymentID string) (*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	result := new(data.Payment)
	err := session.FindByID(paymentID, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindAllPaymentInfoForPatient finds all payment info for patient.
func FindAllPaymentInfoForPatient(ctx context.Context, patientID string) ([]*data.Payment, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{"patientId": patientID}, new(data.Payment))
	if err != nil {
		return nil, err
	}

	return result.([]*data.Payment), nil
}

//SaveDoctor - save doctor
func SaveDoctor(doctor *data.Doctor) (*data.Doctor, error) {
	session := data.Session()
	defer session.Close()

	t := util.NowUTC()
	doctor.CreatedDate = t
	doctor.UpdatedDate = t
	doctor.PhoneNumber = util.SanitizePhoneNumber(doctor.PhoneNumber)
	err := session.Save(doctor)
	if err != nil {
		return nil, err
	}
	return doctor, nil
}

// RevertReverseInsurance - cancels reverse insurance action
func RevertReverseInsurance(rx *data.Prescription) {
	o := rx.CurrentOrder
	if o.InsuranceMilestone != nil && o.InsuranceMilestone.Status == data.MSReverseInsurance {
		o.InsuranceMilestone.Status = data.MSCompleted
	}
}
