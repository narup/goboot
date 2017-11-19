package data

import "time"

//Shipment data representation for shipment
type Shipment struct {
	BaseData              `bson:",inline"`
	UserID                string     `json:"userId" bson:"userId"` //Indexed
	Orders                []string   `json:"orders" bson:"orders" pson:"orders"`
	TrackingNumber        string     `json:"trackingNumber" bson:"trackingNumber,omitempty" pson:"tracking_number"` //Indexed
	TrackingURL           string     `json:"trackingUrl" bson:"trackingUrl,omitempty"`
	TrackingURLForMobile  string     `json:"trackingUrlForMobile" bson:"trackingUrlForMobile,omitempty"`
	ShippingLabelURL      string     `json:"shippingLabelUrl" bson:"shippingLabelUrl,omitempty"`
	ShippingServiceName   string     `json:"shippingServiceName" bson:"shippingServiceName,omitempty" pson:"shipping_service_name"`
	ShippingLabelCost     string     `json:"shippingLabelCost" bson:"shippingLabelCost,omitempty" pson:"shipping_label_cost"`
	PharmacyID            string     `json:"pharmacyId" bson:"pharmacyId,omitempty" pson:"pharmacy_id"`
	ShipmentType          string     `json:"shipmentType" bson:"shipmentType,omitempty" pson:"shipment_type"` //BUNDLE, SINGLE
	ShippedDate           *time.Time `json:"shippedDate" bson:"shippedDate,omitempty" pson:"shipped_date"`    //Sparse Indexed
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate" bson:"estimatedDeliveryDate,omitempty" pson:"estimated_delivery_date"`
	EstimatedDeliveryDays int        `json:"estimatedDeliveryDays" bson:"estimatedDeliveryDays,omitempty" pson:"estimated_delivery_days"`
}

// CollectionName function from gmgo.Document interface
func (s Shipment) CollectionName() string {
	return "shipment"
}

//ShipmentErrorLog stores all the shipment with error or unknown status
type ShipmentErrorLog struct {
	BaseData       `bson:",inline"`
	OrderNumber    string `json:"orderNumber" bson:"orderNumber,omitempty"`
	ShipmentCount  int    `json:"shipmentCount" bson:"shipmentCount"`
	TrackingNumber string `json:"trackingNumber" bson:"trackingNumber,omitempty"`
	PharmacyID     string `json:"pharmacyId" bson:"pharmacyId,omitempty"`
	PharmacyName   string `json:"pharmacyName" bson:"pharmacyName,omitempty"`
	CheckCount     int    `json:"checkCount" bson:"checkCount"`
	ShipmentType   string `json:"shipmentType" bson:"shipmentType,omitempty"` //BUNDLE, SINGLE
}

// CollectionName function from gmgo.Document interface
func (se ShipmentErrorLog) CollectionName() string {
	return "shipmentErrorLog"
}

//DeliveryExceptionLog stores all the deivery exceptions for analytics
type DeliveryExceptionLog struct {
	BaseData         `bson:",inline"`
	OrderNumbers     []string `json:"orderNumbers" bson:"orderNumbers,omitempty" pson:"order_numbers"`
	ShipmentID       string   `json:"shipmentId" bson:"shipmentId,omitempty" pson:"shipment_id"`
	TrackingNumber   string   `json:"trackingNumber" bson:"trackingNumber,omitempty" pson:"tracking_number"`
	ExceptionDetails string   `json:"exceptionDetails" bson:"exceptionDetails,omitempty" pson:"exception_details"`
	ExceptionType    string   `json:"exceptionType" bson:"exceptionType,omitempty" pson:"exception_type"`
}

// CollectionName function from gmgo.Document interface
func (se DeliveryExceptionLog) CollectionName() string {
	return "deliveryExceptionLog"
}

//DeliveryPickupToken stores the delivery pickup token for same day delivery
type DeliveryPickupToken struct {
	BaseData       `bson:",inline"`
	UserID         string          `json:"userId" bson:"userId"`
	OrderIDs       []string        `json:"orderId" bson:"orderIds"`
	Token          string          `json:"token" bson:"token"`
	ShipmentID     string          `json:"shipmentId" bson:"shipmentId"`
	TokenActive    bool            `json:"tokenActive" bson:"tokenActive"`
	DeliveryWindow *DeliveryWindow `json:"deliveryWindow" bson:"deliveryWindow,omitempty"`
	DeliveryID     string          `json:"deliveryId" bson:"deliveryId,omitempty"`
	ReadyBy        *time.Time      `json:"readyBy" bson:"readyBy,omitempty"`
	CreatedBy      string          `json:"createdBy" bson:"createdBy,omitempty"`
}

// CollectionName function from gmgo.Document interface
func (d DeliveryPickupToken) CollectionName() string {
	return "deliveryPickupToken"
}

//DeliveryWindow stores the delivery pickup window for same day delivery
type DeliveryWindow struct {
	ID       string     `json:"id"`
	StartsAt *time.Time `json:"startsAt,omitempty"`
	EndsAt   *time.Time `json:"endsAt,omitempty"`
}
