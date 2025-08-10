package database

// Invites defines the interface to manage invitations and RSVP state.
// Interface-first: we never expose the underlying struct.
type Invites interface {
	FindByCode(code string) (Invite, bool)
	ListAccepted(code string) ([]string, error)
	Accept(code, name string) error
	Decline(code, name string) error
	EnsureSchema() error
}

// Invite holds read-only data for rendering an invitation page.
// Struct is unexported outside the package via returned interface only.
type Invite struct {
	Code    string
	Name    string
	PDF     string
	Members []string
}

type invites struct {
	m        map[string]Invite
	accepted map[string]map[string]struct{}
}

// NewInvites returns an in-memory implementation with a few demo codes.
func NewInvites() Invites {
	data := map[string]Invite{
		// Demo/test codes. You kan replace/add as needed.
		"jensen":     {Code: "jensen", Name: "Familien Jensen", Members: []string{"Rita", "Hans"}},
		"bek":        {Code: "bek", Name: "Familien Bek", Members: []string{"Johnny", "Cecilie"}},
		"hfj934e":    {Code: "hfj934e", Name: "Familien test3", Members: []string{"Peter"}},
		"child0001":  {Code: "child0001", Name: "Sofie (barn)", Members: []string{"Sofie (barn)"}},
		"guest-demo": {Code: "guest-demo", Name: "Gæst", Members: []string{"Gæst"}},
	}
	return &invites{m: data, accepted: make(map[string]map[string]struct{})}
}

func (i *invites) FindByCode(code string) (Invite, bool) {
	if code == "" {
		return Invite{}, false
	}
	v, ok := i.m[code]
	return v, ok
}

func (i *invites) ListAccepted(code string) ([]string, error) {
	st := i.accepted[code]
	if st == nil {
		return []string{}, nil
	}
	acc := make([]string, 0, len(st))
	for name := range st {
		acc = append(acc, name)
	}
	return acc, nil
}

func (i *invites) Accept(code, name string) error {
	if code == "" || name == "" {
		return nil
	}
	if _, ok := i.accepted[code]; !ok {
		i.accepted[code] = make(map[string]struct{})
	}
	i.accepted[code][name] = struct{}{}
	return nil
}

func (i *invites) Decline(code, name string) error {
	if code == "" || name == "" {
		return nil
	}
	if st, ok := i.accepted[code]; ok {
		delete(st, name)
	}
	return nil
}

func (i *invites) EnsureSchema() error { return nil }
