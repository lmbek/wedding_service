package website

import (
	"wedding_service/config"
	"wedding_service/webserver/website/frontend"
)

type Render struct {
	config   config.Config
	frontend frontend.Frontend
	//db database.DB // or whatever your DB interface is
}

func NewRender(config config.Config, frontend frontend.Frontend) *Render {
	return &Render{
		config:   config,
		frontend: frontend,
		//db:       db,
	}
}
