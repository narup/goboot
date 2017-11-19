package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"strings"

	"log"

	"reflect"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/narup/gconfig"
	"github.com/phil-inc/plib/core/data"
	"github.com/pkg/errors"
)

const localDateFormat = "2006-01-02" // yyyy-mm-dd
const localDateFormatWithTime = "2006-01-02 15:04 MST"

type sessionUser struct {
	Key string
}

// SessionUserKey key for context
var SessionUserKey = sessionUser{Key: "SessionUser"}

///// Common functions ///////

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Empty empty string
var Empty = ""

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Manager reads current manager from the context
func Manager(ctx context.Context) *data.User {
	uid := ManagerIDFromContext(ctx)

	manager := new(data.User)
	manager.ID = data.ObjectID(uid)

	return manager
}

//ManagerIDFromContext fetch user Id from context
func ManagerIDFromContext(ctx context.Context) string {
	if ctx.Value(SessionUserKey) != nil {
		jwtClaims := ctx.Value(SessionUserKey).(jwt.MapClaims)
		return jwtClaims["uid"].(string)
	}
	return ""
}

// BasicAuthHeader returns authorization header for Phil API server
func BasicAuthHeader() string {
	authHeader := fmt.Sprintf("%s:%s", Config("api.server.clientID"), Config("api.server.clientSecret"))
	return B64Encode(authHeader)
}

//B64Encode - encodes the given string with base 64
func B64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// HandleRefLoadError wraps all the given errors
func HandleRefLoadError(message string, errs []error) error {
	var err error
	for i := range errs {
		if errs[i] != nil {
			if err == nil {
				err = errs[i]
			}
			err = errors.Wrapf(err, message)
		}
	}
	return err
}

// WrapErrors wraps all the given errors
func WrapErrors(message string, errs []error) error {
	var err error
	for i := range errs {
		if errs[i] != nil {
			if err == nil {
				err = errs[i]
			}
			err = errors.Wrapf(err, message)
		}
	}
	return err
}

// EmptyRxList returns the empty prescription data list
func EmptyRxList() []*data.Prescription {
	return []*data.Prescription{}
}

//EmptyUserList returns empty user data list
func EmptyUserList() []*data.User {
	return []*data.User{}
}

// Config returns string configuration for given key
func Config(key string) string {
	return gconfig.Gcg.GetString(key)
}

// BoolConfig returns config boolean value for the given key
func BoolConfig(key string) bool {
	return gconfig.Gcg.GetBool(key)
}

// IntConfig returns config integer value for the given key
func IntConfig(key string) int {
	return gconfig.Gcg.GetInt(key)
}

//SafeConfig returns empty string if the key doesn't exists
func SafeConfig(key string) string {
	if gconfig.Gcg.Exists(key) {
		return gconfig.Gcg.GetString(key)
	}
	return ""
}

// RandomToken generates random string token
func RandomToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// FormatPhone formats the given phone number
func FormatPhone(phoneNumber string) string {
	if len(phoneNumber) < 10 {
		return phoneNumber
	}

	phoneNumber = SanitizePhoneNumber(phoneNumber)
	return fmt.Sprintf("(%s) %s-%s", phoneNumber[:3], phoneNumber[3:6], phoneNumber[6:10])
}

//SanitizePhoneNumber cleans up the phone number
func SanitizePhoneNumber(phoneNumber string) string {
	phoneNumber = strings.Replace(phoneNumber, "(", "", -1)
	phoneNumber = strings.Replace(phoneNumber, ")", "", -1)
	phoneNumber = strings.Replace(phoneNumber, " ", "", -1)
	phoneNumber = strings.Replace(phoneNumber, "-", "", -1)
	if strings.HasPrefix(phoneNumber, "+1") {
		phoneNumber = strings.Replace(phoneNumber, "+1", "", -1)
	}
	return phoneNumber
}

//USDFormat formats the given currency as 2 decimal place USD
func USDFormat(v string) string {
	val, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return ""
	}
	if val < 0 {
		return fmt.Sprintf("-$%.2f", val)
	}
	return fmt.Sprintf("$%.2f", val)
}

//AddStringValues adds 2 string formatted float numbers
func AddStringValues(v1, v2 string) string {
	val1, err := strconv.ParseFloat(v1, 32)
	if err != nil {
		val1 = 0.0
	}

	val2, err := strconv.ParseFloat(v2, 32)
	if err != nil {
		val2 = 0.0
	}

	result := val1 + val2
	return fmt.Sprintf("%.2f", result)
}

//SubtractStringValues subtracts v2 by v1
func SubtractStringValues(v1, v2 string) string {
	val1, err := strconv.ParseFloat(v1, 32)
	if err != nil {
		val1 = 0.0
	}

	val2, err := strconv.ParseFloat(v2, 32)
	if err != nil {
		val2 = 0.0
	}

	result := val1 - val2
	return fmt.Sprintf("%.2f", result)
}

