// +build unit

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


func SetupDB() { // or *gorm.DB
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = false

	db, err := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string
	if err != nil {
		panic(err)
	}

	dbmgr = &DBManager{
		DB: db,
	}
}

func TestGetImageData(t *testing.T) {
	SetupDB()

    // Create a request to pass to our handler. We don't have any query parameters for now, so we'll
    // pass 'nil' as the third parameter.
    req, err := http.NewRequest("GET", "/image/1", nil)
    if err != nil {
        t.Fatal(err)
    }

	t.Run("NotFound", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getImageData)
	
		handler.ServeHTTP(rr, req)
	
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})
    
	t.Run("Found", func(t *testing.T) {
		commonReply := []map[string]interface{}{{"id":1}}
		mocket.Catcher.Reset().NewMock().WithID(1).WithReply(commonReply)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getImageData)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
}

func TestGetAvailableModels(t *testing.T) {
	req, err := http.NewRequest("GET", "/models", nil)
    if err != nil {
        t.Fatal(err)
	}
	
	t.Run("sdf", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getImageData)
	
		handler.ServeHTTP(rr, req)
	
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
}