package ror

import (
	"context"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/service/prescription"
	"github.com/phil-inc/plib/core/util"
)

//RecordRoREvent - records the Run out risk event
func RecordRoREvent(ctx context.Context, rxID string, eventType string) {
	rorEvent := new(data.RorEvent)
	rorEvent.EventType = eventType
	rorEvent.TimeStamp = util.NowUTC()
	rorEvent.RxID = rxID

	SaveRorEvent(ctx, rorEvent)
}

//RecordRorEventForUser record ROR event for user
func RecordRorEventForUser(ctx context.Context, user *data.User, eventType string) error {
	session := data.Session()
	defer session.Close()

	event := new(data.RorEvent)
	event.InitData()
	event.ManagerID = user.StringID()
	event.TimeStamp = util.NowUTC()
	event.EventType = eventType

	err := session.Save(event)
	if err != nil {
		log.Printf("[ERROR] saving event. Error: %s", err)
	}

	return saveRorEventPostgres(event)
}

//SaveRorEvent - save the Run out risk event
func SaveRorEvent(ctx context.Context, rorEvent *data.RorEvent) error {
	session := data.Session()
	defer session.Close()

	rxID := rorEvent.RxID
	if rxID == "" {
		return errors.New("RxId is empty")
	}
	rx, err := prescription.FindByIDInSession(rxID, session)
	if err != nil {
		return err
	}

	event := new(data.RorEvent)
	event.InitData()
	event.RxID = rorEvent.RxID
	event.TimeStamp = rorEvent.TimeStamp
	event.EventType = rorEvent.EventType
	event.ManagerID = data.RefID(rx.ManagerRef)
	event.PatientID = data.RefID(rx.PatientRef)

	if rx.CurrentOrder != nil {
		event.OrderID = rx.CurrentOrder.ID.Hex()
		event.OrderNumber = rx.CurrentOrder.OrderNumber
	}

	err = session.Save(event)
	if err != nil {
		log.Printf("[ERROR] saving event. Error: %s", err)
	}

	return saveRorEventPostgres(event)
}

//RecordRoREventForRx - records the Run out risk event for given prescription
func RecordRoREventForRx(ctx context.Context, rx *data.Prescription, eventType string) error {
	session := data.Session()
	defer session.Close()

	event := new(data.RorEvent)
	event.InitData()
	event.ManagerID = data.RefID(rx.ManagerRef)
	event.PatientID = data.RefID(rx.PatientRef)
	event.RxID = rx.StringID()
	event.TimeStamp = util.NowUTC()
	event.EventType = eventType

	if rx.CurrentOrder != nil {
		event.OrderID = rx.CurrentOrder.ID.Hex()
		event.OrderNumber = rx.CurrentOrder.OrderNumber
	}

	err := session.Save(event)
	if err != nil {
		log.Printf("[ERROR] saving event. Error: %s", err)
	}

	return saveRorEventPostgres(event)
}

//FindAllRorEvents - finds all ror events
func FindAllRorEvents(ctx context.Context) ([]*data.RorEvent, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.RorEvent))
	if err != nil {
		return nil, err
	}
	return results.([]*data.RorEvent), nil
}

func saveRorEventPostgres(event *data.RorEvent) error {
	if util.IsLocal() {
		return nil
	}

	tx, err := data.PostgresDb.Begin()
	if err != nil {
		return err
	}

	statement, err := tx.Prepare(pq.CopyIn("ror_event", "id", "rxid", "managerid", "patientid", "timestamp", "eventtype", "ordernumber", "orderid", "createddate", "updateddate"))
	defer statement.Close()

	if err != nil {
		return err
	}

	_, err = statement.Exec(event.StringID(), event.RxID, event.ManagerID, event.PatientID, event.TimeStamp, event.EventType, event.OrderNumber, event.OrderID, event.CreatedDate, event.UpdatedDate)
	if err != nil {
		return err
	}

	_, err = statement.Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}
