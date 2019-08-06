package impl

import (
	"github.com/mylxsw/glacier"
	"github.com/mylxsw/go-toolkit/container"
)

type ServiceProvider struct{}

func (s ServiceProvider) Register(app *container.Container) {
	app.MustSingleton(NewSequenceRepo)
	app.MustSingleton(NewKVRepo)
	app.MustSingleton(NewMessageRepo)
	app.MustSingleton(NewMessageGroupRepo)
	app.MustSingleton(NewUserRepo)
	app.MustSingleton(NewRuleRepo)
}

func (s ServiceProvider) Boot(app *glacier.Glacier) {}