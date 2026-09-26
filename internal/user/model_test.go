package user

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {
	user := User{
		ID:           1,
		Username:     "alice",
		PasswordHash: "secret-hash",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	jsonString := string(jsonData)
	if !strings.Contains(jsonString, "alice") {
		t.Errorf("JSON does not contain username: %s", jsonString)
	}
	if strings.Contains(jsonString, "secret-hash") {
		t.Errorf("JSON contains password hash: %s", jsonString)
	}
	if strings.Contains(jsonString, "password_hash") {
		t.Errorf("JSON contains password_hash field: %s", jsonString)
	}
}
