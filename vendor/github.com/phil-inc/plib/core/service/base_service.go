package base

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/narup/gmgo"
	"github.com/phil-inc/plib/core/util"

	"log"

	"fmt"

	"strings"

	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/network"
)

// FindAdminConfig find admin config
func FindAdminConfig(ctx context.Context) (*data.AdminConfig, error) {
	session := data.Session()
	defer session.Close()

	cfg := new(data.AdminConfig)
	err := session.Find(gmgo.Q{}, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// FindAppConfig find admin config
func FindAppConfig(ctx context.Context) (*data.AppConfig, error) {
	session := data.Session()
	defer session.Close()

	cfg := new(data.AppConfig)
	err := session.Find(gmgo.Q{}, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func IsWhiteListedPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	if phoneNumber == "" {
		return false, errors.New("phone number is empty")
	}
	cfg, err := FindAppConfig(ctx)
	if err != nil {
		return false, err
	}
	if cfg == nil {
		return false, errors.New("app config is empty")
	}

	for _, num := range cfg.WhiteListedPhoneNumbers {
		if num == phoneNumber {
			return true, nil
		}
	}

	return false, nil

}

//GetUserMaskNumber returns twilio number that is from user's state
func GetUserMaskNumber(ctx context.Context, usr *data.User) string {
	if usr.FromNumber != "" {
		return usr.FromNumber
	}
	cfg, err := FindAdminConfig(ctx)
	if err != nil {
		return util.Config("twilio.userNumber")
	}

	num := cfg.OutboundNumberMap[usr.State]
	if num == "" {
		return util.Config("twilio.userNumber")
	}
	return num
}

// LookupZipCode does a database lookup for given zipcode to find city and state
func LookupZipCode(ctx context.Context, zipCode string) (*data.Zipcode, error) {
	session := data.Session()
	defer session.Close()

	z := new(data.Zipcode)
	err := session.Find(gmgo.Q{"zip": zipCode}, z)
	if err != nil {
		return nil, err
	}
	return z, nil
}

// VerifyAddress verify given address using Shippo
func VerifyAddress(ctx context.Context, address *data.Address) (*data.Address, error) {
	shippoURL := "https://api.goshippo.com/addresses/"
	shippoAuthKey := fmt.Sprintf("ShippoToken %s", util.Config("shippo.apiKey"))
	headers := map[string]string{"Authorization": shippoAuthKey}

	// Build out the data for our message
	v := url.Values{}
	v.Set("object_purpose", "QUOTE")
	v.Set("street1", address.Street1)
	v.Set("street2", address.Street2)
	v.Set("city", address.City)
	v.Set("state", address.State)
	v.Set("zip", address.ZipCode)
	v.Set("country", "US")
	v.Set("validate", "true")

	body, err := network.HTTPPost(shippoURL, v, headers)
	if err != nil {
		log.Printf("Error validating address: %s\n", err)
		return nil, err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("Invalid response from Shippo %s\n", err)
		return nil, err
	}
	objectState, ok := resp["object_state"].(string)
	if ok && objectState != "INVALID" && objectState != "INCOMPLETE" {
		validAddr := new(data.Address)
		validAddr.Street1 = resp["street1"].(string)
		validAddr.Street2 = resp["street2"].(string)
		validAddr.State = resp["state"].(string)
		validAddr.City = resp["city"].(string)

		fullZipCode := resp["zip"].(string)
		elements := strings.Split(fullZipCode, "-")
		validAddr.ZipCode = elements[0]

		return validAddr, nil
	}

	errMessages := resp["messages"].([]interface{})
	if len(errMessages) > 0 {
		errMessage := errMessages[0].(map[string]interface{})
		errText := errMessage["text"].(string)

		return nil, errors.New(errText)
	}

	return nil, errors.New("Unknown error validating address")
}

// SaveInsuranceFile saves insurance file with given name and conent type
func SaveInsuranceFile(fileName, contentType string, fileData []byte) (string, error) {
	session := data.Session()
	defer session.Close()

	f := gmgo.File{}
	f.Name = fileName
	f.ContentType = contentType
	f.Data = fileData

	fileID, err := session.SaveFile(f, "rex_files")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

// GetCacheValue get value from cache for given key
func GetCacheValue(key string) (string, error) {
	return data.Cache().GetStringValue(key)
}

// SetCacheValue set given value from cache for given key
func SetCacheValue(key, value string, expire int) error {
	return data.Cache().SetStringValue(key, value, expire)
}
