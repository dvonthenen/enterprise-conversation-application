# Asynchronous Enterprise Conversation Application

The code in this repository implements an Enterprise Conversation Application for Asynchronous use cases based on the architecture described on the README at the root of the repo.

There are 2 components within this repo that correspond to the components in the Overview section are:

- Symbl.ai Asynchronous REST/Dataminer Service
- Example Realtime Middleware Plugin

If you place these component within the architecture diagram, you get the implementation diagram for your application below:

![Enterprise Conversation Application](https://github.com/dvonthenen/enterprise-conversation-application/blob/main/docs/images/enterprise-asynchronous-architecture.png?raw=true)

And that implementation leverages these external software platforms:

- [Neo4J](https://neo4j.com/) for Conversation Insight Persistence
  - [What is Neo4J?](https://neo4j.com/docs/getting-started/current/)
  - [Feature Basics](https://graphacademy.neo4j.com/courses/neo4j-fundamentals/)
- [RabbitMQ](https://rabbitmq.com/) for Message Bus/Exchange
  - [What is RabbitMQ?](https://blog.iron.io/what-is-rabbitmq/)
  - [Platform Support](https://www.rabbitmq.com/devtools.html)

## Symbl.ai Asynchronous REST/Dataminer Service

This [Symbl Asynchronous REST/Dataminer Service](https://github.com/dvonthenen/enterprise-conversation-application/tree/main/cmd/symbl-rest-dataminer) component can and should be re-used as-is. This does a lot of the heavy-lifting and will save you some time to implement a custom component and database schema that essentially can be reused as an off-the-shelf component. It is highly encouraged to use this component for reasons we will get into later on.

This REST/Dataminer Service makes use of the [Asynchronous API](https://docs.symbl.ai/docs/async-api) to receive conversation insights from the Symbl Platform. This component exposes a REST endpoint which takes in a conversation ID and then retreives all of the conversation insights for that pre-processed conversation ID. For each insight retrieved from the platform, this component will send a corresponding RabbitMQ message for each type of insight. This message will then be consumed by your middleware plugin to perform some action of your chosing.

## Example Asynchronous Middleware Plugin

The [Example Asynchronous Middleware Plugin](https://github.com/dvonthenen/enterprise-conversation-application/tree/main/cmd/example-asynchronous-plugin) contained in the repo is only a starting point for your own implementation of a plugin. This scaffold code should be modified to capture your business rules to fit your specific business needs.

Since these conversations aren't tied to conversation currently happening (aka ones happening in realtime) but rather conversation that occurred asynchronously in the past, this generic scaffold implementation does **NOT** send Application Specific Messages for each type Symbl insight (Topics, Trackers, Entities, etc). There **could** be use cases for processing past conversations which may require an immediate notification of someone in a meeting, but right now it's hard to imagine that use case.

The purpose of this component is to receive conversation insights from the RabbitMQ messages sent by the REST/Dataminer Service and perform some desired action based on your application needs or business goals. For example, maybe you have a [Tracker](https://docs.symbl.ai/docs/trackers) that corresponds to wanting more Widgets from your Widget company. You could send an email to that user with current offers available at your Widget company.

## How to Deploy the Example Asynchronous Application in this Repository

This section goes over how to deploy and run the components required for Asynchronous/REST use cases within this repo.

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

### Running the Example Enterprise Application for Asynchronous Conversations

Make sure the Prerequisite Component instances are running before proceeding! Clone the [Enterprise Conversation Application](https://github.com/dvonthenen/enterprise-conversation-application) repo and the change your working directory to the root of that repo.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-conversation-application.git
foo@bar:~$ cd enterprise-conversation-application
```

Once you are at the root of the repo, you can start the [Symbl Asynchronous REST/Dataminer Service](https://github.com/dvonthenen/enterprise-conversation-application/tree/main/cmd/symbl-asynchronous-dataminer) in the console by running the following commands:

```bash
foo@bar:~$ cd ./cmd/symbl-asynchronous-dataminer
foo@bar:~$ go run cmd.go
```

In a new console window when you are at the root of the repo, you can start the [Example Middleware Plugin](https://github.com/dvonthenen/enterprise-conversation-application/tree/main/cmd/example-asynchronous-plugin) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-asynchronous-plugin
foo@bar:~$ go run cmd.go
```

To invoke this workflow, you are going to need a past conversation process. If you don't have one, you want log into the [Symbl Platform](https://platform.symbl.ai/) and go to the [API Explorer](https://platform.symbl.ai/#/explorer/topics) to process one of the example conversations provided.

To run that conversation or any previously processed conversation through this application, you will need that conversation ID to invoke that following command:

```bash
foo@bar:~$ curl -k https://127.0.0.1/conversation/[conversation-id]
```

You should see all of the insights get dumped into the console where Example Middleware Plugin is running. These are the insights to be consumed by your business specific application and where you would want to perform some action to do.

## Running Plugins from Conversation Plugins Repo

Are you looking for a more meaningful demo or example of the real power of this architecture? We previously hinted at this implementation of a pluggable framework where you can start various Middleware Plugins to provide off-the-shelf capabilities. We have an [Enterprise Conversation Plugins](https://github.com/dvonthenen/enterprise-conversation-plugins) repo that serves as an App Store of pre-built functionality in the form individual plugins.

The first couple of plugins implemented in this repo are:

- the [Email Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/asynchronous/email) sends an email when a configured Topic, Tracker or Entity is encountered
- the [Webhook Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/asynchronous/webhook) sends a JSON of the entire conversation to a specified URI when a configured Topic, Tracker or Entity is encountered

### How Do I Launch These Plugins

To try these plugins out, make sure the Prerequisite Component instances are running before proceeding! Assuming you have already cloned the [Enterprise Conversation Application](https://github.com/dvonthenen/enterprise-conversation-application), clone the [Enterprise Conversation Plugins](hhttps://github.com/dvonthenen/enterprise-conversation-plugins) repo to your local laptop.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-conversation-plugins.git
foo@bar:~$ cd enterprise-conversation-application
```

Like before, start the [Symbl Asynchronous REST/Dataminer Service](https://github.com/dvonthenen/enterprise-conversation-application/tree/main/cmd/symbl-asynchronous-dataminer) in the console by running the following commands:

In your first console window, run:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-application
foo@bar:~$ cd ./cmd/symbl-asynchronous-dataminer
foo@bar:~$ go run cmd.go
```

Instead of starting the Example Middleware Plugin, start the [Email Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/asynchronous/email). For this plugin,...

**-----------------------------**
**TODO: ENV VAR**
**-----------------------------**

In a second console window, run the `email` plugin by executing:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-plugins
foo@bar:~$ cd ./plugins/asynchronous/email
foo@bar:~$ go run cmd.go
```

**---------------------------------------------**
**TODO: No Async Client... use CURL**
**---------------------------------------------**

Finally, create a third console window. To being processing a previously recorded conversation, you can call the REST API in the [TODO Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/asynchronous/email) by running the following commands:

```bash
foo@bar:~$ curl -k https://127.0.0.1/conversation/[conversation-id]
```

**---------------------------------------------**
**TODO: Describe Process**
**---------------------------------------------**

You should start to see some historical and statistical data flowing through to the client when those past insights are triggered.
