package broker

import (
	"go.uber.org/fx"
)


func NewApp() *fx.App {
	app := fx.New(
		fx.Provide(
			newConfig,
			newLogger,
			newStanClient,
			newBroker,
		),
	)

	return app
}
