package address

import (
	"context"
	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/util"
	"time"
)

// FindAddress find all the orders for prescription
func FindAddress(ctx context.Context, addressID string) (*data.Address, error) {
	session := data.Session()
	defer session.Close()

	result := new(data.Address)
	err := session.FindByID(addressID, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateInSession updates the given order using passed in session
func UpdateInSession(a *data.Address, session *gmgo.DbSession) error {
	a.UpdatedDate = util.NowUTC()
	return session.Update(gmgo.Q{"_id": a.ID}, a)
}

// Update save or update the given order
func Update(a *data.Address) error {
	session := data.Session()
	defer session.Close()

	return UpdateInSession(a, session)
}

// SaveInSession saves the order in a given session. Returns error if any
func SaveInSession(a *data.Address, session *gmgo.DbSession) error {
	t := time.Now().UTC()
	if a.ID.Hex() == "" {
		a.CreatedDate = &t
	}
	a.UpdatedDate = &t
	return session.Save(a)
}

// FindAllForPatient finds all addresses of patient.
func FindAllForPatient(ctx context.Context, patientID string) ([]*data.Address, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{"patientId": patientID}, new(data.Address))
	if err != nil {
		return nil, err
	}

	return result.([]*data.Address), nil
}
