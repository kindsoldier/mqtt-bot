/*
 * Copyright:  Pixel Networks <support@pixel-networks.com> 
 * Author: Oleg Borodin <oleg.borodin@pixel-networks.com>
 */


package transport

import (
	"fmt"
	"log"
	//"net/url"
	//"os"
	"time"
    "strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
    clientId    string          = "pix-bot"

    keepalive   time.Duration   = 1 // sec
    waitTimeout time.Duration   = 3 // sec
    pingTimeout time.Duration   = 1 // sec
)


const (
    QosL1       byte            = 1
    QosL2       byte            = 2
    QosL3       byte            = 4

    //clientId    string          = "pm-xxxx"
)

type Transport struct {
    client  mqtt.Client
}

func New() *Transport {
    return &Transport{}
}

func (this *Transport) Bind(hostname string, port int, username string, password string) error {
    var err error

	opts := mqtt.NewClientOptions()

    uri := fmt.Sprintf("tcp://%s:%d", hostname, port)
	opts.AddBroker(uri)
	opts.SetUsername(username)
	opts.SetPassword(password)
    clientId := strconv.FormatInt(time.Now().Unix(), 10)
	opts.SetClientID(clientId)
    opts.SetAutoReconnect(true)

    opts.SetKeepAlive(keepalive)
    opts.SetPingTimeout(pingTimeout)

    onConnectHandler := func(client mqtt.Client) {
        log.Println("connect to broker")
    }
    opts.SetOnConnectHandler(onConnectHandler)

    onReconnectHandler := func(client mqtt.Client, opts *mqtt.ClientOptions) {
        log.Println("reconnect to broker")
    }
    opts.SetReconnectingHandler(onReconnectHandler)

	this.client = mqtt.NewClient(opts)

	token := this.client.Connect()
	for !token.WaitTimeout(waitTimeout * time.Second) {}

    err = token.Error()
	if err != nil {
		return err
	}
    return err
}

func (this *Transport) Publish(topic string, message string) {
        if this.client.IsConnected() {
            this.client.Publish(topic, 0, false, message)
        }
}

type Handler = func(subject string, message []byte)

func (this *Transport) Subscribe(topic string, callback Handler) error {
    var err error

    mqttHandler := func(mqttClient mqtt.Client, mqttMessage mqtt.Message) {
        callback(mqttMessage.Topic(), mqttMessage.Payload()) 
    }

    token := this.client.Subscribe(topic, QosL1, mqttHandler)
    for !token.WaitTimeout(waitTimeout * time.Second) {}

    err = token.Error()
    if err != nil {
        return err
    }
    return err
}
