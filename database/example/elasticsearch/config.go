package main

type (
	elasticsearch struct {
		Hosts    []string `yaml:"hosts"`
		Index    []string `yaml:"index"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
		Size     int      `yaml:"size"`
	}
)

var (
	Elasticsearch = &elasticsearch{}
)