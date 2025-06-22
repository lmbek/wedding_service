package website

import (
	_ "embed"
	"html/template"
	"net/http"
)

//go:embed frontend/out/private/pages/front_page.gohtml
var frontPageData string

//go:embed frontend/out/private/pages/invitation_page.gohtml
var invitationPageData string

func FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	ExecutePage(
		w,
		template.New("frontpage"),
		frontPageData,
		map[string]string{"Name": "Lars"},
	)
}

func InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	ExecutePage(
		w,
		template.New("invitation"),
		invitationPageData,
		map[string]string{"Name": "Lars2"},
	)
}
