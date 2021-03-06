package diegonats

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pivotal-golang/lager"
)

type NATSClientRunner struct {
	addresses string
	username  string
	password  string
	logger    lager.Logger
	client    NATSClient
}

func NewClientRunner(addresses, username, password string, logger lager.Logger, client NATSClient) NATSClientRunner {
	return NATSClientRunner{
		addresses: addresses,
		username:  username,
		password:  password,
		logger:    logger.Session("nats-runner"),
		client:    client,
	}
}

func (runner NATSClientRunner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	natsMembers := []string{}
	for _, addr := range strings.Split(runner.addresses, ",") {
		uri := url.URL{
			Scheme: "nats",
			User:   url.UserPassword(runner.username, runner.password),
			Host:   addr,
		}
		natsMembers = append(natsMembers, uri.String())
	}

	unexpectedConnClosed, err := runner.client.Connect(natsMembers)
	for err != nil {
		runner.logger.Error("connecting-to-nats-failed", err)
		select {
		case <-signals:
			return nil
		case <-time.After(time.Second):
			unexpectedConnClosed, err = runner.client.Connect(natsMembers)
		}
	}

	runner.logger.Info("connecting-to-nats-succeeeded")
	close(ready)

	select {
	case <-signals:
		// BUG(tedsuo): until we have ordered group shutdown, closing the client
		// causes unexpected behavior. https://www.pivotaltracker.com/story/show/81411580

		//	runner.client.Close()
		runner.logger.Info("shutting-down")
		return nil
	case <-unexpectedConnClosed:
		runner.logger.Error("unexpected-nats-close", nil)
		return errors.New("nats closed unexpectedly")
	}
}
