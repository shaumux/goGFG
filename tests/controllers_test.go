package tests

import (
	"fmt"
	"goGFG/controllers"
	"goGFG/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAuthFail(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.GetProducts)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusUnauthorized)
	}

}

func TestAuthSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.GetProducts)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}
}

func TestGetProduct(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.GetProducts)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	expected := `[]`
	if strings.Compare(strings.TrimSpace(rr.Body.String()), expected) != 0 {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIndexProduct(t *testing.T) {
	body := `{"title":"shoe","brand":"roadster","price":1000,"stock":10}`
	req, err := http.NewRequest("POST", "/api/v1/products", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.IndexProduct)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(body) {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr.Body.String(), body)
	}
	time.Sleep(15 * time.Second)
	req1, err := http.NewRequest("GET", "http:///api/v1/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	req1.SetBasicAuth("shaumux", "secretPassword")
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	handlerFunc1 := http.HandlerFunc(controllers.GetProducts)
	handler1 := middlewares.BasicAuthMiddleWare(handlerFunc1)
	handler1.ServeHTTP(rr1, req1)
	if status := rr1.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf("[%s]", body)
	if strings.TrimSpace(rr1.Body.String()) != expected {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr1.Body.String(), expected)
	}
}

func TestFilter(t *testing.T) {
	body := `{"title":"slippers","brand":"indoors","price":1000,"stock":10}`

	req, err := http.NewRequest("POST", "/api/v1/products", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.IndexProduct)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(body) {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr.Body.String(), body)
	}

	time.Sleep(15 * time.Second)

	req, err = http.NewRequest("GET", "/api/v1/products?filter=brand:indoors", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handlerFunc = http.HandlerFunc(controllers.GetProducts)
	handler = middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	expected := fmt.Sprintf("[%s]", body)
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSearch(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/products?q=slippers", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("shaumux", "secretPassword")
	rr := httptest.NewRecorder()
	handlerFunc := http.HandlerFunc(controllers.GetProducts)
	handler := middlewares.BasicAuthMiddleWare(handlerFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Incorrect Status Code, Got: %v expected: %v", status, http.StatusOK)
	}

	expected := `[{"title":"slippers","brand":"indoors","price":1000,"stock":10}]`
	if strings.Compare(strings.TrimSpace(rr.Body.String()), expected) != 0 {
		t.Errorf("handler returned unexpected body, got %v want %v",
			rr.Body.String(), expected)
	}
}
