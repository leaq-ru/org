package config

type Config struct {
	ServiceName string
	Grpc        grpc
	MongoDB     mongodb
	STAN        stan
	NATS        nats
	LogLevel    string `envconfig:"LOGLEVEL"`
}

type grpc struct {
	Port string `envconfig:"GRPC_PORT"`
}

type mongodb struct {
	URL string `envconfig:"MONGODB_URL"`
}

type stan struct {
	ClusterID string `envconfig:"STAN_CLUSTERID"`
}

type nats struct {
	URL string `envconfig:"NATS_URL"`
}
