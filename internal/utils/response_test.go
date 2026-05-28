package utils

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestErrorResponse(t *testing.T) {
    gin.SetMode(gin.TestMode)

    recorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(recorder)

    ErrorResponse(c, http.StatusBadRequest, "invalid input")

    if recorder.Code != http.StatusBadRequest {
        t.Errorf("expected status %d, got %d",
            http.StatusBadRequest,
            recorder.Code,
        )
    }

    var response map[string]interface{}

    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }

    if response["error"] != "invalid input" {
        t.Errorf("unexpected error message")
    }

    if int(response["code"].(float64)) != http.StatusBadRequest {
        t.Errorf("unexpected status code")
    }
}

func TestSuccessResponse(t *testing.T) {
    gin.SetMode(gin.TestMode)

    recorder := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(recorder)

    payload := gin.H{
        "message": "success",
    }

    SuccessResponse(c, http.StatusOK, payload)

    if recorder.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d",
            http.StatusOK,
            recorder.Code,
        )
    }

    var response map[string]interface{}

    err := json.Unmarshal(recorder.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }

    data := response["data"].(map[string]interface{})

    if data["message"] != "success" {
        t.Errorf("unexpected success message")
    }

    if int(response["code"].(float64)) != http.StatusOK {
        t.Errorf("unexpected status code")
    }
}