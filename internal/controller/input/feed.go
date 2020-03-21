package input

import "github.com/nebiros/krss/internal/model/entity"

type AddFeedInput struct {
	Title string `json:"title" form:"title" validate:"required"`
	URL   string `json:"url" form:"url" validate:"required,url"`
}

func (in *AddFeedInput) ToCreateFeed() entity.CreateFeed {
	return entity.CreateFeed{
		UserID: 0,
		Title:  "",
		URL:    "",
	}
}
