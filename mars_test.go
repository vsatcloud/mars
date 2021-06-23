package mars

import (
	"fmt"
	"net/http"
	"testing"

	"gorm.io/gorm"

	"github.com/vsatcloud/mars/casbin"
)

func TestMars(t *testing.T) {
	mars := New()

	var db = &gorm.DB{}

	casbin.Casbin(db)

	server := &http.Server{
		Addr:    "localhost:9999",
		Handler: mars,
	}

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Server closed under request")
		} else {
			panic(err)
		}
	}

}
