package main

import (
	"io/ioutil"
	"log"

	"github.com/cam-stitt/qago"
	"gopkg.in/yaml.v2"
)

func main() {
	data, err := ioutil.ReadFile("fixtures/helloworld.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	suite := qago.Suite{}
	err = yaml.Unmarshal(data, &suite)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println(suite)

	qago.RunSuite(&suite)
}
