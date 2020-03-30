package input

import "github.com/nebiros/krss/internal/model/entity"

type NewFeedInput struct {
	Title string `json:"title" form:"title" validate:"omitempty"`
	URL   string `json:"url" form:"url" validate:"required,url"`
}

func (in *NewFeedInput) ToCreateFeed() entity.CreateFeed {
	return entity.CreateFeed{
		UserID: -1,
		Title:  in.Title,
		URL:    in.URL,
	}
}
