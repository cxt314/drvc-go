package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cxt314/drvc-go/internal/forms"
	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

// MemberList displays a list of all members
func (m *Repository) MemberList(w http.ResponseWriter, r *http.Request) {
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["members"] = members

	render.Template(w, r, "member-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// MemberCreate displays the page to create a new member
func (m *Repository) MemberCreate(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// MemberCreatePost processes the POST request for creating a new member
func (m *Repository) MemberCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse received form values into member object
	v := models.Member{}
	err := helpers.ParseFormToMember(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["member"] = v

		render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertMember(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/members", http.StatusSeeOther)
}

// MemberEdit shows the edit form for a member by id
func (m *Repository) MemberEdit(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get member from database
	v, err := m.DB.GetMemberByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["member"] = v

	render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// MemberEditPost processes the POST request for editing a member by id
func (m *Repository) MemberEditPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get member from database
	v, err := m.DB.GetMemberByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form into fetched member
	err = helpers.ParseFormToMember(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["member"] = v

		render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateMember(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Updated member successfully")
	http.Redirect(w, r, fmt.Sprintf("/members/%d", id), http.StatusSeeOther)
}

// MemberDeactivate updates a member's active status by id
func (m *Repository) MemberDeactivate(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// update member set is_active to false
	err = m.DB.UpdateMemberActiveByID(id, false)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// redirect to members list after making member inactive
	http.Redirect(w, r, "/members", http.StatusSeeOther)
}

func (m *Repository) AddAlias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
		<div class="col-3 mb-2 alias-group">
			<div class="input-group">
				<input class="form-control" type="text" name="aliases">
				<span class="input-group-text">
					<button type="button" class="btn-close" aria-label="Close"
						hx-get="/remove-item" hx-trigger="click" hx-target="closest .alias-group" hx-swap="outerHTML" hx-confirm="Delete alias?"></button>
				</span>
			</div>
		</div>
		`

	w.Write([]byte(html))
}

/*// MemberDelete deletes the member with the given id
func (m *Repository) MemberDelete(w http.ResponseWriter, r *http.Request) {
	// exploded := strings.Split(r.RequestURI, "/")
	// id, err := strconv.Atoi(exploded[2])
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	render.Template(w, r, "member-list.page.tmpl", &models.TemplateData{})
}*/
