package model

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	//"strconv"
	"strings"
	"time"

	//"github.com/google/uuid"
	nano_uuid "github.com/aidarkhanov/nanoid/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserName     string             `json:"user_name" bson:"user_name"`
	Password     string             `json:"password" bson:"password"`
	FullName     string             `json:"full_name" bson:"full_name"`
	ClientUserId string             `json:"client_user_id" bson:"client_user_id"`
	LastLogin    *time.Time         `json:"last_login" bson:"last_login"`
	LastAttempt  *time.Time         `json:"last_attempt" bson:"last_attempt"`
	Attempts     int                `json:"attempts" bson:"attempts"`
	CreatedAt    *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at" bson:"updated_at"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	SessionId  string
	AccessUuid string
	UserId     string
	Token      *jwt.Token
	//ExpiresAt   int64
}

// func CreateAuth(userId string, td *TokenDetails) error {
// 	at := time.Unix(td.AtExpires, 0)
// 	rt := time.Unix(td.RtExpires, 0)

// }

func Login(userName, password string) (string, error) {
	var err error
	pwd := EncryptPassword(password)
	filter := bson.D{{"user_name", userName}, {"password", pwd}}

	collection, _ := storage.GetCollection("users")
	var user User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Debugf("Login:70-- Find Session for UserName: %s returned ERROR: %v", userName, err)
		return "", fmt.Errorf("not Authorized")
	}
	session, err := GetSessionForUserID(user.ID)
	if err != nil { // Session does not exist. Create new session for user
		log.Info("Login:75 - No Session exists, creating")
		_, err = CreateSessionForUser(&user)
		if err != nil {
			return "", fmt.Errorf("unable to create session - %s", err.Error())
		}
	} else {

		session.UpdateSessionAccessedAt()

	}
	session, err = GetSessionForUserID(user.ID)
	if err != nil {
		log.Errorf("Re-Reading session: %s", err.Error())
		return "", err
	}
	//fmt.Printf("Login:90 Returning session: %s\n", spew.Sdump(session))
	return session.ID.Hex(), nil
}

//EncryptPassword: returns an encrypted password for storeage or find
func EncryptPassword(passwd string) string {
	return passwd //TODO: Actually encrypt the password

}

func CreateSessionId() string {
	id, _ := nano_uuid.New()
	return id
	// if err!= nil {
	// 	return "", fmt.Errorf("auth_session:91 -- Could not generate uuid: %s\n", err.Error())
	// }
}

// func ValidateSession(id string) (*AuthSession, error) {
// 	return nil, errors.New("Not Implemented")

// }

func CreateToken(userid string) (string, error) {
	var err error

	id, err := nano_uuid.New()
	if err != nil {
		return "", fmt.Errorf("user:98 -- Could not generate uuid: %s", err.Error())
	}
	//Creating Session

	atExpires := time.Now().Add(time.Minute * 30).Unix()
	td := &TokenDetails{}
	td.AtExpires = atExpires //time.Now().Add(time.Minute * 30).Unix()  //TODO: Token Expiration should come from config
	td.AccessUuid, _ = nano_uuid.New()
	//td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix() //TODO: Token Refresh Expiration should come from config
	//td.RefreshUuid = uuid.NewV4().String()
	td.RefreshUuid, _ = nano_uuid.New()

	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = atExpires
	atClaims["session_id"] = id
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	return td.AccessToken, nil
}

func TokenSignedString(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	return VerifyTokenString(tokenString)
}

func VerifyTokenString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		log.Debugf("ExtractToken:153 - It is a bearer token")
		token := strArr[1]
		return token
	}
	log.Debugf("ExtractToken:156 - It is not a bearer token")
	return ""
}

// func TokenValid(token *jwt.Token) error {
// 	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 		return errors.New("token is invalid")
// 	}
// 	return nil
// }

//   func TokenValid(tkn *jwt.Token) error {
// 	token, err := VerifyToken(r)
// 	if err != nil {
// 	   return err
// 	}
// 	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
// 	   return err
// 	}
// 	return nil
//   }

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	return GetTokenMetaData(token)
}

func GetTokenMetaData(token *jwt.Token) (*AccessDetails, error) {
	var err error
	claims, ok := token.Claims.(jwt.MapClaims)
	//if ok && token.Valid {
	if ok {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, errors.New("access not found in token")
		}
		sessionId, ok := claims["session_id"].(string)
		if !ok {
			return nil, errors.New("session not found in token")
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("user not found in token")
		}

		//    expiresAt, ok := claims["exp"].(int64)
		//    if !ok {
		// 	  return nil, errors.New("exp not found in token")
		//    }

		return &AccessDetails{
			AccessUuid: accessUuid,
			SessionId:  sessionId,
			UserId:     userId,
			Token:      token,
			//ExpiresAt: expiresAt,
		}, nil
	}
	return nil, err
}

func GetClaimItem(tkn *jwt.Token, claim string) (string, bool) {
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	return claims[claim].(string), ok
}

func SetTokenSession(tkn *jwt.Token, sessionId string) *jwt.Token {
	atClaims := tkn.Claims.(jwt.MapClaims)
	//atClaims := jwt.MapClaims{}
	atClaims["session_id"] = sessionId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at
	// newToken, err :=
	// td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	// if err != nil {
	//    return  nil, err
	// }
	// log.Debugf("claims[session_id] = %s", atClaims["session_id"])
	// tkn.Claims = atClaims
	// //log.Debugf("Token.claims[session_id] = %s", atClaims["session_id"])
	// return tkn
}

func UpdateTokenExpire(tkn *jwt.Token) *jwt.Token {
	atClaims := tkn.Claims.(jwt.MapClaims)
	//atClaims := jwt.MapClaims{}
	atClaims["authorized"] = false
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() //TODO: Token Expiration should come from config
	fmt.Printf("\n\n\n####UpdTokenExpire: %s\n", spew.Sdump(atClaims))
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at
}
