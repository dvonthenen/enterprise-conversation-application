module github.com/dvonthenen/enterprise-reference-implementation

go 1.18

require (
	github.com/dvonthenen/symbl-go-sdk v0.1.3-0.20221201180251-284ac126bab0
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	github.com/neo4j/neo4j-go-driver/v5 v5.3.0
	k8s.io/klog/v2 v2.80.1
)

require (
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gordonklaus/portaudio v0.0.0-20220320131553-cc649ad523c1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace github.com/koding/websocketproxy => github.com/dvonthenen/websocketproxy v0.0.0-20221130030525-2d35aa6682e7

// replace github.com/koding/websocketproxy => ../../koding/websocketproxy
