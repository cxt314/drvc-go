package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cxt314/drvc-go/internal/forms"
	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

// UserLogin shows the user login page
func (m *Repository) UserLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

// UserLoginPost handles logging the user in
func (m *Repository) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)

	// run form validations
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

// UserList lists all users
func (m *Repository) UserList(w http.ResponseWriter, r *http.Request) {
	users, err := m.DB.AllUsers()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["users"] = users

	render.Template(w, r, "user-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// UserEditIndex redirects to the currently logged-in user's edit page
func (m *Repository) UserEditIndex(w http.ResponseWriter, r *http.Request) {
	id := m.App.Session.GetInt(r.Context(), "user_id")
	http.Redirect(w, r, fmt.Sprintf("/users/update/%d", id), http.StatusSeeOther)
}

// UserEdit displays form to update user
func (m *Repository) UserEdit(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get user from database
	v, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["user"] = v

	render.Template(w, r, "edit-user.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// UserEditPost updates user info from submitted form
func (m *Repository) UserEditPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get user from database
	v, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form into fetched user
	err = helpers.ParseFormToUser(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks
	form.Required("email", "first-name", "last-name")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = v

		render.Template(w, r, "edit-user.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateUser(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Updated user successfully")
	http.Redirect(w, r, fmt.Sprintf("/users/update/%d", id), http.StatusSeeOther)
}

// UserCreate displays form to create a new user
func (m *Repository) UserCreate(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-user.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// UserCreatePost creates a new user
func (m *Repository) UserCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse received form values into user object
	v := models.User{}
	err := helpers.ParseFormToUser(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks
	form.Required("email", "password", "first-name", "last-name")
	form.IsEmail("email")
	form.MinLength("password", 8, r)
	form.IsEqual("password", "password-confirm")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = v

		render.Template(w, r, "edit-user.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertUser(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UserEditPassword displays form to update user password
func (m *Repository) UserEditPassword(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get user from database
	v, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["user"] = v

	render.Template(w, r, "edit-user-pw.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// UserEditPasswordPost updates user password
func (m *Repository) UserEditPasswordPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get user from database
	v, err := m.DB.GetUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks
	form.Required("current-password", "password", "password-confirm")
	_, _, err = m.DB.Authenticate(v.Email, r.Form.Get("current-password"))
	if err != nil {
		log.Println(err)
		form.Errors.Add("current-password", "Invalid password")
	}
	form.MinLength("password", 8, r)
	form.IsEqual("password", "password-confirm")


	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = v

		render.Template(w, r, "edit-user-pw.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// update user with new password from form
	v.Password = r.Form.Get("password")
	
	err = m.DB.UpdateUserPassword(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Updated password successfully")
	http.Redirect(w, r, fmt.Sprintf("/users/update/%d", id), http.StatusSeeOther)
}

func (m *Repository) UserDelete(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if id == userID {
		// don't allow deleting current user
		m.App.Session.Put(r.Context(), "error", "Cannot delete current user")
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	err = m.DB.DeleteUserByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "User deleted")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}