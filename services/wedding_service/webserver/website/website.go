package website

import (
	"net/http"
	"strings"
)

type invitationData struct {
	Name            string
	Information     string
	AcceptedInvites []string
	InviteValid     bool
	InviteCode      string
}

func (render *Render) FrontPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/front_page.gohtml", map[string]string{"Name": "Lars"})
}

func (render *Render) InvitationPageHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	name := ""
	valid := false
	accepted := []string{}
	if v, ok := render.invites.FindByCode(code); ok {
		valid = true
		name = v.Name
		// Best-effort preload of accepted list; frontend will refresh via API
		acc, err := render.invites.ListAccepted(code)
		if err == nil {
			accepted = acc
		}
	}
	// Audit: log that an invite link was accessed
	ip := r.RemoteAddr
	ua := r.Header.Get("User-Agent")
	ref := r.Header.Get("Referer")
	path := r.URL.Path
	_ = render.invites.TrackVisit(code, ip, ua, ref, path)
	data := invitationData{
		Name:            name,
		Information:     "Din plads er reserveret til ceremonien i Farre Kirke kl. 13:00 og fest i fælleshuset efterfølgende. Kontakt os med sms eller opkald, hvis du har allergier på +45 26 23 25 55.",
		AcceptedInvites: accepted,
		InviteValid:     valid,
		InviteCode:      code,
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

func (render *Render) ReceptionPageHandler(w http.ResponseWriter, r *http.Request) {
	executePage(w, render, "pages/reception_page.gohtml", map[string]string{"Name": "Lars"})
}
