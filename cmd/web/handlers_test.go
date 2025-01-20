package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/chriskoorzen/go-subscription-webapp/cmd/web/db"
)

var pageTests = []struct {
	testName           string
	url                string
	httpVerb           string
	handler            http.HandlerFunc
	sessionData        map[string]interface{}
	expectedStatusCode int
	expectedHTML       string
}{
	{
		testName:           "home page",
		url:                "/",
		httpVerb:           "GET",
		handler:            testApp.GETHomePage,
		expectedStatusCode: http.StatusOK,
	},
	{
		testName:           "login page",
		url:                "/login",
		httpVerb:           "GET",
		handler:            testApp.GETLoginPage,
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `<h1 class="mt-5">Login</h1>`,
	},
	{
		testName: "logout page",
		url:      "/logout",
		httpVerb: "GET",
		handler:  testApp.GETLogout,
		sessionData: map[string]interface{}{
			"userID": 1,
			"user":   db.User{ID: 1, Active: 1},
		},
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		testName:           "register page",
		url:                "/register",
		httpVerb:           "GET",
		handler:            testApp.GETRegisterPage,
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `<h1 class="mt-5">Register</h1>`,
	},
	// { // TODO test token validity
	// 	testName:           "activate account page",
	// 	url:                "/activate-account",
	// 	httpVerb:           "GET",
	// 	handler:            testApp.GETActivateAccount,
	// 	expectedStatusCode: http.StatusSeeOther,
	// },
	{
		testName: "subscription plans page",
		url:      "/members/plans",
		httpVerb: "GET",
		handler:  testApp.GETSubscriptionPlans,
		sessionData: map[string]interface{}{
			"userID": 1,
			"user":   db.User{ID: 1, Active: 1},
		},
		expectedStatusCode: http.StatusOK,
		expectedHTML:       `<h1 class="mt-5">Plans</h1>`,
	},
}

func Test_Pages(t *testing.T) {
	for _, e := range pageTests {
		req, _ := http.NewRequest(e.httpVerb, e.url, nil) // build a request to test
		ctx := getCtx(req)                                // add session to request context
		req = req.WithContext(ctx)
		res := httptest.NewRecorder() // create a response recorder

		// add session data to the request context
		if len(e.sessionData) > 0 {
			for k, v := range e.sessionData {
				testApp.Session.Put(ctx, k, v)
			}
		}

		// execute the handler
		e.handler.ServeHTTP(res, req)

		// test results
		if res.Code != e.expectedStatusCode {
			t.Errorf("%s failed - expected status %d, got %d", e.testName, e.expectedStatusCode, res.Code)
		}
		// we are looking for html
		if len(e.expectedHTML) > 0 {
			html := res.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("%s failed - expected html not found:\n%s", e.testName, e.expectedHTML)
			}
		}
	}

}

func TestConfig_POSTLoginPage(t *testing.T) {

	postedData := strings.NewReader(url.Values{
		"email":    {"testman@example.com"},
		"password": {"abc12345"},
	}.Encode())

	req, _ := http.NewRequest("POST", "/login", postedData) // build a request to test
	ctx := getCtx(req)                                      // add session to request context
	req = req.WithContext(ctx)
	res := httptest.NewRecorder() // create a response recorder

	handler := http.HandlerFunc(testApp.POSTLoginPage)
	handler.ServeHTTP(res, req)

	// test results
	if res.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", res.Code)
	}
	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("did not find userID in session")
	}
}

func TestConfig_GETSubscribeToPlan(t *testing.T) {
	req, _ := http.NewRequest("GET", "/members/subscribe?plan=1", nil) // build a request to test
	ctx := getCtx(req)                                                 // add session to request context
	req = req.WithContext(ctx)
	res := httptest.NewRecorder() // create a response recorder

	// add session data to the request context
	testApp.Session.Put(ctx, "user", db.User{
		ID:        1,
		Active:    1,
		Email:     "testUser@example.com",
		FirstName: "Test",
		LastName:  "User",
	})

	handler := http.HandlerFunc(testApp.GETSubscribeToPlan)
	handler.ServeHTTP(res, req)

	// test results
	if res.Code != http.StatusSeeOther {
		t.Errorf("expected status 303, got %d", res.Code)
	}
}
