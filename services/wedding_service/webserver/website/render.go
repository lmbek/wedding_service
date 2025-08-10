package website

import (
	"wedding_service/config"
	"wedding_service/webserver/database"
	"wedding_service/webserver/website/frontend"
)

type Render struct {
	config   config.Config
	frontend frontend.Frontend
	invites  database.Invites
}

func NewRender(config config.Config, frontend frontend.Frontend, invites database.Invites) *Render {
	return &Render{
		config:   config,
		frontend: frontend,
		invites:  invites,
	}
}
