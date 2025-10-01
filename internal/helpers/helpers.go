package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/cxt314/drvc-go/internal/config"
)

var app *config.AppConfig

// NewHelpers set up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// func TripToURLValues(t models.Trip) (url.Values, error) {
// 	v := url.Values{}

// 	st := reflect.TypeOf(t)
// 	for i := 0; i < st.NumField(); i++ {
// 		field := st.Field(i)
// 		v.Add(field.Tag.Get("form"), )
// 	}
// }
