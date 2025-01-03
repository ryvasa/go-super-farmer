package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/internal/e2e_test/helper"
	"github.com/ryvasa/go-super-farmer/internal/model/domain"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	baseURL := "http://localhost:8080/api"

	// Mock request body
	userRequest := map[string]string{
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"password": "securepassword",
	}
	requestBody, _ := json.Marshal(userRequest)

	t.Run("should register user successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()
		DBHelper.CreateRole()

		// Kirim request ke server nyata
		resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Verifikasi respons body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Ambil email dari nested map "data"
		data := responseBody["data"].(map[string]interface{})
		email := data["email"].(string)

		error := responseBody["errors"]
		status := responseBody["status"].(float64)
		message := responseBody["message"]
		success := responseBody["success"]
		logrus.Log.Info(error)
		// Assertion
		assert.Nil(t, error)
		assert.Equal(t, 201, int(status))
		assert.Equal(t, "success", message)
		assert.Equal(t, true, success)
		assert.Equal(t, "John Doe", data["name"])
		assert.Equal(t, "john.doe@example.com", email)
	})

	t.Run("should return error if email already exists", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()
		DBHelper.CreateRole()

		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")
		// Kirim request ke server nyata
		resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(requestBody))
		logrus.Log.Info(resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verifikasi respons body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)
		logrus.Log.Info(resp)

		error := responseBody["errors"].(map[string]interface{})
		status := responseBody["status"].(float64)
		message := responseBody["message"]
		success := responseBody["success"]

		code := error["code"]

		assert.NotNil(t, error)
		assert.Equal(t, 400, int(status))
		assert.Equal(t, "failed", message)
		assert.Equal(t, false, success)
		assert.Equal(t, "BAD_REQUEST", code)
		assert.Equal(t, "email already exists", error["message"])
	})

	t.Run("should return error if request body is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()
		DBHelper.CreateRole()

		// Kirim request ke server nyata
		resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer([]byte("invalid request body")))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verifikasi respons body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		error := responseBody["errors"].(map[string]interface{})
		status := responseBody["status"].(float64)
		message := responseBody["message"]
		success := responseBody["success"]

		code := error["code"]

		assert.NotNil(t, error)
		assert.Equal(t, 400, int(status))
		assert.Equal(t, "failed", message)
		assert.Equal(t, false, success)
		assert.Equal(t, "BAD_REQUEST", code)
	})
}

