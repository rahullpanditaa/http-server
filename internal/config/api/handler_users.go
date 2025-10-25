package api

import (
	"log"
	"net/http"

	"github.com/rahullpanditaa/http-server/internal/auth"
	"github.com/rahullpanditaa/http-server/internal/database"
	"github.com/rahullpanditaa/http-server/internal/handlers"
	"github.com/rahullpanditaa/http-server/internal/handlers/helpers"
)

// HandlerCreateUser is the handler function for the endpoint POST /api/users.
// Get email, password in request body
// create a hashed password from received password
// create a new user in users table in chirpy database using received email and hashed password
// send back a json response with details of user created in previous step
func (handler *ApiConfigHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	userDetailsReceived := helpers.ReadRequestJSON[handlers.User](w, r)

	emailInRequest := userDetailsReceived.Email
	passwordInRequest := userDetailsReceived.Password

	hashedPassword, err := auth.HashPassword(passwordInRequest)
	if err != nil {
		helpers.RespondWithError(w, 500, "unable to hash given password")
		log.Printf("Error: %v\n", err)
		log.Println("Method erred: HandlerCreateUser", "File: handler_users.go", "1")
		return
	}

	user, err := handler.cfg.DbQueries.CreateUser(r.Context(),
		database.CreateUserParams{
			Email:          emailInRequest,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		helpers.RespondWithError(w, 500, "database err, cannot create a user record")
		log.Printf("Error: %v\n", err)
		log.Println("Method erred: HandlerCreateUser", "File: handler_users.go", "2")
		return
	}

	userToReturn := handlers.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email:     user.Email,
		Password:  passwordInRequest,
	}

	helpers.RespondWithJson(w, http.StatusCreated, userToReturn)
}
