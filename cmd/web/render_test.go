package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/some-url", nil) // build a request to test
	ctx := getCtx(req)                                 // add session to request context
	req = req.WithContext(ctx)

	// add flash, warning, and error messages to session
	testApp.Session.Put(ctx, "flash", "flash message")
	testApp.Session.Put(ctx, "warning", "warning message")
	testApp.Session.Put(ctx, "error", "error message")

	// construct a render template
	td := testApp.AddDefaultData(&TemplateData{}, req)

	// check if the session values are set in the template data
	if td.Flash != "flash message" {
		t.Error("flash value not set in template data")
	}
	if td.Warning != "warning message" {
		t.Error("warning value not set in template data")
	}
	if td.Error != "error message" {
		t.Error("error value not set in template data")
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/some-url", nil) // build a request to test
	ctx := getCtx(req)                                 // add session to request context
	req = req.WithContext(ctx)

	// Case 1 - session is not authenticated
	auth := testApp.IsAuthenticated(req)
	if auth {
		t.Error("expected false, got true - session should not be authenticated")
	}

	// Case 2 - session is authenticated
	testApp.Session.Put(ctx, "userID", 37)
	auth = testApp.IsAuthenticated(req)
	if !auth {
		t.Error("expected true, got false - session should be authenticated")
	}
}

func TestConfig_render(t *testing.T) {
	req, _ := http.NewRequest("GET", "/some-url", nil) // build a request to test
	ctx := getCtx(req)                                 // add session to request context
	req = req.WithContext(ctx)

	res := httptest.NewRecorder() // create a response recorder

	// render
	testApp.render(res, req, "home.page.gohtml", &TemplateData{})

	// check if the response status code is 200
	if res.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.Code)
	}
}
