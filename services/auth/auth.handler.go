package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"viraj_golang/services/user"
	"viraj_golang/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	router    *http.ServeMux
	authModel *user.UserModel
}

func InitAuthRoutes(database *sqlx.DB) *http.ServeMux {
	authRouter := http.NewServeMux()

	authHandler := &AuthHandler{
		router:    authRouter,
		authModel: &user.UserModel{DB: database},
	}
	authHandler.RegisterAuthRoutes()
	return authRouter
}

func (a *AuthHandler) RegisterAuthRoutes() {
	a.router.HandleFunc("/signup", a.CreateUser)
	a.router.HandleFunc("/login", a.LoginUser)
}

func (a *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var receivedUser *user.User = &user.User{}
	if err := json.NewDecoder(r.Body).Decode(receivedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(receivedUser)
	//check for existence of other user wqith sa,me email
	conflictingusers, err := a.authModel.FindByEmail(receivedUser.Email)
	if err != nil || len(conflictingusers) > 0 {
		// fmt.Println(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			http.Error(w, "user already exists with the given email", http.StatusBadRequest)
			return

		}
	}
	hashPass, err := utils.HashPassword(receivedUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if hashPass != "" {
		receivedUser.Password = hashPass
		result, err := a.authModel.InsertOne(receivedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jwt, err := utils.GenJWT(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		utils.SuccessfullyCreated(w, 201, "Successfully created a user", map[string]any{"insertedId": strconv.Itoa(result), "token": jwt})
		return
	}
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	userDetails := &user.User{}
	if err := json.NewDecoder(r.Body).Decode(userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	usersInDb, err := a.authModel.FindByEmail(userDetails.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(usersInDb) == 1 {
		user := usersInDb[0]
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDetails.Password)); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := utils.GenJWT(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := &utils.JsonResponse{Code: 200, Message: "Logged In User", Body: map[string]any{"user_id": user.ID, "token": token}}
		utils.SuccessfullyFoundOne(w, resp)
	} else {
		http.Error(w, fmt.Errorf("issue finind user with given email").Error(), http.StatusBadRequest)
	}
}
