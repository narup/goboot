package order

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/util"
)

var digits = "1234567890"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// OrderLoadConfig defines configuration to load order data
type OrderLoadConfig struct {
	LoadAddress,
	LoadPayment,
	LoadOriginPharmacy,
	LoadPartnerPharmacy bool
}

//UpdateInSession updates the given order using passed in session
func UpdateInSession(o *data.Order, session *gmgo.DbSession) error {
	o.UpdatedDate = util.NowUTC()
	return session.Update(gmgo.Q{"_id": o.ID}, o)
}

// Update save or update the given order
func Update(o *data.Order) error {
	session := data.Session()
	defer session.Close()

	t := time.Now().UTC()
	o.UpdatedDate = &t
	return session.Update(gmgo.Q{"_id": o.ID}, o)
}

// Save saves the order in a given session. Returns error if any
func Save(session *gmgo.DbSession, o *data.Order) error {
	t := time.Now().UTC()
	if o.ID.Hex() == "" {
		o.CreatedDate = &t
	}
	o.UpdatedDate = &t
	return session.Save(o)
}

// SaveInSession saves the given order using passed in session
func SaveInSession(o *data.Order, session *gmgo.DbSession) error {
	o.CreatedDate = util.NowUTC()
	o.UpdatedDate = util.NowUTC()

	return session.Save(o)
}

// FindByID finds order data and all its references by id. It returns
// complete data with all the DBRefs loaded and populated.
func FindByID(ctx context.Context, orderID string) (*data.Order, error) {
	session := data.Session()
	defer session.Close()

	return FindByIDInSession(orderID, session)
}

//FindByIDInSession finds order based on given order ID and db session
func FindByIDInSession(orderID string, session *gmgo.DbSession) (*data.Order, error) {
	o := new(data.Order)
	err := session.FindByID(orderID, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// FindOrdersForPrescription find all the orders for prescription
func FindOrdersForPrescription(ctx context.Context, rxID string) ([]*data.Order, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"rxId": rxID}, new(data.Order))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Order), nil
}

// FindByOrderNumber finds the order by order number. References are not loaded.
func FindByOrderNumber(ctx context.Context, orderNumber string) (*data.Order, error) {
	session := data.Session()
	defer session.Close()

	return FindByOrderNumberInSession(orderNumber, session)
}

