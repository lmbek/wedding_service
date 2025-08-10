package database

// Invites defines the interface to manage invitations and RSVP state.
// Interface-first: we never expose the underlying struct.
type Invites interface {
	FindByCode(code string) (Invite, bool)
	ListAccepted(code string) ([]string, error)
	// ListAllAccepted returns all member names that have accepted across all invites.
	ListAllAccepted() ([]string, error)
	Accept(code, name string) error
	Decline(code, name string) error
	EnsureSchema() error
	TrackVisit(code, ip, userAgent, referer, path string) error
}

// Invite holds read-only data for rendering an invitation page.
// Struct is unexported outside the package via returned interface only.
type Invite struct {
	Code    string
	Name    string
	Members []string
}

type invites struct {
	m map[string]Invite
}

// NewInvites returns an empty in-memory implementation (no hardcoded demo data).
// This is used only as a fallback when no database is available.
func NewInvites() Invites {
	data := map[string]Invite{}
	return &invites{m: data}
}

func (i *invites) FindByCode(code string) (Invite, bool) {
	if code == "" {
		return Invite{}, false
	}
	v, ok := i.m[code]
	return v, ok
}

func (i *invites) ListAccepted(code string) ([]string, error) {
	// In-memory fallback keeps no RSVP state to avoid non-persistent storage.
	return []string{}, nil
}

func (i *invites) Accept(code, name string) error {
	// No in-memory mutation. Persist only supported by DB-backed implementation.
	return nil
}

func (i *invites) Decline(code, name string) error {
	// No in-memory mutation. Persist only supported by DB-backed implementation.
	return nil
}

func (i *invites) ListAllAccepted() ([]string, error) {
	// In-memory fallback keeps no RSVP state.
	return []string{}, nil
}

func (i *invites) EnsureSchema() error { return nil }

func (i *invites) TrackVisit(code, ip, userAgent, referer, path string) error {
	// no-op in-memory auditor
	return nil
}
