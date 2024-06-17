package conf

type Kafka struct {
	GroupId        string        `yaml:"groupId"`
	ExampleTopic   string        `yaml:"exampleTopic"`
	KafkaConsumerA KafkaConsumer `yaml:"consumer-a"`
	KafkaProducer  KafkaProducer `yaml:"producer"`
}

type KafkaConsumer struct {
	Brokers []string `yaml:"brokers"`
}

type KafkaProducer struct {
	Brokers []string `yaml:"brokers"`
}