//NowUTC returns current date time in UTC format
func NowUTC() *time.Time {
	ts := time.Now().UTC()
	return &ts
}

//NowPST returns local pacific time
func NowPST() *time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	ts := time.Now().In(loc)
	return &ts
}

// YesterdayPST returns local pacific time yesterday
func YesterdayPST() *time.Time {
	now := NowPST()
	t := now.Add(-24 * time.Hour)
	return &t
}

// YesterdayStartEndTimePST start and end time for yesterday
func YesterdayStartEndTimePST() (*time.Time, *time.Time) {
	yest := YesterdayPST()

	ystart := time.Date(yest.Year(), yest.Month(), yest.Day(), 0, 0, 0, 0, yest.Location())
	yend := time.Date(yest.Year(), yest.Month(), yest.Day(), 23, 59, 59, 0, yest.Location())

	return &ystart, &yend
}

//DayStartTimePSTFor returns time based on year, month and day
func DayStartTimePSTFor(year, month, day int) time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

//DayEndTimePSTFor returns time based on year, month and day
func DayEndTimePSTFor(year, month, day int) time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	return time.Date(year, time.Month(month), day, 23, 59, 59, 0, loc)
}

// DayStartTimePST day start time PST
func DayStartTimePST() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	//set timezone,
	pst := time.Now().In(loc)
	return time.Date(pst.Year(), pst.Month(), pst.Day(), 0, 0, 0, 0, loc)
}

// DayEndTimePST day end time PST
func DayEndTimePST() time.Time {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	pst := time.Now().In(loc)
	return time.Date(pst.Year(), pst.Month(), pst.Day(), 23, 59, 59, 0, loc)
}

//DayStartTime returns start date and time of the day
func DayStartTime() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

//DayEndTime returns end date and time of the day
func DayEndTime() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
}

//MonthStartEndTimePST start and end time for the given year and month in PST
func MonthStartEndTimePST(year, month int) (time.Time, time.Time) {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	monthEnd := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, loc)

	return monthStart, monthEnd
}

// FormatDate formats the date to string using simple format yyyy-mm-dd
func FormatDate(date *time.Time) string {
	return date.Format(localDateFormat)
}

// FormatDateWithTime formats the date to string using simple format yyyy-mm-dd h:mm a tz
func FormatDateWithTime(date *time.Time) string {
	return date.Format(localDateFormatWithTime)
}

//USFormatDate date in US format
func USFormatDate(date *time.Time) string {
	return fmt.Sprintf("%s/%d", date.Format("01/02"), date.Year())
}

// HumanDate returns human readable date
func HumanDate(date *time.Time) string {
	return date.Format("Mon, Jan 2, 2006")
}

//IsLocal checks if it's local environment
func IsLocal() bool {
	return Config("app.environment") == "local"
}

//IsDev returns if the application is running in dev environment
func IsDev() bool {
	return Config("app.environment") == "dev" || Config("app.environment") == "local"
}

//IsProd check if application is running in prod env
func IsProd() bool {
	return Config("app.environment") == "prod"
}

// IsDebugMode check if app is running in debug mode. Does heavy logging
func IsDebugMode() bool {
	return gconfig.Gcg.GetBool("app.debugMode")
}

// IsClosed checks if passed in channel is closed
func IsClosed(ch <-chan string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

//ToJSON to JSON string
func ToJSON(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON %s\n", err)
		return ""
	}
	return string(jsonBytes)
}

//TypeName returns the string value of the interface v
func TypeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

//FormatPrice formats the price in string format as 2 decimal value
func FormatPrice(price string) string {
	if price == "" || price == "-1" {
		return "0.00"
	}
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		val = 0.0
	}
	return fmt.Sprintf("%.2f", val)
}

//FirstName returns the first name only for the given full name
func FirstName(fullName string) string {
	s := strings.Split(fullName, " ")
	if len(s) > 0 {
		fn := strings.ToLower(s[0])
		return strings.Title(fn)
	}
	return fullName
}

//LastName returns the last name only for the given full name
func LastName(fullName string) string {
	s := strings.Split(fullName, " ")
	if len(s) > 0 {
		ln := strings.ToLower(s[len(s)-1])
		return strings.Title(ln)
	}
	return fullName
}

//PartialName returns the full first name and first letter of the middle or last name
func PartialName(fullName string) string {
	s := strings.Split(fullName, " ")
	if len(s) > 0 {
		fn := strings.ToLower(s[0])
		ln := strings.ToLower(s[len(s)-1])

		firstChar := string(ln[0])
		pn := fmt.Sprintf("%s %s", fn, firstChar)

		return strings.Title(pn)
	}
	return fullName
}
