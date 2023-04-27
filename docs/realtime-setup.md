# Realtime/Streaming Enterprise Conversation Application

The code in this repository implements an Enterprise Conversation Application for Realtime/Streaming use cases based on the architecture described on the README at the root of the repo.

There are 3 components within this repo that correspond to the components in the Overview section are:

- Symbl.ai Realtime Proxy/Dataminer Service
- Example Realtime Middleware Plugin
- Example Client Interface

If you place these component within the architecture diagram, you get the implementation diagram for your application below:

![Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation/blob/main/docs/images/enterprise-architecture-implementation.png?raw=true)

And that implementation leverages these external software platforms:

- [Neo4J](https://neo4j.com/) for Conversation Insight Persistence
  - [What is Neo4J?](https://neo4j.com/docs/getting-started/current/)
  - [Feature Basics](https://graphacademy.neo4j.com/courses/neo4j-fundamentals/)
- [RabbitMQ](https://rabbitmq.com/) for Message Bus/Exchange
  - [What is RabbitMQ?](https://blog.iron.io/what-is-rabbitmq/)
  - [Platform Support](https://www.rabbitmq.com/devtools.html)

## Symbl.ai Proxy/Dataminer Service

This [Symbl.ai Proxy/Dataminer Service](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) aims to ingest the data, create any additional metadata to associate with these insights, and then save this context to recall later. There happens to be an excellent platform that does all of the heavy lifting for us, **cough cough** the [Symbl Platform](https://platform.symbl.ai). Using the Symbl Platform, we can extract these conversation insights without having to train models, have expertise in Artificial Intelligence or Machine Learning, or require a team of data scientists.

This Proxy/Dataminer Service makes use of the [Realtime/Streaming API](https://docs.symbl.ai/docs/streaming-api) to receive conversation insights from the Symbl Platform. For each realtime insight emitted by the platform, this component will persist the conversation insight into the graph database and then send a corresponding RabbitMQ message for that insight. This message will then be consumed by your middleware plugin to perform some action of your chosing.

## Example Realtime Middleware Plugin

The [Example Realtime Middleware Plugin](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-plugin) contained in the repo is only a starting point for your own implementation of a plugin. This scaffold code should be modified to capture your business rules to fit your specific business needs.

This example makes use of the [Realtime Streaming API](https://docs.symbl.ai/docs/streaming-api) to receive conversation insights from the Symbl Platform. This generic scaffold implementation sends an example **Application Specific Message** for each type Symbl insight (Topics, Trackers, Entities, etc). These application specific messages are unique to your business and what you care about. These messages are created by your Middleware Plugin and sent to your UI or Client Interface.

## Example Simple Client

Instead of building out a full blown web client using a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, this repo provides a simple client that takes your local laptop's microphone input to provide the conversation. That's what this [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-simulated-client) does.

> **_NOTE:_** A web UI based example is coming shortly!

## How to Deploy the Example Realtime/Streaming Application in this Repository

This section goes over how to deploy and run the components required for realtime/streaming use cases within this repo.

### Prerequisite Software

If you don't have this software running on your laptop, you are going to need to install:

- [Docker Desktop](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/doc/install)

> **_NOTE:_** Container images will be made available soon. This should dramatically decrease the barrier of entry to get this stuff up an running.

### Prerequisite Components

The Example Application in this repo, makes use of a Graph database called [Neo4j](https://neo4j.com/) and also an AMQP server called [RabbitMQ](https://rabbitmq.com/).

To launch an instance of a [Neo4j](https://neo4j.com/) database, run the following Docker command:

```bash
foo@bar:~$ docker run \
              -d \
              --publish=7474:7474 --publish=7687:7687 \
              -v $HOME/neo4j/data:/data \
              -v $HOME/neo4j/logs:/logs \
              -v $HOME/neo4j/import:/var/lib/neo4j/import \
              -v $HOME/neo4j/plugins:/plugins \
              neo4j:5
```

To launch an instance of a [RabbitMQ](https://rabbitmq.com/) database, run the following Docker command:

```bash
foo@bar:~$ docker run \
                -d -p 5672:5672 -p 15672:15672 \
                --hostname my-rabbit --name my-rabbit rabbitmq:3
```

To clean up both of these instance when you are finished with them, you need to first `stop` then `rm` the instances in Docker.

```bash
foo@bar:~$ docker container ls --all
foo@bar:~$ docker stop <container id>
foo@bar:~$ docker rm <container id>
```

### Running the Example Enterprise Application for Realtime Conversations

Make sure the Prerequisite Component instances are running before proceeding! Clone the [Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation) repo and the change your working directory to the root of that repo.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-reference-implementation.git
foo@bar:~$ cd enterprise-reference-implementation
```

Once you are at the root of the repo, you can start the [Symbl Proxy/Dataminer Service](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) in the console by running the following commands:

```bash
foo@bar:~$ cd ./cmd/symbl-proxy-dataminer
foo@bar:~$ go run cmd.go
```

In a new console window when you are at the root of the repo, you can start the [Example Middleware Plugin](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-plugin) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-realtime-plugin
foo@bar:~$ go run cmd.go
```

Create a third/new console window when you are at the root of the repo, you can start the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-simulated-client) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-realtime-simulated-client
foo@bar:~$ go run cmd.go
```

Once the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-simulated-client) is running, you should be able to speak into your microphone to simulate a conversation that would take place on various platforms like a [Zoom](https://zoom.us/) meeting, on a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, etc.

As you start to speak into your microphone, you should see example application specific messages come through on the example client. Pretty simple!

## Running Plugins from Conversation Plugins Repo

Are you looking for a more meaningful demo or example of the real power of this architecture? We previously hinted at this implementation of a pluggable framework where you can start various Middleware Plugins to provide off-the-shelf capabilities. We have an [Enterprise Conversation Plugins](https://github.com/dvonthenen/enterprise-conversation-plugins) repo that serves as an App Store of pre-built functionality in the form individual plugins.

The first couple of plugins implemented in this repo are:
- the [Historical Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/realtime/historical) which triggers an Application Specific Message of the last 5 mentions of a Topic, Tracker, Entity, etc
- the [Statistical Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/realtime/statistical) which provides the number of times a topic, Tracker, Entity, etc has been mentioned in the past 30 mins, hour, 4 hours, day, 2 days, week and month.

### How Do I Launch These Plugins

To try these plugins out, make sure the Prerequisite Component instances are running before proceeding! Assuming you have already cloned the [Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation), clone the [Enterprise Conversation Plugins](hhttps://github.com/dvonthenen/enterprise-conversation-plugins) repo to your local laptop.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-conversation-plugins.git
foo@bar:~$ cd enterprise-reference-implementation
```

Like before, start the [Symbl Proxy/Dataminer Service](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) in the console by running the following commands:

In your first console windows, run:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-reference-implementation
foo@bar:~$ cd ./cmd/symbl-proxy-dataminer
foo@bar:~$ go run cmd.go
```

Instead of starting the Example Middleware Plugin, start the [Historical Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/realtime/historical) and the [Statistical Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/realtime/statistical) each in their own console window.

In a second console windows, run the `historical` plugin by executing:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-plugins
foo@bar:~$ cd ./plugins/realtime/historical
foo@bar:~$ go run cmd.go
```

In a third console windows, run the `statistical` plugin by executing:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-plugins
foo@bar:~$ cd ./plugins/realtime/statistical
foo@bar:~$ go run cmd.go
```

Finally, create a fourth console window to start the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-realtime-simulated-client) by running the following commands:

```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-reference-implementation
foo@bar:~$ cd ./cmd/example-realtime-simulated-client
foo@bar:~$ go run cmd.go
```

Then start talking into the microphone. Since both of these plugins are about aggregating conversations over time, close the `example-realtime-simulated-client` instance after having mentioned some topics, entities, trackers, etc. Then start up another instance of the `example-realtime-simulated-client` (which can be done in the same console window) and mention the same topics, entities, trackers, etc as in the previous conversation session. You should start to see some historical and statistical data flowing through to the client when those past insights are triggered.
