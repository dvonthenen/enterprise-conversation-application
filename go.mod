module github.com/dvonthenen/enterprise-reference-implementation

go 1.18

require (
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	k8s.io/klog/v2 v2.80.1
)

require (
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
)

replace github.com/koding/websocketproxy => github.com/dvonthenen/websocketproxy v0.0.0-20221130030525-2d35aa6682e7
// replace github.com/koding/websocketproxy => ../../koding/websocketproxy
