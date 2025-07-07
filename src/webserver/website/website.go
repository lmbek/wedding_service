package website

import (
	"net/http"
)

func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, "pages/front_page.gohtml", map[string]string{"Name": "Lars"})
}

func InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, "pages/invitation_page.gohtml", map[string]string{"Name": "Lars2"})
}
