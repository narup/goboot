package insurance

import (
	"context"
	"encoding/base64"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/util"
)

// GenerateMissingInsuranceToken returns missing insurance token
func GenerateMissingInsuranceToken(managerID string) (*data.MissingInsuranceToken, error) {

	session := data.Session()
	defer session.Close()

	token := util.RandomToken(10)

	it := new(data.MissingInsuranceToken)
	it.InitData()
	it.Token = token
	it.UserID = managerID

	err := session.Save(it)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// FindMissingInsuranceTokenForUser returns the missing insurance token from userID
func FindMissingInsuranceTokenForUser(userID string) (*data.MissingInsuranceToken, error) {
	session := data.Session()
	defer session.Close()

	missingInsuranceToken := new(data.MissingInsuranceToken)
	err := session.Find(gmgo.Q{"userId": userID}, missingInsuranceToken)
	if err != nil && err.Error() != data.ErrNotFound {
		return nil, err
	}

	return missingInsuranceToken, nil
}

// GenerateInvalidInsuranceToken returns insurance exception token
func GenerateInvalidInsuranceToken(managerID string) (*data.InvalidInsuranceToken, error) {

	session := data.Session()
	defer session.Close()

	token := util.RandomToken(10)

	it := new(data.InvalidInsuranceToken)
	it.InitData()
	it.Token = token
	it.UserID = managerID

	err := session.Save(it)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// FindInvalidInsuranceTokenForUser returns the insurance exception token from userID
func FindInvalidInsuranceTokenForUser(userID string) (*data.InvalidInsuranceToken, error) {
	session := data.Session()
	defer session.Close()

	token := new(data.InvalidInsuranceToken)
	err := session.Find(gmgo.Q{"userId": userID}, token)
	if err != nil && err.Error() != data.ErrNotFound {
		return nil, err
	}

	return token, nil
}

// FindInsurance returns the insurance with the id.
func FindInsurance(ctx context.Context, insuranceID string) (*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	ins := new(data.Insurance)
	err := session.FindByID(insuranceID, ins)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

// GetEncodedInsuranceImageFile
func GetEncodedInsuranceImageFile(imageID string) (string, string, error) {
	file, err := GetInsuranceImageFile(imageID)
	if err != nil {
		return "", "", err
	}

	imageData := base64.StdEncoding.EncodeToString(file.Data)
	contentType := file.ContentType
	return imageData, contentType, nil
}

// GetInsuranceImageFile
func GetInsuranceImageFile(imageID string) (*gmgo.File, error) {
	session := data.Session()
	defer session.Close()

	file := new(gmgo.File)
	err := session.ReadFile(imageID, "rex_files", file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// FindAllForPatient finds all the insurances of patient.
func FindAllForPatient(ctx context.Context, patientID string) ([]*data.Insurance, error) {
	session := data.Session()
	defer session.Close()

	result, err := session.FindAll(gmgo.Q{"patientId": patientID}, new(data.Insurance))
	if err != nil {
		return nil, err
	}

	return result.([]*data.Insurance), nil
}
