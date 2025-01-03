package helper

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ryvasa/go-super-farmer/pkg/logrus"
)

type AuthHelper interface {
	GetTokenAdmin() string
	GetTokenFarmer() string
}

type AuthHelperImpl struct {
	dbHelper DBHelper
}

func NewAuthHelper(dbHelper DBHelper) AuthHelper {
	return &AuthHelperImpl{dbHelper}
}

func (e *AuthHelperImpl) GetTokenAdmin() string {
	e.dbHelper.CreateRole()
	baseURL := "http://localhost:8080/api"

	// Mock request body
	userRequest := map[string]string{
		"name":     "Test Admin Helper",
		"email":    "testadminhelper@example.com",
		"password": "securepassword",
	}
	loginRequest := map[string]string{
		"email":    "testadminhelper@example.com",
		"password": "securepassword",
	}
	userBody, _ := json.Marshal(userRequest)
	authBody, _ := json.Marshal(loginRequest)

	// Kirim request ke server nyata
	_, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(userBody))
	if err != nil {
		panic(err)
	}
	auth, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(authBody))
	if err != nil {
		panic(err)
	}

	// Verifikasi respons body
	var responseBody map[string]interface{}
	json.NewDecoder(auth.Body).Decode(&responseBody)

	// Ambil email dari nested map "data"
	data := responseBody["data"].(map[string]interface{})
	return data["token"].(string)
}

func (e *AuthHelperImpl) GetTokenFarmer() string {
	e.dbHelper.CreateRole()

	baseURL := "http://localhost:8080/api"

	// Mock request body
	userRequest := map[string]string{
		"name":     "Test Farmer Helper",
		"email":    "testfarmerhelper@example.com",
		"password": "securepassword",
	}
	loginRequest := map[string]string{
		"email":    "testfarmerhelper@example.com",
		"password": "securepassword",
	}
	userBody, _ := json.Marshal(userRequest)
	authBody, _ := json.Marshal(loginRequest)

	// Kirim request ke server nyata
	user, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(userBody))
	if err != nil {
		panic(err)
	}
	defer user.Body.Close()

	// Parse respons body untuk user
	var userResponse map[string]interface{}
	if err := json.NewDecoder(user.Body).Decode(&userResponse); err != nil {
		panic(err)
	}

	// Ambil data dari respons
	data, ok := userResponse["data"].(map[string]interface{})
	if !ok {
		panic("invalid response structure, data is not a map")
	}

	id, ok := data["id"].(string)
	if !ok {
		panic("invalid response structure, id is not a string")
	}

	token := e.GetTokenAdmin()

	// Kirim request ke server nyata
	UpdateUserRole(baseURL, id, token)

	// Login untuk mendapatkan token
	auth, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(authBody))
	if err != nil {
		panic(err)
	}
	defer auth.Body.Close()

	// Parse respons body untuk login
	var responseBody map[string]interface{}
	if err := json.NewDecoder(auth.Body).Decode(&responseBody); err != nil {
		panic(err)
	}

	// Ambil token dari respons login
	authData, ok := responseBody["data"].(map[string]interface{})
	if !ok {
		panic("invalid response structure, data is not a map")
	}

	token, ok = authData["token"].(string)
	if !ok {
		panic("invalid response structure, token is not a string")
	}

	return token
}

func UpdateUserRole(baseURL, id, token string) {
	logrus.Log.Infof("Update user role with ID: %s", id)
	// Body untuk update request
	updateBody := map[string]interface{}{
		"role_id": 2,
	}
	updateBodyBytes, _ := json.Marshal(updateBody)

	// Buat request PATCH
	req, err := http.NewRequest("PATCH", baseURL+"/users/"+id, bytes.NewBuffer(updateBodyBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+token) // Tambahkan header Authorization
	req.Header.Set("Content-Type", "application/json")

	logrus.Log.Infof("Request body: %+v", updateBody)
	// Kirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Verifikasi respons
	if resp.StatusCode != http.StatusOK {
		logrus.Log.Info(resp.StatusCode)
		panic("failed to update user role, status: " + resp.Status)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		panic(err)
	}

	logrus.Log.Infof("Update successful: %+v", responseBody)
}
