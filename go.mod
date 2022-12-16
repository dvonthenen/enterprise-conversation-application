module github.com/dvonthenen/enterprise-reference-implementation

go 1.18

require (
	github.com/dvonthenen/symbl-go-sdk v0.1.4-0.20221213174244-8c4877c3ed9f
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	github.com/neo4j/neo4j-go-driver/v5 v5.3.0
	github.com/rabbitmq/amqp091-go v1.5.0
	k8s.io/klog/v2 v2.80.1
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
)

replace (
	github.com/gorilla/websocket => github.com/dvonthenen/websocket v1.5.1-0.20221123154619-09865dbf1be2
	github.com/koding/websocketproxy => github.com/dvonthenen/websocketproxy v0.0.0-20221207172044-14b1fc90f46a
)

// replace github.com/gorilla/websocket => ../../dvonthenen/websocket
// replace github.com/koding/websocketproxy => ../../koding/websocketproxy
