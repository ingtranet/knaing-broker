package broker

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

func newStanClient(config *viper.Viper) (stan.Conn, error) {
	clusterID := config.GetString("stan_cluster_id")
	var clientID string
	hostname, err := os.Hostname()
	if err != nil {
		clientID = hostname
	} else {
		clientID = uuid.New().String()
	}
	natsURL := config.GetString("nats_url")

	logger.Info().Msg(fmt.Sprintf("creating connection with %s %s %s", clusterID, clientID, natsURL))
	client, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.Pings(3, 20))
	if err != nil {
		return nil, errors.Wrap(err, "creating stan client failed: " + natsURL)
	}
	return client, nil
}


