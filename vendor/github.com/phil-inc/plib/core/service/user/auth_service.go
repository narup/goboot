package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	mrand "math/rand"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/narup/gmgo"

	"github.com/phil-inc/plib/core/data"
	"github.com/phil-inc/plib/core/util"
	"golang.org/x/crypto/bcrypt"
)

var digits = "1234567890"

func init() {
	mrand.Seed(time.Now().UnixNano())
}

// JWTToken auth token
type JWTToken struct {
	UserID       string `json:"userID"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshToken represents request data for refresh token
type RefreshToken struct {
	Token string `json:"refreshToken"`
}

// Error which indicates an account is locked
var AccountLockedError = errors.New("account is locked")

// Login check if user email and password is valid
func Login(ctx context.Context, email, password string, isPasswordRaw bool) (JWTToken, error) {
	session := data.Session()
	defer session.Close()

	savedUser := new(data.User)
	if err := session.Find(gmgo.Q{"email": email}, savedUser); err != nil {
		return JWTToken{}, err
	}

	return getJWTTokenAfterEmailOrPhoneVerification(session, savedUser, password, isPasswordRaw)
}

//EmailOrPhoneLogin login using either email or phone
func EmailOrPhoneLogin(ctx context.Context, email, phoneNumber, pwd string) (JWTToken, error) {
	var authToken JWTToken
	var err error

	if email != "" {
		//login with email
		authToken, err = Login(ctx, email, pwd, true)
	} else {
		//login with phone number
		authToken, err = LoginWithPhoneNumber(ctx, phoneNumber, pwd, true)
	}

	return authToken, err
}

// LockAccount adds a timestamp for the time in which the account was locked.
// Note: Updating the user model must be handled in session outside this function.
func LockAccount(ctx context.Context, usr *data.User) {
	if usr.UserAttributes == nil {
		return
	}
	currTime := time.Now()
	usr.UserAttributes.LoginAttemptsExceeded = &currTime
}

// UnlockAccount removes the lock timestamp, and resets the Login attempts to 0.
// Note: Updating the user model must be handled in session outside this function.
func UnlockAccount(ctx context.Context, usr *data.User) {
	if usr.UserAttributes == nil {
		return
	}
	usr.UserAttributes.LoginAttemptsCount = 0
	usr.UserAttributes.LoginAttemptsExceeded = nil
}

// AccountIsLocked checks if the maximum attempts has been reached.
// The maximum attempts will be reset to 0, if enough time has elapsed.
// Note: Updating the user model must be handled in session outside this function.
func AccountIsLocked(ctx context.Context, usr *data.User) bool {
	if usr.UserAttributes == nil {
		return false
	}
	maxAttempts, _ := strconv.Atoi(util.Config("auth.maxLoginAttempts"))

	if usr.UserAttributes.LoginAttemptsCount >= maxAttempts {
		// If count is greater than or equal to 3, check when they were locked out
		hourFromLocked := usr.UserAttributes.LoginAttemptsExceeded.Add(time.Hour)
		if time.Now().After(hourFromLocked) {
			// If an hour has elapsed, unlock account
			UnlockAccount(ctx, usr)

			return false
		}
		// If hour has not yet passed, account is still locked
		return true
	}

	return false
}

// LoginAttempt checks the number of attempts a user has made at logging in,
// and properly locks/unlocks their account.
func LoginAttempt(ctx context.Context, email, phoneNumber, pwd string) (JWTToken, error) {
	var authToken JWTToken

	session := data.Session()
	defer session.Close()

	usr, err := FindByEmailOrPhoneInSession(ctx, session, email, phoneNumber)
	if err != nil {
		return authToken, err
	}

	defer func() {
		UpdateInSession(usr, session)
	}()

	if usr.UserAttributes == nil {
		usr.UserAttributes = new(data.UserAttributes)
	}

	accountLocked := AccountIsLocked(ctx, usr)
	if accountLocked {
		return authToken, AccountLockedError
	}

	usr.UserAttributes.LoginAttemptsCount++ // increment attempts

	authToken, err = EmailOrPhoneLogin(ctx, email, phoneNumber, pwd)
	if err != nil {
		maxAttempts, _ := strconv.Atoi(util.Config("auth.maxLoginAttempts"))

		if usr.UserAttributes.LoginAttemptsCount >= maxAttempts { // 3 failed attempts is the max allowed
			LockAccount(ctx, usr)

			return authToken, AccountLockedError
		}

		return authToken, err
	}

	UnlockAccount(ctx, usr)
	return authToken, nil
}

// LoginWithPhoneNumber check if user email and password is valid
func LoginWithPhoneNumber(ctx context.Context, phoneNumber, password string, isPasswordRaw bool) (JWTToken, error) {
	session := data.Session()
	defer session.Close()

	results, err := session.FindAll(gmgo.Q{"phoneNumber": phoneNumber}, new(data.User))
	if err != nil {
		return JWTToken{}, err
	}

	users := results.([]*data.User)
	if len(users) == 1 {
		return getJWTTokenAfterEmailOrPhoneVerification(session, users[0], password, isPasswordRaw)
	}
	return JWTToken{}, errors.New("multiple accounts or no account")
}

func getJWTTokenAfterEmailOrPhoneVerification(session *gmgo.DbSession, savedUser *data.User, password string, isPasswordRaw bool) (JWTToken, error) {
	//check password
	if isPasswordRaw {
		matchError := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(password))
		if matchError != nil {
			log.Printf("[ERROR] Password hash comparision failed for user ID: %s, Error: %s", savedUser.StringID(), matchError)
			return JWTToken{}, errors.New("Password doesn't match")
		}
	} else {
		if savedUser.Password != password {
			log.Printf("[ERROR] Password hash comparision failed for user ID: %s", savedUser.StringID())
			return JWTToken{}, errors.New("Password doesn't match")
		}
	}

	//generate token values
	roles := getUserRoles(savedUser, session)
	atValue, err := generateJWTAccessToken(savedUser, roles)
	rtValue := util.RandomToken(16)

	newUserToken := new(data.UserToken)
	err = session.Find(gmgo.Q{"userId": savedUser.StringID()}, newUserToken)
	if err != nil && err.Error() != data.ErrNotFound {
		return JWTToken{}, err
	}

	if newUserToken.ID == "" {
		newUserToken.InitData()
		newUserToken.RefreshToken = rtValue
		newUserToken.UserID = savedUser.StringID()
		newUserToken.BlackListed = false

		err = setTokenExpiration(newUserToken)
		if err == nil {
			err = session.Save(newUserToken)
		}
	} else {
		newUserToken.RefreshToken = rtValue
		setTokenExpiration(newUserToken)

		err = session.Update(gmgo.Q{"_id": newUserToken.ID}, newUserToken)
	}

	if err != nil {
		return JWTToken{}, err
	}

	return JWTToken{UserID: savedUser.StringID(), AccessToken: atValue, RefreshToken: rtValue}, nil
}

// Log user out by removing their refresh token from the db
func Logout(ctx context.Context, userID string) error {
	session := data.Session()
	defer session.Close()

	savedUserToken := new(data.UserToken)
	err := session.Find(gmgo.Q{"userId": userID}, savedUserToken)
	if err != nil {
		return err
	}

	err = session.Remove(gmgo.Q{"_id": savedUserToken.ID}, savedUserToken)
	return err
}

func setTokenExpiration(ut *data.UserToken) error {
	// Use refresh token expiration time from configuration
	refreshTokenExpiration, err := strconv.Atoi(util.Config("auth.refreshTokenExpiration"))
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Minute * time.Duration(refreshTokenExpiration))
	ut.TokenExpiration = &exp

	return nil
}

// RefreshAccessToken check if user email and password is valid
func RefreshAccessToken(ctx context.Context, userID, refreshToken string) (JWTToken, error) {
	session := data.Session()
	defer session.Close()

	savedUserToken := new(data.UserToken)
	err := session.Find(gmgo.Q{"userId": userID}, savedUserToken)
	if err != nil {
		return JWTToken{}, err
	}

	if savedUserToken.BlackListed == true || savedUserToken.RefreshToken != refreshToken {
		return JWTToken{}, errors.New("Invalid refresh token")
	}

	// Check if refresh token has expired
	if savedUserToken.TokenExpiration != nil && savedUserToken.TokenExpiration.Before(time.Now()) {
		// If token has expired, clear it from db so user can create a new session by logging in
		session.Remove(gmgo.Q{"_id": savedUserToken.ID}, savedUserToken)
		return JWTToken{}, errors.New("Expired refresh token")
	}

	currentUser := new(data.User)
	err = session.FindByID(userID, currentUser)
	if err != nil {
		return JWTToken{}, err
	}

	roles := getUserRoles(currentUser, session)
	accessTokenValue, err := generateJWTAccessToken(currentUser, roles)
	if err != nil {
		return JWTToken{}, err
	}

	savedUserToken.RefreshToken = util.RandomToken(16)
	err = session.Update(gmgo.Q{"_id": savedUserToken.ID}, savedUserToken)
	if err != nil {
		return JWTToken{}, err
	}

	return JWTToken{UserID: userID, AccessToken: accessTokenValue, RefreshToken: savedUserToken.RefreshToken}, nil
}

//CreateSessionToken creates single sign on session token.
//This token is only valid for 2 mins for security reasons
func CreateSessionToken(ctx context.Context, userID string) (string, error) {
	session := data.Session()
	defer session.Close()

	st := new(data.SessionToken)
	st.InitData()
	st.UserID = userID
	token, err := GenerateRandomString(32)
	if err != nil {
		log.Printf("[ERROR] error generating secure session token. Error: %s", err)
		st.Token = util.RandomToken(32)
	} else {
		st.Token = token
	}

	return st.Token, session.Save(st)
}

// LoginWithSessionToken returns JWT token based on single sign-on session token
func LoginWithSessionToken(ctx context.Context, sessionToken string) (JWTToken, error) {
	session := data.Session()
	defer session.Close()

	q := gmgo.Q{
		"token": sessionToken,
	}

	st := new(data.SessionToken)
	err := session.Find(q, st)
	if err != nil {
		return JWTToken{}, fmt.Errorf("Session token invalid")
	}

	createdDate := *st.CreatedDate
	sinceCreated := time.Since(createdDate)
	if sinceCreated.Minutes() > 5 {
		return JWTToken{}, fmt.Errorf("Session token expired")
	}

	usr, err := FindByIDInSession(ctx, st.UserID, session)
	if err != nil {
		return JWTToken{}, err
	}

	return getJWTTokenForUser(ctx, usr, session)
}

// LoginWithUserID returns JWT token based on userID
func LoginWithUserID(ctx context.Context, userID string) (JWTToken, error) {
	session := data.Session()
	defer session.Close()

	currentUser := new(data.User)
	err := session.FindByID(userID, currentUser)
	if err != nil {
		return JWTToken{}, err
	}

	return getJWTTokenForUser(ctx, currentUser, session)
}

//FindTwoStepTokenByCode - fetch user by user ID and 2-factor authentication token
func FindTwoStepTokenByCode(ctx context.Context, userID, code string) (*data.TwoStepToken, error) {
	session := data.Session()
	defer session.Close()

	twoStepToken := new(data.TwoStepToken)
	return twoStepToken, session.Find(gmgo.Q{"code": code, "userId": userID}, twoStepToken)
}

//GenerateTwoStepToken generate 2-step token for JWT user access token
func GenerateTwoStepToken(ctx context.Context, jwtToken JWTToken) (*data.TwoStepToken, error) {
	session := data.Session()
	defer session.Close()

	st := new(data.TwoStepToken)
	st.InitData()
	st.UserID = jwtToken.UserID
	st.AccessToken = jwtToken.AccessToken
	st.RefreshToken = jwtToken.RefreshToken
	st.Code = generateTwoStepAuthCode()

	return st, session.Save(st)
}

//FindUserByResetToken - fetch user by password reset token
func FindUserByResetToken(ctx context.Context, token string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	resetToken := new(data.PasswordResetToken)
	err := session.Find(gmgo.Q{"token": token}, resetToken)
	if err != nil {
		return nil, err
	}

	usr := new(data.User)
	err = session.FindByRef(resetToken.UserRef, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

//GenerateIdentityToken creates identity token for the user ID
func GenerateIdentityToken(ctx context.Context, userID string) (string, error) {
	session := data.Session()
	defer session.Close()

	st := new(data.IdentityToken)
	st.InitData()
	st.UserID = userID
	token, err := GenerateRandomString(8)
	if err != nil {
		log.Printf("[ERROR] error generating identity token. Error: %s", err)
		st.Token = util.RandomToken(8)
	} else {
		st.Token = token
	}

	return st.Token, session.Save(st)
}

//FindUserByIdentityToken - fetch user by identity token
func FindUserByIdentityToken(ctx context.Context, token string) (*data.User, error) {
	session := data.Session()
	defer session.Close()

	it := new(data.IdentityToken)
	err := session.Find(gmgo.Q{"token": token}, it)
	if err != nil {
		return nil, err
	}

	usr := new(data.User)
	err = session.FindByID(it.UserID, usr)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

//MarkIdentityTokenInactive mark identity token for the token as inactive
func MarkIdentityTokenInactive(ctx context.Context, token string) error {
	session := data.Session()
	defer session.Close()

	it := data.IdentityToken{}
	return session.UpdateFieldValue(gmgo.Q{"token": token}, it.CollectionName(), "status", "Inactive")
}

func getJWTTokenForUser(ctx context.Context, savedUser *data.User, session *gmgo.DbSession) (JWTToken, error) {
	//generate token values
	roles := getUserRoles(savedUser, session)
	atValue, err := generateJWTAccessToken(savedUser, roles)
	rtValue := util.RandomToken(16)

	newUserToken := new(data.UserToken)
	err = session.Find(gmgo.Q{"userId": savedUser.StringID()}, newUserToken)
	if err != nil && err.Error() != data.ErrNotFound {
		return JWTToken{}, err
	}

	if newUserToken.ID == "" {
		newUserToken.InitData()
		newUserToken.RefreshToken = rtValue
		newUserToken.UserID = savedUser.StringID()
		newUserToken.BlackListed = false

		// Use refresh token expiration time from configuration
		rte, rerr := strconv.Atoi(util.Config("auth.refreshTokenExpiration"))
		if rerr != nil {
			return JWTToken{}, rerr
		}
		exp := time.Now().Add(time.Minute * time.Duration(rte))

		newUserToken.TokenExpiration = &exp
		err = session.Save(newUserToken)
	} else {
		newUserToken.RefreshToken = rtValue
		err = session.Update(gmgo.Q{"_id": newUserToken.ID}, newUserToken)
	}

	if err != nil {
		return JWTToken{}, err
	}

	return JWTToken{UserID: savedUser.StringID(), AccessToken: atValue, RefreshToken: rtValue}, nil
}

func generateJWTAccessToken(user *data.User, roles []string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain

	//iss: The issuer of the token
	//nbf: Defines the time before which the JWT MUST NOT be accepted for processing
	//iat: The time the JWT was issued. Can be used to determine the age of the JWT
	//exp: This will probably be the registered claim most often used. This will define the expiration in NumericDate value.
	//     the expiration MUST be after the current date/time.

	// Use access token expiration time from configuration
	accessTokenExpiration, err := strconv.Atoi(util.Config("auth.accessTokenExpiration"))
	if err != nil {
		return "", err
	}
	exp := time.Now().Add(time.Minute * time.Duration(accessTokenExpiration)).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   "phil.us",
		"uid":   user.StringID(),
		"roles": roles,
		"nbf":   time.Now().Unix(),
		"iat":   time.Now().Unix(),
		"exp":   exp,
	})

	// Sign and get the complete encoded token as a string using the secret
	key := []byte(util.Config("auth.hmacToken"))
	tokenString, err := accessToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//FindUserRoles returns the list of user roles for the given user
func FindUserRoles(ctx context.Context, user *data.User) []string {
	session := data.Session()
	defer session.Close()

	return getUserRoles(user, session)
}

func getUserRoles(user *data.User, session *gmgo.DbSession) []string {

	rs := make([]string, 0)
	for _, rRef := range user.RoleRefs {
		role := new(data.UserRole)
		err := session.FindByRef(rRef, role)
		if err != nil {
			continue
		}
		rs = append(rs, role.Name)
	}
	return rs
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// generateTwoStepAuthCode generates the random order number of format xx-xxx
func generateTwoStepAuthCode() string {
	size := 6
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		if i == 2 {
			buf[i] = '-'
		} else {
			buf[i] = digits[mrand.Intn(len(digits))]
		}
	}
	return string(buf)
}
