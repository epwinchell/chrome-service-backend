package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RedHatInsights/chrome-service-backend/rest/models"
	"github.com/RedHatInsights/chrome-service-backend/rest/service"
	"github.com/RedHatInsights/chrome-service-backend/rest/util"
	"github.com/go-chi/chi/v5"
)

type AddVisitedBundlePayload struct {
	Bundle string `json:"bundle"`
}

// Use the user obj in context to pull full data row from DB
func GetUserIdentity(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(util.USER_CTX_KEY).(models.UserIdentity)
	updatedUser, err := service.GetUserIdentityData(user)
	if err != nil {
		panic(err)
	}

	resp := util.EntityResponse[models.UserIdentity]{
		Data: updatedUser,
	}

	json.NewEncoder(w).Encode(resp)
}

func AddVisitedBundle(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(util.USER_CTX_KEY).(models.UserIdentity)
	var request AddVisitedBundlePayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}
	updatedUser, err := service.AddVisitedBundle(user, request.Bundle)
	if err != nil {
		panic(err)
	}

	resp := util.EntityResponse[models.UserIdentity]{
		Data: updatedUser,
	}

	json.NewEncoder(w).Encode(resp)
}

func GetVisitedBundles(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(util.USER_CTX_KEY).(models.UserIdentity)
	bundle, err := service.GetVisitedBundles(user)
	if err != nil {
		panic(err)
	}
	resp := util.EntityResponse[map[string]bool]{
		Data: bundle,
	}

	json.NewEncoder(w).Encode(resp)
}

func GetIntercomHash(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(util.USER_CTX_KEY).(models.UserIdentity)
	bundleParam := r.URL.Query()["bundle"]
	bundle := "fallback"

	if len(bundleParam) > 0 {
		bundle = bundleParam[0]
	}
	hash, err := service.GetUserIntercomHash(user.AccountId, service.IntercomBundle(bundle))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error."))
		return
	}

	resp := util.EntityResponse[string]{
		Data: hash,
	}

	json.NewEncoder(w).Encode(resp)
}

func MakeUserIdentityRoutes(sub chi.Router) {
	sub.Get("/", GetUserIdentity)
	sub.Get("/intercom", GetIntercomHash)
	sub.Route("/visited-bundles", func(r chi.Router) {
		r.Post("/", AddVisitedBundle)
		r.Get("/", GetVisitedBundles)
	})
}
