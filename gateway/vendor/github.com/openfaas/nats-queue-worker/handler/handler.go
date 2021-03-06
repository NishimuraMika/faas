package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/go-nats-streaming"
	"github.com/openfaas/faas/gateway/queue"
)

// NatsQueue queue for work
type NatsQueue struct {
	nc stan.Conn
}

type NatsConfig interface {
	GetClientID() string
}

type DefaultNatsConfig struct {
}

func (DefaultNatsConfig) GetClientID() string {
	val, _ := os.Hostname()
	return "faas-publisher-" + val
}

// CreateNatsQueue ready for asynchronous processing
func CreateNatsQueue(address string, port int, clientConfig NatsConfig) (*NatsQueue, error) {
	queue1 := NatsQueue{}
	var err error
	natsURL := fmt.Sprintf("nats://%s:%d", address, port)
	log.Printf("Opening connection to %s\n", natsURL)

	clientID := clientConfig.GetClientID()
	clusterID := "faas-cluster"

	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	queue1.nc = nc

	return &queue1, err
}

// Queue request for processing
func (q *NatsQueue) Queue(req *queue.Request) error {
	var err error

	fmt.Printf("NatsQueue - submitting request: %s.\n", req.Function)

	out, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	err = q.nc.Publish("faas-request", out)

	return err
}
