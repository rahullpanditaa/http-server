package api

import (
	"net/http"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/config"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/helpers"
)

// HandlerCreateUser is the handler function for the endpoint POST /api/users.
// Get email, password in request body
// create a hashed password from received password
// create a new user in users table in chirpy database using received email and hashed password
// send back a json response with details of user created in previous step
func (handler *ApiConfigHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	userDetailsReceived := helpers.ReadRequestJSON[config.User](w, r)

	emailInRequest := userDetailsReceived.Email
	passwordInRequest := userDetailsReceived.Password

	hashedPassword, err := auth.HashPassword(passwordInRequest)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "unable to hash given password")
		helpers.LogErrorWithRequest(err, r, "unable to hash given password")
		return
	}

	user, err := handler.Cfg.DbQueries.CreateUser(r.Context(),
		database.CreateUserParams{
			Email:          emailInRequest,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		helpers.RespondWithError(w, 500, "database error, cannot create a user record")
		helpers.LogErrorWithRequest(err, r, "database error, cannot create a user record")
		return
	}

	userToReturn := config.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email:     user.Email,
		Password:  passwordInRequest,
	}

	helpers.RespondWithJson(w, http.StatusCreated, userToReturn)
}
