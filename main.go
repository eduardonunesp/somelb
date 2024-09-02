package main

import (
	"bytes"
)

var yamlConfig = `
servers:
  - backend:
      hosts:
        - url: "http://localhost:8080"
        - url: "http://localhost:8081"
    frontend:
      port: 1337
`

func main() {
	buf := bytes.NewBufferString(yamlConfig)
	config, err := ReadConfig(buf)
	if err != nil {
		panic(err)
	}

	s, err := NewServer(config.Servers[0])
	if err != nil {
		panic(err)
	}

	err = s.Run()
	if err != nil {
		panic(err)
	}
}
