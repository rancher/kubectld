package events

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubectld/events/util"
	"github.com/rancher/kubectld/stack"
)

type EventHandlerWithData func(*revents.Event, *client.RancherClient) (map[string]interface{}, error)

type StackHandler struct {
	Server string
}

func (h *StackHandler) Create(event *revents.Event, cli *client.RancherClient) error {
	return h.wrap(event, cli, h.create)
}

func (h *StackHandler) FinishUpgrade(event *revents.Event, cli *client.RancherClient) error {
	return h.wrap(event, cli, h.finishUpgrade)
}

func (h *StackHandler) Rollback(event *revents.Event, cli *client.RancherClient) error {
	return h.wrap(event, cli, h.rollback)
}

func (h *StackHandler) Upgrade(event *revents.Event, cli *client.RancherClient) error {
	return h.wrap(event, cli, h.upgrade)
}

func (h *StackHandler) Remove(event *revents.Event, cli *client.RancherClient) error {
	logrus.Infof("Received event: Name: %s, Event Id: %s, Resource Id: %s", event.Name, event.ID, event.ResourceID)
	if err := h.remove(event, cli); err != nil {
		logrus.WithField("EventId", event.ID).Errorf("Failed to delete: %v", err)
	}
	return util.PublishReply(util.NewReply(event), cli)
}

func (h *StackHandler) wrap(event *revents.Event, cli *client.RancherClient, fn EventHandlerWithData) error {
	logrus.Infof("Received event: Name: %s, Event Id: %s, Resource Id: %s", event.Name, event.ID, event.ResourceID)
	data, err := fn(event, cli)
	resp := util.NewReply(event)
	if err == nil {
		resp.Data = data
	} else {
		resp.TransitioningMessage = err.Error()
		resp.Transitioning = "error"
	}
	return util.PublishReply(resp, cli)
}

func (h *StackHandler) getInput(event *revents.Event, cli *client.RancherClient) stack.Input {
	uuid := util.GetString(event.Data, "environment", "uuid")
	templates := util.GetStringMap(event.Data, "environment", "data", "fields", "templates")
	environment := util.GetStringMap(event.Data, "environment", "data", "fields", "environment")
	namespace := util.GetString(event.Data, "environment", "data", "fields", "namespace")

	if environment == nil {
		environment = map[string]string{}
	}

	return stack.Input{
		Server:      h.Server,
		Files:       templates,
		Environment: environment,
		Namespace:   namespace,
		Labels: map[string]string{
			"io.rancher.stack.uuid": uuid,
		},
	}
}

func (h *StackHandler) create(event *revents.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	input := h.getInput(event, cli)
	if len(input.Files) == 0 {
		return nil, nil
	}
	err := stack.Create(input)
	return nil, err
}

func (h *StackHandler) rollback(event *revents.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	_, err := h.create(event, cli)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"environment": map[string]interface{}{
			"upgrade": nil,
		},
	}, nil
}

func (h *StackHandler) remove(event *revents.Event, cli *client.RancherClient) error {
	input := h.getInput(event, cli)
	if len(input.Files) == 0 {
		return nil
	}
	return stack.Remove(input)
}

func (h *StackHandler) finishUpgrade(event *revents.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	fmt.Printf("Finish Upgrade: %#v\n", event.Data)
	upgradeData := util.GetMap(event.Data, "environment", "data", "fields", "upgrade")
	if upgradeData == nil {
		return nil, nil
	}
	upgradeData["upgrade"] = nil

	return map[string]interface{}{
		"environment": upgradeData,
	}, nil
}

func (h *StackHandler) upgrade(event *revents.Event, cli *client.RancherClient) (map[string]interface{}, error) {
	templates := util.GetStringMap(event.Data, "processData", "templates")
	environment := util.GetStringMap(event.Data, "processData", "environment")
	externalID := util.GetString(event.Data, "processData", "externalId")
	previousExternalID := util.GetString(event.Data, "environment", "externalId")
	if externalID == "" {
		externalID = previousExternalID
	}

	input := h.getInput(event, cli)
	input.Files = templates
	if input.Environment == nil {
		input.Environment = map[string]string{}
	}

	for k, v := range environment {
		input.Environment[k] = v
	}

	if len(input.Files) == 0 {
		return nil, nil
	}

	return map[string]interface{}{
		"environment": map[string]interface{}{
			"upgrade": map[string]interface{}{
				"environment": environment,
				"templates":   templates,
				"externalId":  externalID,
			},
		},
	}, stack.Create(input)
}
