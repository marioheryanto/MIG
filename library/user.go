package library

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
	"github.com/marioheryanto/MIG/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserLibrary struct {
	Repo repository.UserRepositoryInterface
}

type UserLibraryInterface interface {
	Register(user model.UserRequest) error
	Login(user model.UserRequest) (fiber.Cookie, error)
	Home(tokenString string) (interface{}, error)
}

func NewUserLibrary(repo repository.UserRepositoryInterface) UserLibraryInterface {
	return UserLibrary{
		Repo: repo,
	}
}

func (l UserLibrary) Register(request model.UserRequest) error {
	user := model.User{}
	user.Name = request.Name
	user.Email = strings.ToLower(request.Email)

	exist, err := l.Repo.CheckUserExistWithEmail(request.Email)
	if err != nil {
		return err
	}

	if exist {
		return helper.NewServiceError(http.StatusBadRequest, "email already taken")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	err = l.Repo.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (l UserLibrary) Login(request model.UserRequest) (fiber.Cookie, error) {
	user := model.User{}
	user.Email = strings.ToLower(request.Email)
	cookie := fiber.Cookie{}

	err := l.Repo.GetUserWithEmail(&user)
	if err != nil {
		return cookie, helper.ReplaceServiceErrorForLogin(err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
	if err != nil {
		return cookie, helper.ReplaceServiceErrorForLogin(err)
	}

	ttl := time.Now().Add(15 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: ttl.Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return cookie, err
	}

	cookie.Name = "jwt"
	cookie.Value = tokenString
	cookie.Expires = ttl
	cookie.HTTPOnly = true

	return cookie, nil
}

func (l UserLibrary) Home(tokenString string) (interface{}, error) {
	user := model.User{}

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return nil, err
	}

	err = l.Repo.GetUserWithID(claims.Issuer, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
