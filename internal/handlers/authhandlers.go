package handlers

import (
	"net/http"

	"github.com/cxt314/drvc-go/internal/forms"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

func (m *Repository) UserLogin(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: forms.New(nil)})
}

func (m *Repository) UserLoginPost(w http.ResponseWriter, r *http.Request) {

}
