package user

import (
	"net/http"
	"strconv"

	"github.com/virajBaswana/go_App/middlewares"
	"github.com/virajBaswana/go_App/utils"

	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	router    *http.ServeMux
	userModel *UserModel
}

func InitRoutes(database *sqlx.DB) *http.ServeMux {
	userRouter := http.NewServeMux()
	middlewares.CheckAuth(userRouter)

	user := &UserHandler{
		router:    userRouter,
		userModel: &UserModel{DB: database},
	}
	user.RegisterUserRoutes()
	return userRouter
}

func (user *UserHandler) RegisterUserRoutes() {

	user.router.HandleFunc("GET /getAllUsers", user.GetAllUsers)
	user.router.HandleFunc("PUT /updateUser", user.UpdateUser)
	user.router.HandleFunc("DELETE /users/{id}", user.DeleteUser)
	user.router.HandleFunc("GET /getUser/{id}", user.GetUser)
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

	userId := utils.ExtractClaimsFromRequest(r.Context())
	if userId == "" {
		http.Error(w, "login and process jwt properly", http.StatusUnauthorized)
	} else {
		r.URL.Query()
		if err != nil {
			http.Error(w, "Could not fetch all users", http.StatusInternalServerError)
		}
		resp := &utils.JsonResponse{Code: 200, Message: "Fetched all users", Body: map[string]any{"users": users}}

		utils.SuccessfullyFetchedAll(w, resp)

	}

}
func (user *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
func (user *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
