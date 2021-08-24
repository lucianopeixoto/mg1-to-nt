package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// Here I'll add the part that handles what the server will do when a message is
	// received on a subscribed topic.
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	fmt.Printf("Received message from topic: %s\n", msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("MQTT Client Connected to Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	fmt.Println("Minew G1 adapter for nTopus Indoor Location.")
	fmt.Println("v 0.00")
	fmt.Println("Starting App...")

	var mqttBrokerHost string
	flag.StringVar(&mqttBrokerHost, "mqtt-host", "localhost", "IP or hostname of the MQTT Broker")
	var mqttBrokerPort = 1883
	flag.IntVar(&mqttBrokerPort, "mqtt-port", 1883, "MQTT Broker port")
	var mqttSubscribeTopic string
	flag.StringVar(&mqttSubscribeTopic, "mqtt-topic", "#", "MQTT topic to subscribe to")
	flag.Parse()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerHost, mqttBrokerPort))
	opts.SetClientID("minew-g1-ntopus-adapter")
	//opts.SetUsername("ntopus")
	//opts.SetPassword("ntopus")
	opts.SetDefaultPublishHandler(messagePubHandler)
	fmt.Printf("Connecting to %s\n", mqttBrokerHost)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client, mqttSubscribeTopic)

	fmt.Printf("Press Ctrl+C to end\n")
	WaitForCtrlC()
	fmt.Printf("\n")

	client.Disconnect(250)
}

// To use latter just in case the app needs to send some command to the G1s
//
/*func publish(client mqtt.Client) {
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
	}
}*/

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func WaitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}
