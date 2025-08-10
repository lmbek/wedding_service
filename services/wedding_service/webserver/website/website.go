package website

import (
	"net/http"
)

func (render *Render) FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/front_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/invitation_page.gohtml", map[string]string{"Name": "Lars2"})
}
