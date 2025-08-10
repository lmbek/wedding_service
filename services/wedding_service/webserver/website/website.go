package website

import (
	"net/http"
)

type invitationData struct {
	Name            string
	Information     string
	AcceptedInvites []string
	InviteValid     bool
	InviteCode      string
	InvitePDF       string
}

func (render *Render) FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/front_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	name := "Kære"
	pdf := ""
	valid := false
	if v, ok := render.invites.FindByCode(code); ok {
		valid = true
		name = v.Name
		pdf = v.PDF
	}
	data := invitationData{
		Name:            name,
		Information:     "Din plads er reserveret til ceremonien i Farre Kirke kl. 13:00 og fest i fælleshuset efterfølgende. Kontakt os med sms eller opkald, hvis du har allergier på +45 26 23 25 55.",
		AcceptedInvites: []string{},
		InviteValid:     valid,
		InviteCode:      code,
		InvitePDF:       pdf,
	}
	executePage(w, render, "pages/invitation_page.gohtml", data)
}

func (render *Render) MenuPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/menu_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) RSVPPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/rsvp_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) InfoPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/info_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) WishesPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/wishes_page.gohtml", map[string]string{"Name": "Lars"})
}
