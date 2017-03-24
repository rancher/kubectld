package events

import (
	revents "github.com/rancher/event-subscriber/events"
)

func StartEventHandler(url, accessKey, secretKey string, workers int) error {
	eventHandlers := map[string]revents.EventHandler{
		"kubernetesStack.create":        create,
		"kubernetesStack.upgrade":       upgrade,
		"kubernetesStack.rollback":      rollback,
		"kubernetesStack.remove":        remove,
		"kubernetesStack.finishupgrade": finishUpgrade,
		"ping": ping,
	}

	router, err := revents.NewEventRouter("", 0, url, accessKey, secretKey, nil, eventHandlers, "", workers, revents.DefaultPingConfig)
	if err != nil {
		return err
	}

	return router.StartWithoutCreate(nil)
}
