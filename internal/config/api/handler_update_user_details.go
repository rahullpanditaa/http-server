package api

import (
	"net/http"
)

// HandlerUpdateUserDetails is the handler function for the endpoint PUT /api/users.
// Via this endpoint, users can update their own email and password.
// Requires an access token (jwt) in header.
// new email, password in request body
func (handler *ApiConfigHandler) HandlerUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	// refreshTokenReceivedInHeader, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	helpers.RespondWithError(w, http.StatusBadRequest, "Authorization header wrong")
	// 	return
	// }

	// requestBody := helpers.ReadRequestJSON[config.User](w, r)
	// newEmailReceievedFromUser := requestBody.Email
	// newPasswordReceivedFromUser := requestBody.Password

	// hash, err := auth.HashPassword(newPasswordReceivedFromUser)
	// if err != nil {
	// 	helpers.RespondWithError(w, http.StatusInternalServerError, "unable to hash password")
	// 	helpers.LogErrorWithRequest(err, r, "unable to hash password")
	// 	return
	// }

	//
}
