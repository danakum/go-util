package mqtt

import (
	"github.com/danakum/go-util/config"
	"github.com/danakum/go-util/log"
	PahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
	"time"
)

type conf struct {
	Brokers              []string `yaml:"brokers" json:"brokers"`
	ClientID             string   `yaml:"client_id" json:"client_id"`
	User                 string   `yaml:"user" json:"user"`
	Password             string   `yaml:"password" json:"password"`
	PingTimeout          int      `yaml:"ping_timeout" json:"ping_timeout"`
	MaxReconnectInterval int      `yaml:"max_reconnect_interval" json:"max_reconnect_interval"`
	ConnectTimeout       int      `yaml:"connect_timeout" json:"connect_timeout"`
	MessageChannelDepth  int      `yaml:"message_channel_depth" json:"message_channel_depth"`
}

//var mqttConf conf

func Init(clientId string, filePath string, onConnect PahoMqtt.OnConnectHandler) PahoMqtt.Client {
	conf := parseConfig(filePath)
	opts := PahoMqtt.NewClientOptions()
	opts.Username = conf.User
	opts.Password = conf.Password
	opts.CleanSession = false
	opts.AutoReconnect = true
	opts.PingTimeout = time.Duration(conf.PingTimeout) * time.Second
	opts.MaxReconnectInterval = time.Duration(conf.MaxReconnectInterval) * time.Second
	opts.ConnectTimeout = time.Duration(conf.ConnectTimeout) * time.Second
	opts.OnConnectionLost = func(client PahoMqtt.Client, e error) {
		log.Error(`mqtt client disconnected `, e)
	}
	opts.MessageChannelDepth = uint(conf.MessageChannelDepth)
	opts.OnConnect = onConnect
	for _, addr := range conf.Brokers {
		opts.AddBroker(`tcp://` + addr)
	}

	opts.ClientID = conf.ClientID
	if clientId != `` {
		opts.ClientID = clientId
	}

	client := PahoMqtt.NewClient(opts)

	go func(c PahoMqtt.Client) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		select {
		case sig := <-signals:
			c.Disconnect(200)
			log.Info(`Mqtt connection aborted : `, sig)
			break
		}
	}(client)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(`Cannot connect to the broker : `, token.Error())
		return client
	}

	log.Info(`MQTT Connection establish for client `, opts.ClientID)

	return client
}

func parseConfig(filePath string) (conf) {
	mqttConf := conf{}
	config.DefaultConfigurator.Load(filePath, &mqttConf, func(config interface{}) {

	})

	return mqttConf
}
