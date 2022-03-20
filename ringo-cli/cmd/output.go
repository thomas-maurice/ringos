package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func output(obj interface{}) {
	switch marshaller {
	case "json":
		b, err := json.MarshalIndent(&obj, "", "  ")
		if err != nil {
			logrus.Fatalf("could not marshall struct: %s", err)
		}
		fmt.Println(string(b))
	case "yaml":
		b, err := yaml.Marshal(&obj)
		if err != nil {
			logrus.Fatalf("could not marshall struct: %s", err)
		}
		fmt.Println(string(b))
	default:
		fmt.Println(obj)
	}
}
