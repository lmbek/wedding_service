package website

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontPageHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	FrontPageHandler(w, r)
	checkStatusOK(t, w)
}

func TestInvitationPageHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/invitation/", nil)

	InvitationPageHandler(w, r)
	checkStatusOK(t, w)
}

func checkStatusOK(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, w.Code)
		return
	}
}
