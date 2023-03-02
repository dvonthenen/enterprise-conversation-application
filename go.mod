module github.com/dvonthenen/enterprise-reference-implementation

go 1.18

require (
	github.com/dvonthenen/rabbitmq-manager v0.1.1
	github.com/dvonthenen/symbl-go-sdk v0.1.5-0.20230302195324-986cb6f026d8
	github.com/dvonthenen/websocketproxy v0.1.0-dyv.1
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/neo4j/neo4j-go-driver/v5 v5.3.0
	github.com/r3labs/sse/v2 v2.9.0
	k8s.io/klog/v2 v2.90.0
)

require (
	github.com/dvonthenen/websocket v1.5.1-dyv.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gordonklaus/portaudio v0.0.0-20220320131553-cc649ad523c1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace github.com/r3labs/sse/v2 => github.com/dvonthenen/sse/v2 v2.0.0-20221222171132-1daa5f8b774c

// replace github.com/dvonthenen/symbl-go-sdk => ../../dvonthenen/symbl-go-sdk
// replace github.com/dvonthenen/websocket => ../../dvonthenen/websocket
// replace github.com/dvonthenen/websocketproxy => ../../dvonthenen/websocketproxy
// replace github.com/dvonthenen/sse => ../../r3labs/sse