// FindByOrderNumberInSession finds the order by order number in session. References are not loaded.
func FindByOrderNumberInSession(orderNumber string, session *gmgo.DbSession) (*data.Order, error) {
	o := new(data.Order)
	err := session.Find(gmgo.Q{"orderNumber": orderNumber}, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// UpdateOrderMilestoneStatus updates order
func UpdateOrderMilestoneStatus(ctx context.Context, orderNumber, milestoneName, status string) error {
	found := false
	for _, mn := range data.MilestoneNames {
		if mn == milestoneName {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Invalid milestone name")
	}
	found = false
	for _, st := range data.MilestoneStatusList {
		if st == status {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Invalid milestone status value")
	}

	session := data.Session()
	defer session.Close()

	key := fmt.Sprintf("%s.status", milestoneName)
	return session.UpdateFieldValue(gmgo.Q{"orderNumber": orderNumber}, new(data.Order).CollectionName(), key, status)
}

// UpdatePaymentStatusForOrder for given order update payment status
func UpdatePaymentStatusForOrder(o *data.Order, session *gmgo.DbSession) {
	if !o.IsPendingPaymentApproval() && !o.IsPendingPaymentApprovalAndSignup() {
		return
	}

	pm := o.PaymentMilestone
	pm.Status = data.MSPaymentAuthorized
	pm.LastProcessedDate = util.NowUTC()

	mis := o.MissingInfoState
	if mis != nil {
		status := *mis.MissingStatus
		if status == data.PaymentError ||
			status == data.CopayIncreased ||
			status == data.CopayDecreased ||
			status == data.OriginalCopayNotAvailable {

			o.MissingInfoState = nil
		}
	}

	//update the order
	o.RemoveSleepStatus(data.PaymentApproval)
	o.RemoveSleepStatus(data.PaymentApprovalUntilFreeTrialDelivered)
	UpdateInSession(o, session)
}

// FindShipmentByTrackingNumber find the shipment for given tracking number.
func FindShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*data.Shipment, error) {
	session := data.Session()
	defer session.Close()

	s := new(data.Shipment)
	return s, session.Find(gmgo.Q{"trackingNumber": trackingNumber}, s)
}

//GetShipmentsCountForUser returns the shipment count for the given user ID
func GetShipmentsCountForUser(ctx context.Context, userID string) (int, error) {
	session := data.Session()
	defer session.Close()

	return GetShipmentsCountForUserInSession(ctx, userID, session)
}

//GetShipmentsCountForUserInSession returns the shipment count for the given user ID in session
func GetShipmentsCountForUserInSession(ctx context.Context, userID string, session *gmgo.DbSession) (int, error) {
	result, err := session.FindAllWithFields(gmgo.Q{"userId": userID}, []string{"_id"}, new(data.Shipment))
	if err != nil {
		return -1, err
	}

	ss := result.([]*data.Shipment)
	return len(ss), nil
}

// FindShipmentByID find the shipment for given tracking number.
func FindShipmentByID(ctx context.Context, shipmentID string) (*data.Shipment, error) {
	session := data.Session()
	defer session.Close()

	result := new(data.Shipment)
	err := session.FindByID(shipmentID, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindAllLastDayShipments return all the shipments for last day
func FindAllLastDayShipments(ctx context.Context) ([]*data.Shipment, error) {
	session := data.Session()
	defer session.Close()

	st, _ := util.YesterdayStartEndTimePST()
	q := gmgo.Q{
		"shippedDate": gmgo.Q{
			"$gte": st.UTC(),
		},
	}
	results, err := session.FindAll(q, new(data.Shipment))
	if err != nil {
		return nil, err
	}

	return results.([]*data.Shipment), nil
}

//SaveFaxContent saves the fax content for the fax type and order number
func SaveFaxContent(ctx context.Context, faxContent, faxType, orderNumber string) error {
	session := data.Session()
	defer session.Close()

	fc := new(data.FaxContent)
	fc.InitData()
	fc.Content = faxContent
	fc.FaxType = faxType
	fc.OrderNumber = orderNumber

	return session.Save(fc)
}

// GenerateOrderNumber generates the random order number of format xxxx-xxxx-xxxx
func GenerateOrderNumber() string {
	session := data.Session()
	defer session.Close()

	orderNumber := generateOrderNumber()
	o, _ := FindByOrderNumberInSession(orderNumber, session)
	for o != nil {
		orderNumber = generateOrderNumber()
		o, _ = FindByOrderNumberInSession(orderNumber, session)
	}
	return orderNumber
}

func generateOrderNumber() string {
	size := 14
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		if i == 4 || i == 9 {
			buf[i] = '-'
		} else {
			buf[i] = digits[rand.Intn(len(digits))]
		}
	}
	return string(buf)
}

// InitializeMilestones sets up the order milestones
func InitializeMilestones(originType string, o *data.Order) {
	if originType == "Transfer" {
		tm := new(data.TransferMilestone)
		tm.Status = data.MSNotInitiated
		o.TransferMilestone = tm
	} else {
		nm := new(data.NewRxMilestone)
		nm.Status = data.MSNotInitiated
		o.NewRxMilestone = nm
	}

	im := new(data.InsuranceMilestone)
	im.Status = data.MSNotInitiated
	im.SendCopayConfirmEmail = true
	o.InsuranceMilestone = im

	pm := new(data.PaymentMilestone)
	pm.Status = data.MSNotInitiated
	pm.AppliedDummyCharge = false
	o.PaymentMilestone = pm

	sm := new(data.StockCheckMilestone)
	sm.Status = data.MSNotInitiated
	o.StockCheckMilestone = sm

	dm := new(data.DeliveryMilestone)
	dm.Status = data.MSNotInitiated
	dm.DeliveryConfirmed = false
	o.DeliveryMilestone = dm
}

// LoadReferences loads all the order data references such as address, payment, pharmacy
func LoadReferences(o *data.Order, manager *data.User, session *gmgo.DbSession) error {
	errList := make([]error, 0)
	//fetch address
	if o.AddressRef != nil {
		addr, err := data.AddressByRef(o.AddressRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.Address = addr
		}
	}
	//fetch payment
	if o.PaymentRef != nil {
		p, err := data.PaymentByRef(o.PaymentRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.Payment = p
		}
	}
	//fetch origin pharmacy
	if o.OriginPharmacyRef != nil {
		op, err := data.PharmacyByRef(o.OriginPharmacyRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.OriginPharmacy = op
		}
	}
	//fetch partner pharmacy
	if o.PartnerPharmacyRef != nil {
		pp, err := data.PharmacyByRef(o.PartnerPharmacyRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.PartnerPharmacy = pp
		}
	}
	if o.Payment != nil {
		o.Payment.Manager = manager
	}
	if o.Address != nil {
		o.Address.Manager = manager
	}
	return util.HandleRefLoadError("ERRORS", errList)
}

// LoadReferencesWithConfig loads all the order data references such as address, payment, pharmacy with config
func LoadReferencesWithConfig(o *data.Order, manager *data.User, config *OrderLoadConfig, session *gmgo.DbSession) error {
	errList := make([]error, 0)
	//fetch address
	if config.LoadAddress && o.AddressRef != nil && o.Address == nil {
		addr, err := data.AddressByRef(o.AddressRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.Address = addr
		}
	}
	//fetch payment
	if config.LoadPayment && o.PaymentRef != nil && o.Payment == nil {
		p, err := data.PaymentByRef(o.PaymentRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.Payment = p
		}
	}
	//fetch origin pharmacy
	if config.LoadOriginPharmacy && o.OriginPharmacyRef != nil && o.OriginPharmacy == nil {
		op, err := data.PharmacyByRef(o.OriginPharmacyRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.OriginPharmacy = op
		}
	}
	//fetch partner pharmacy
	if config.LoadPartnerPharmacy && o.PartnerPharmacyRef != nil && o.PartnerPharmacy == nil {
		pp, err := data.PharmacyByRef(o.PartnerPharmacyRef, session)
		if err != nil && err.Error() != data.ErrNotFound {
			errList = append(errList, err)
		} else {
			o.PartnerPharmacy = pp
		}
	}
	if o.Payment != nil {
		o.Payment.Manager = manager
	}
	if o.Address != nil {
		o.Address.Manager = manager
	}
	return util.HandleRefLoadError("ERRORS", errList)
}

//FindAllOrders returns all orders
func FindAllOrders(ctx context.Context) ([]*data.Order, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Order))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Order), nil
}
