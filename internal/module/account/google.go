package account

import (
	"github.com/gin-gonic/gin"
)


type _google struct {
	Email     string `json:"email"`
	Id        string `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Url       string `json:"url"`
	Name      string `json:"name"`
}

func NewGoogle(ctx *gin.Context) (g *_google, err error) {
	err = ctx.ShouldBindJSON(&g)
	return
}

func (g *_google) GetName() string {
	return g.Name
}

func (g *_google) GetId() string {
	return g.Id
}

func (g *_google) GetUrl() string {
	return g.Url
}

func (g *_google) GetEmail() string {
	return g.Email
}

func (g *_google) GetAvatarUrl() string {
	return g.AvatarUrl
}

func (g *_google) Source() string {
	return "google"
}
