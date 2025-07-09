package website

import "wedding_service/webserver/website/frontend"

type Render struct {
	frontend frontend.Frontend
	//db       database.DB // or whatever your DB interface is
}

func NewRender(frontend frontend.Frontend) *Render {
	return &Render{
		frontend: frontend,
		//db:       db,
	}
}
