package doctor

import (
	"context"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
)

// Update - updates doctor
func Update(ctx context.Context, doctorInfo *data.Doctor) (*data.Doctor, error) {
	session := data.Session()
	defer session.Close()

	err := session.Update(gmgo.Q{"_id": doctorInfo.ID}, doctorInfo)
	if err != nil {
		return nil, err
	}

	return doctorInfo, nil
}

// FindByID find the doctor by ID
func FindByID(ctx context.Context, doctorID string) (*data.Doctor, error) {
	session := data.Session()
	defer session.Close()

	result := new(data.Doctor)
	err := session.FindByID(doctorID, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindByIDInSession find doctor by id using session
func FindByIDInSession(ctx context.Context, doctorID string, session *gmgo.DbSession) (*data.Doctor, error) {
	result := new(data.Doctor)
	err := session.FindByID(doctorID, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//FindAllMDPartners returns the list of all MD partners of Phil
func FindAllMDPartners(ctx context.Context) ([]*data.MDPartner, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{}, new(data.MDPartner))
	if err != nil {
		return nil, err
	}
	return result.([]*data.MDPartner), nil
}

//FindMDPartnerForDoctorID returns the MDPartner data for the given Doctor ID.
func FindMDPartnerForDoctorID(ctx context.Context, doctorID string) (*data.MDPartner, error) {
	session := data.Session()
	defer session.Close()

	return FindMDPartnerForDoctorIDInSession(ctx, doctorID, session)
}

//FindMDPartnerForDoctorIDInSession returns the MDPartner data for the given Doctor ID using given session
func FindMDPartnerForDoctorIDInSession(ctx context.Context, doctorID string, session *gmgo.DbSession) (*data.MDPartner, error) {
	mdPartner := new(data.MDPartner)
	return mdPartner, session.Find(gmgo.Q{"doctor.$id": data.ObjectID(doctorID)}, mdPartner)
}

//FindMDOfficeUsersForDoctorID returns the list of MDOfficeUser for given doctorID.
func FindMDOfficeUsersForDoctorID(ctx context.Context, doctorID string) ([]*data.MDOfficeUser, error) {
	session := data.Session()
	defer session.Close()

	return FindMDOfficeUsersForDoctorIDInSession(ctx, doctorID, session)
}

//FindMDOfficeUsersForDoctorIDInSession returns the list of MDOfficeUser for given doctorID.
func FindMDOfficeUsersForDoctorIDInSession(ctx context.Context, doctorID string, session *gmgo.DbSession) ([]*data.MDOfficeUser, error) {
	q := gmgo.Q{
		"doctorIds": doctorID,
	}

	result, err := session.FindAll(q, new(data.MDOfficeUser))
	if err != nil {
		return nil, err
	}

	return result.([]*data.MDOfficeUser), nil
}

//FindMDOfficeUserForUserID find MDOfficeUser data for Phil user ID
func FindMDOfficeUserForUserID(ctx context.Context, userID string) (*data.MDOfficeUser, error) {
	session := data.Session()
	defer session.Close()

	return FindMDOfficeUserForUserIDInSession(ctx, userID, session)
}

//FindMDOfficeUserForUserIDInSession find MDOfficeUser data for Phil user ID in a given session
func FindMDOfficeUserForUserIDInSession(ctx context.Context, userID string, session *gmgo.DbSession) (*data.MDOfficeUser, error) {
	q := gmgo.Q{
		"userId": userID,
	}

	ou := new(data.MDOfficeUser)
	return ou, session.Find(q, ou)
}
