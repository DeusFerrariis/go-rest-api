package main

import (
// "bytes"
// "encoding/json"
// "net/http"
// "net/http/httptest"
// "testing"
)

// func TestHandleCreateUser(t *testing.T) {
// 	body := map[string]string{"username": "wcleg4"}

// 	var b bytes.Buffer
// 	if err := json.NewEncoder(&b).Encode(body); err != nil {
// 		t.Fatalf("error serializing request body %s", err.Error())
// 	}

// 	req := httptest.NewRequest("POST", "/user/new", &b)
// 	wr := httptest.NewRecorder()
// 	users := make([]User, 0)
// 	wu := WithUsers(&users)

// 	wu(wr, req, HandleCreateUser)

// 	if wr.Code != http.StatusAccepted {
// 		t.Errorf("expected statuscode 202, got %d", wr.Code)
// 	}
// 	if len(users) != 1 {
// 		t.Errorf("expected 1 new user, have %d", len(users))
// 	}
// 	// TODO: add id check
// }
