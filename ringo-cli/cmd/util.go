package cmd

import (
	"github.com/sirupsen/logrus"
	ringo "github.com/thomas-maurice/ringos/go-client"
)

func MustGetClient() *ringo.Client {
	client, err := ringo.NewClient(&ringo.ClientConfig{
		TargetAddress: address,
		Username:      username,
		Password:      password,
		Debug:         debug,
		AuthEnabled:   password != "",
	})

	if err != nil {
		logrus.WithError(err).Fatalf("could not create a client")
	}

	return client
}
