package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns true if there are no errors, otherwise it returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in POST and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	// if x == "" {
	// 	return false
	// }
	// return true
	return x != ""
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for a valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) IsValidEndMileage(field string, originalStartMileage int, originalEndMileage int, nextTripEndMileage int) {
	newEndMileage, err := strconv.Atoi(f.Get(field))
	if err != nil {
		f.Errors.Add(field, "This field must be a number")
	}

	if newEndMileage < originalStartMileage {
		f.Errors.Add(field, "End mileage cannot be less than start mileage")
	}

	diff := newEndMileage - originalEndMileage
	if diff == 1000 {
		// no error, all future trips will have mileages adjusted by 1000
	} else if diff > 0 && diff < 1000 {
		if newEndMileage > nextTripEndMileage {
			f.Errors.Add(field, fmt.Sprintf("End Mileage must be less than next trip's end mileage: %d", nextTripEndMileage))
		}
	}
}
