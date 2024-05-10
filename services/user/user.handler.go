package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"viraj_golang/middlewares"
	"viraj_golang/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	router    *http.ServeMux
	userModel *UserModel
}

func InitUserRoutes(mux *http.ServeMux, database *sql.DB) {
	user := &UserHandler{
		router:    mux,
		userModel: &UserModel{dB: database},
	}
	user.RegisterUserRoutes()
	user.RegisterAuthRoutes()
}

func (u *UserHandler) RegisterAuthRoutes() {
	middlewareStackNonAuth := middlewares.CreateMiddlewareStack(
		middlewares.RecoverPanic,
		middlewares.RequestLogger,
		middlewares.SecureHeaders,
	)

	u.router.HandleFunc("/auth/signup", middlewareStackNonAuth(u.CreateUser))
	u.router.HandleFunc("/auth/login", middlewareStackNonAuth(u.LoginUser))
}

func (user *UserHandler) RegisterUserRoutes() {
	middlewareStack := middlewares.CreateMiddlewareStack(
		middlewares.RecoverPanic,
		middlewares.RequestLogger,
		middlewares.CheckAuth,
		middlewares.SecureHeaders,
	)
	user.router.HandleFunc("GET /users", middlewareStack(user.GetAllUsers))
	user.router.HandleFunc("PUT /users", middlewareStack(user.UpdateUser))
	user.router.HandleFunc("DELETE /users/{id}", middlewareStack(user.DeleteUser))
	user.router.HandleFunc("GET /users/{id}", middlewareStack(user.GetUser))
}

func (user *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var receivedUser *User = &User{}
	if err := json.NewDecoder(r.Body).Decode(receivedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(receivedUser)
	//check for existence of other user wqith sa,me email
	conflictingusers, err := user.userModel.FindByEmail(receivedUser.Email)
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
		result, err := user.userModel.InsertOne(receivedUser)
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
func (user *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "pass valid user id", http.StatusBadRequest)
	}
	if user_id != 0 {
		foundUser, err := user.userModel.GetById(user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if foundUser != nil {
			utils.SuccessfullyFoundOne(w, &utils.JsonResponse{Code: http.StatusFound, Message: "User found", Body: map[string]any{"user": foundUser}})
		}

	}

}
func (user *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users = []*User{}
	users, err := user.userModel.GetAll()
	r.URL.Query()
	if err != nil {
		http.Error(w, "Could not fetch all users", http.StatusInternalServerError)
	}
	resp := &utils.JsonResponse{Code: 200, Message: "Fetched all users", Body: map[string]any{"users": users}}

	utils.SuccessfullyFetchedAll(w, resp)

}
func (user *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
func (user *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}

func (user *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	userDetails := &User{}
	if err := json.NewDecoder(r.Body).Decode(userDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	usersInDb, err := user.userModel.FindByEmail(userDetails.Email)
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
