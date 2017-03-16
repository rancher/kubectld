package events

import (
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/client"
)

func create(event *events.Event, cli *client.RancherClient) error {
	return wrap(event, cli, installCatalog)
}

func upgrade(event *events.Event, cli *client.RancherClient) error {
	return wrap(event, cli, upgradeCatalog)
}

func rollback(event *events.Event, cli *client.RancherClient) error {
	return wrap(event, cli, rollbackCatalog)
}

func remove(event *events.Event, cli *client.RancherClient) error {
	return wrap(event, cli, removeCatalog)
}

func finishUpgrade(event *events.Event, cli *client.RancherClient) error {
	return createAndPublishReply(event, cli)
}

func ping(event *events.Event, cli *client.RancherClient) error {
	return createAndPublishReply(event, cli)
}
