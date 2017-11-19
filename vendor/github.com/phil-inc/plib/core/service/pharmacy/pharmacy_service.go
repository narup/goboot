package pharmacy

import (
	"context"
	"time"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/util"
)

// FindByID finds pharmacy by ID
func FindByID(ctx context.Context, pharmacyID string) (*data.Pharmacy, error) {
	session := data.Session()
	defer session.Close()

	p := new(data.Pharmacy)
	return p, session.FindByID(pharmacyID, p)
}

// FindPartnersByState find all partner pharmacies by state
func FindPartnersByState(state string) ([]*data.Pharmacy, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{"partner": true, "state": state}, new(data.Pharmacy))
	if err != nil {
		return nil, err
	}
	return result.([]*data.Pharmacy), nil
}

// FindByNameAndPhone find pharmacy by name and phone number
func FindByNameAndPhone(name, phone string) (*data.Pharmacy, error) {
	session := data.Session()
	defer session.Close()

	ph := new(data.Pharmacy)
	return ph, session.Find(gmgo.Q{"name": name, "phoneNumber": phone}, ph)
}

// FindAll returns all pharmacies
func FindAll(ctx context.Context) ([]*data.Pharmacy, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{}, new(data.Pharmacy))
	if err != nil {
		return nil, err
	}
	return results.([]*data.Pharmacy), nil
}

//FindPharmacyUserByUserID find Pharmacy user for the given Phil user ID
func FindPharmacyUserByUserID(ctx context.Context, userID string) (*data.PharmacyUser, error) {
	session := data.Session()
	defer session.Close()

	q := gmgo.Q{
		"rexUser.$id": data.ObjectID(userID),
	}

	pu := new(data.PharmacyUser)
	err := session.Find(q, pu)
	if err != nil {
		return nil, err
	}

	ph := new(data.Pharmacy)
	err = session.FindByRef(pu.PharmacyRef, ph)
	if err != nil {
		return nil, err
	}

	pu.Pharmacy = ph

	return pu, nil
}

// FindAllPharmacyUsers finds all the pharmacy users for the given pharmacy ID
func FindAllPharmacyUsers(ctx context.Context, pharmacyID string) ([]*data.User, error) {
	session := data.Session()
	defer session.Close()

	q := gmgo.Q{
		"pharmacy.$id": data.ObjectID(pharmacyID),
	}

	results, err := session.FindAll(q, new(data.PharmacyUser))
	if err != nil {
		return nil, err
	}

	finalList := make([]*data.User, 0)
	pUsers := results.([]*data.PharmacyUser)
	for _, pUser := range pUsers {
		usr := new(data.User)
		err := session.FindByRef(pUser.UserRef, usr)
		if err != nil {
			//ignore
		} else {
			finalList = append(finalList, usr)
		}
	}

	return finalList, nil
}

// DefaultPartnerPharmacy returns the default partner pharmacy
func DefaultPartnerPharmacy() *data.Pharmacy {
	session := data.Session()
	defer session.Close()

	pharmacy := new(data.Pharmacy)
	err := session.Find(gmgo.Q{"name": "Phil", "phoneNumber": "8559770975"}, pharmacy)
	if err != nil {
		pharmacy.InitData()
		pharmacy.Name = "Phil"
		pharmacy.PhoneNumber = "8559770975"
		if err := session.Save(pharmacy); err != nil {
			return pharmacy
		}
	}
	return pharmacy
}

// SaveInSession saves the order in a given session. Returns error if any
func SaveInSession(pharmacy *data.Pharmacy, session *gmgo.DbSession) error {
	t := time.Now().UTC()
	if pharmacy.ID.Hex() == "" {
		pharmacy.CreatedDate = &t
	}
	pharmacy.UpdatedDate = &t
	return session.Save(pharmacy)
}

// Save saves the order in a given session. Returns error if any
func Save(pharmacy *data.Pharmacy) error {
	session := data.Session()
	defer session.Close()

	t := time.Now().UTC()
	if pharmacy.ID.Hex() == "" {
		pharmacy.CreatedDate = &t
	}
	pharmacy.UpdatedDate = &t
	return session.Save(pharmacy)
}

// Update - update the pharmacy
func Update(pharmacy *data.Pharmacy) (*data.Pharmacy, error) {
	session := data.Session()
	defer session.Close()

	savedPh := new(data.Pharmacy)
	err := session.Find(gmgo.Q{"_id": pharmacy.ID}, savedPh)
	if err != nil {
		return nil, err
	}

	savedPh.Name = pharmacy.Name
	savedPh.Street1 = pharmacy.Street1
	savedPh.Email = pharmacy.Email
	savedPh.City = pharmacy.City
	savedPh.State = pharmacy.State
	savedPh.ZipCode = pharmacy.ZipCode
	savedPh.PartnershipStatus = pharmacy.PartnershipStatus
	savedPh.Partner = pharmacy.Partner
	savedPh.ContactName = pharmacy.ContactName
	savedPh.PhoneNumber = pharmacy.PhoneNumber
	savedPh.FaxNumber = pharmacy.FaxNumber
	savedPh.NPI = pharmacy.NPI
	savedPh.DEA = pharmacy.DEA

	savedPh.UpdatedDate = util.NowUTC()
	return savedPh, session.Update(gmgo.Q{"_id": pharmacy.ID}, savedPh)
}