func TestGetAllUsersE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	autHelper := helper.NewAuthHelper(DBHelper)
	baseURL := "http://localhost:8080/api"
	logrus.Log.Info("args ...interface{}")
	t.Run("should get all users successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()
		logrus.Log.Info("haha")

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()
		logrus.Log.Info("haha")

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "Jane Doe", "jane.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parsing respons menggunakan helper generik
		responseBody, err := helper.ParseResponseBodyPagination[domain.User](resp)
		assert.NoError(t, err)

		// Validasi meta respons
		assert.Nil(t, responseBody.Errors)
		assert.Equal(t, 200, responseBody.Status)
		assert.Equal(t, "success", responseBody.Message)
		assert.True(t, responseBody.Success)

		// Validasi data pengguna di dalam pagination
		assert.Equal(t, 1, len(responseBody.Data.Data))
		user := responseBody.Data.Data[0]
		assert.Equal(t, "Jane Doe", user.Name)
		assert.Equal(t, "jane.doe@example.com", user.Email)
	})

	t.Run("should return error if request params is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users?page=invalid", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users", nil)
		assert.NoError(t, err)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return error if forbidden", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenFarmer()
		logrus.Log.Infof("Token: %s", token)

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestGetUserByIDE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	autHelper := helper.NewAuthHelper(DBHelper)
	baseURL := "http://localhost:8080/api"

	t.Run("should get user successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "Jane Doe", "jane.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users/"+id.String(), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parsing respons menggunakan helper generik
		responseBody, err := helper.ParseResponseBody[domain.User](resp)
		assert.NoError(t, err)

		// Validasi meta respons
		assert.Nil(t, responseBody.Errors)
		assert.Equal(t, 200, responseBody.Status)
		assert.Equal(t, "success", responseBody.Message)
		assert.Equal(t, true, responseBody.Success)

		user := responseBody.Data
		assert.Equal(t, "Jane Doe", user.Name)
		assert.Equal(t, "jane.doe@example.com", user.Email)
	})

	t.Run("should return error if request params is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users/invalid", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("GET", baseURL+"/users/"+uuid.New().String(), nil)
		assert.NoError(t, err)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestDeleteUserE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	autHelper := helper.NewAuthHelper(DBHelper)
	baseURL := "http://localhost:8080/api"

	t.Run("should delete user successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("DELETE", baseURL+"/users/"+id.String(), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parsing respons menggunakan helper generik
		responseBody, err := helper.ParseResponseBody[domain.User](resp)
		assert.NoError(t, err)

		// Validasi meta respons
		assert.Nil(t, responseBody.Errors)
		assert.Equal(t, 200, responseBody.Status)
		assert.Equal(t, "success", responseBody.Message)
		assert.Equal(t, true, responseBody.Success)

	})

	t.Run("should return error if request params is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("DELETE", baseURL+"/users/invalid", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("DELETE", baseURL+"/users/"+id.String(), nil)
		assert.NoError(t, err)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return error if forbidden", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenFarmer()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("DELETE", baseURL+"/users/"+id.String(), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestUpdateUserE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	autHelper := helper.NewAuthHelper(DBHelper)
	baseURL := "http://localhost:8080/api"

	t.Run("should update user successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String(), bytes.NewBuffer([]byte(`{"name":"Jane Doe","email":"jane.doe@example.com","password":"securepassword"}`)))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parsing respons menggunakan helper generik
		responseBody, err := helper.ParseResponseBody[domain.User](resp)
		assert.NoError(t, err)

		// Validasi meta respons
		assert.Nil(t, responseBody.Errors)
		assert.Equal(t, 200, responseBody.Status)
		assert.Equal(t, "success", responseBody.Message)
		assert.Equal(t, true, responseBody.Success)

		user := responseBody.Data
		assert.Equal(t, "Jane Doe", user.Name)
		assert.Equal(t, "jane.doe@example.com", user.Email)
	})

	t.Run("should return error if request params is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/invalid", bytes.NewBuffer([]byte(`{"name":"Jane Doe","email":"jane.doe@example.com","password":"securepassword"}`)))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String(), bytes.NewBuffer([]byte(`{"name":"Jane Doe","email":"jane.doe@example.com","password":"securepassword"}`)))
		assert.NoError(t, err)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return error forbiddenif update role with non admin user", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenFarmer()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String(), bytes.NewBuffer([]byte(`{"name":"Jane Doe","email":"jane.doe@example.com","password":"securepassword","role_id":1}`)))
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func TestRestoreUserE2E(t *testing.T) {
	DBHelper := helper.NewDBHelper()

	autHelper := helper.NewAuthHelper(DBHelper)
	baseURL := "http://localhost:8080/api"

	t.Run("should restore user successfully", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")
		DBHelper.DeleteUser(id)

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String()+"/restore", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parsing respons menggunakan helper generik
		responseBody, err := helper.ParseResponseBody[domain.User](resp)
		assert.NoError(t, err)

		// Validasi meta respons
		assert.Nil(t, responseBody.Errors)
		assert.Equal(t, 200, responseBody.Status)
		assert.Equal(t, "success", responseBody.Message)
		assert.Equal(t, true, responseBody.Success)

		user := responseBody.Data
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john.doe@example.com", user.Email)
	})

	t.Run("should return error if request params is invalid", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenAdmin()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")
		DBHelper.DeleteUser(id)

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/invalid/restore", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		DBHelper.DeleteUser(id)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error if unauthorized", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")
		DBHelper.DeleteUser(id)

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String()+"/restore", nil)
		assert.NoError(t, err)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return error if forbidden", func(t *testing.T) {
		defer DBHelper.TeardownTestDB()

		// Dapatkan token dari helper
		token := autHelper.GetTokenFarmer()

		// Setup test data
		id := uuid.New()
		DBHelper.CreateUser(id, "John Doe", "john.doe@example.com", "securepassword")
		DBHelper.DeleteUser(id)

		// Buat request dengan header Authorization
		req, err := http.NewRequest("PATCH", baseURL+"/users/"+id.String()+"/restore", nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		// Kirim request
		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
