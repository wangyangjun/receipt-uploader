package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/wangyangjun/receipt-uploader/graph/model"
	"golang.org/x/crypto/bcrypt"
)

const dateFormat = "2006-01-02 00:00:00"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(username string, password string) (*model.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		log.Panicln("Hash password failed")
		return nil, errors.New("Internal error!")
	}

	user := model.User{
		ID:          fmt.Sprintf("%v", uuid.NewV4()),
		Username:    username,
		Password:    &passwordHash,
		DateCreated: time.Now().Format(dateFormat),
	}

	userMap, err := readUserMap()
	if err != nil {
		return nil, err
	}

	if _, ok := userMap[user.Username]; ok {
		return nil, errors.New("Username is already token, please use another one")
	}

	userMap[user.Username] = user
	userMapBytes, err := json.Marshal(userMap)
	if err != nil {
		log.Println(err)
		return nil, errors.New("User data serialization failed")
	}
	err = ioutil.WriteFile(userFilePath, userMapBytes, 0644)

	if err != nil {
		log.Println(err)
		return nil, errors.New("Internal error")
	}
	return user.WithoutPwd(), nil
}

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func generateToken(username string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}

func Login(username string, password string) (*model.AuthPayload, error) {
	userMap, err := readUserMap()
	if err != nil {
		return nil, err
	}

	user, ok := userMap[username]
	if !ok {
		return nil, errors.New("Wrong username or password")
	}
	if !checkPasswordHash(password, *user.Password) {
		return nil, errors.New("Wrong username or password")
	}

	tokenString, err := generateToken(username)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return nil, errors.New("Internal error")
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	return &model.AuthPayload{
		Token: tokenString,
		User:  user.WithoutPwd(),
	}, nil
}

func readUserMap() (map[string]model.User, error) {
	file, err := ioutil.ReadFile(userFilePath)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can not open user file")
	}
	userMap := make(map[string]model.User)
	json.Unmarshal(file, &userMap)
	return userMap, nil
}

func GetAllUsers() ([]*model.User, error) {
	userMap, err := readUserMap()
	if err != nil {
		return nil, err
	}
	users := []*model.User{}

	for _, user := range userMap {
		users = append(users, &user)
	}

	return users, nil
}

func GetUserByUsername(username string) (*model.User, error) {
	userMap, err := readUserMap()
	if err != nil {
		return nil, err
	}
	if user, ok := userMap[username]; ok {
		return user.WithoutPwd(), nil
	}
	return nil, errors.New("No user with such id")
}
