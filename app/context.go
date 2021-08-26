package app


import (
	"github.com/danakum/go-util/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	BaseUrl string `yaml:"base_url" json:"base_url"`
}

var Config config

func Init() {
	file, err := ioutil.ReadFile(`config/surge.yaml`)
	if err != nil {
		log.Fatal(`Cannot open config file`, `, config/mqtt.yaml, `, err)
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(`Cannot parse config file`, `config/mqtt.yaml`, err)
	}
	log.Info(`Surge config established:`, Config)
}
