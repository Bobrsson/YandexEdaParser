package manager

import (
	"context"

	log "github.com/sirupsen/logrus"

	httpServer "YandexEdaParser/http"
	"YandexEdaParser/structs"
)

func ServerRun(config structs.Config) {
	man := new(YandexManager)
	man.Run(config.DB, config.Location, config.Rating)

	var ctx, cancel = context.WithCancel(context.Background())

	if err = httpServer.Run(ctx, man); err != nil {
		log.Println("failed start public http server error: ", err)
		cancel()
	}
}
