# Reference Implementation for a Symbl.ai Enterprise Application

The goal of this repository is to provide both:

- an architecture for a highly scalable application with Enterprise class features used for Conversation Understanding
- a concrete example to be used as a Reference Implementation for the Architecture described above

## Characteristics of an Enterprise Conversation Application

What are some attributes of an Enterprise Application for Conversation Understanding? Here are some typical qualities that would probably fit this type of application:

- Scalable - all aspects of the application can be scaled out to accommodate more operational capacity (ie microservices)
- Highly Available - the application can absorb outages but still continue to function and be accessible to users
- Data Resiliency - data strategic to the operation of the application is accessible despite disruptions
- Data Recovery and Migration - Has the ability to withstand catastrophic failure and migrate data to other providers

In addition to the above characteristics, Enterprise-class features specifically a Conversation Analytics platform might also include but are not limited to:

- Real-Time Response - Able to react to new information as it's happening
- Historical Data - Conversation data is captured for purposes like dashboards, etc
- Data Aggregation - Has the ability to detect trends or patterns based on historical data over time
- Predictive Actions - Performing autonomous application driven actions based on data patterns discovered

## Architecture Diagram Overview

This is a high-level block diagram, that satisfies the above characteristics, for what the architecture of an Enterprise Conversation could look like...

![Enterprise Reference Architecture](https://github.com/dvonthenen/enterprise-reference-implementation/blob/main/docs/images/enterprise-architecture.png?raw=true)

Feature Benefits for this Enterprise Reference Architecture:

- Build applications with a historical conversation context
- Persist conversation insights (data ownership)
- Build scalable conversation applications
- Company’s business rules/logic pushed into backend server microservices
- Dashboards, dashboards, dashboards
- Historical data is air-gapped and can be backed up
- UI isolation. Change all aspects of UI frameworks without changing the code
- Also supports Asynchronous analysis of data

### Deep Dive of this Architecture

This architecture calls for a minimum of 3 required components. They are:

- Symbl.ai Dataminer/Proxy Service
- Middleware Plugin
- Client Interface (Web, CLI, App, etc)

#### **Symbl.ai Dataminer/Proxy**

This [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) aims to ingest the data, create any additional metadata to associate with these insights, and then save this context to recall later. There happens to be an excellent platform that does all of the heavy lifting for us, **cough cough** the [Symbl Platform](https://platform.symbl.ai). Using the Symbl Platform, we can extract these conversation insights without having to train models, have expertise in Artificial Intelligence or Machine Learning, or require a team of data scientists.

Preserving insights represents the first of two significant pieces of work in this design. In order to aggregate conversation insights from external conversation sources and through historical data, we need to have a method for persisting this data to recall and make associations to conversations happening now. This component bites off this aspect of the design.

#### **Middleware Plugin**

Middleware Plugins are a significant component required in this Enterprise Reference Implementation for Conversation Analysis. It provides associations or defining the relationships between contextual insights and your business.

A Middleware Plugin is deeply tied to what your business cares about. This plugin, either in code or interfacing with another external system, captures your company's specific business rules, goals, or what you hope to achieve. For example, these business rules can then be used to notify others within the company to take action, create events that you might want to pass along to other software systems, or trigger actions you want to perform directly in this component.

In this architecture, these events are called **Application Specific Messages**. These application specific messages are unique to your business and what you care about. These messages are created by your Middleware Plugin and sent to your UI or Client Interface. We will talk more about this in the next section.

To reiterate, since these plugins are codified business rules for **YOUR** business, these plugins will most likely require to be built by **YOU**. Having said that, we decided to create a pluggable framework for your Middleware needs which means that there is a great deal of code reuse that can happen if you fully leverage this framework.

#### **Client Application aka User Interface**

The last piece in this architecture is the Client-side Application. This interface can be an Angular Web UI, another REST service, [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) application, etc.

This interfaces does not and probably will not be built from scratch. It's more likely you already have an interface that your customers or clients already use. You will augment that User Interface to receive these **Application Specific Messages** to display within your interface.

For example, your company might provide a sales enablement platform that uses conversation insights to help increase and drives sales for your customers. Within your customer conversation on your platform, you might have received a Symbl event that a customer mentioned the topic surrounding a product you have within your profolio. An **Application Specific Messages** can be sent to your sales person within your UI to discuss that particular product to your customer.

## Enterprise Reference Implementation

Given the architecture above, the code in this repository implements the following cusing the described software platform or components:

![Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation/blob/main/docs/images/enterprise-architecture-implementation.png?raw=true)

The 3 components within this repo which correspond to the 3 components we talked about in the previous Overview section are:

- Symbl.ai Dataminer/Proxy Service
- Example Middleware Plugin
- Example Client Interface

And they *currently* leverage these external software platforms:

- [Neo4J](https://neo4j.com/) for Conversation Insight Persistence
  - [What is Neo4J?](https://neo4j.com/docs/getting-started/current/)
  - [Feature Basics](https://graphacademy.neo4j.com/courses/neo4j-fundamentals/)
- [RabbitMQ](https://rabbitmq.com/) for Message Bus/Exchange
  - [What is RabbitMQ?](https://blog.iron.io/what-is-rabbitmq/)
  - [Platform Support](https://www.rabbitmq.com/devtools.html)

### Symbl.ai Dataminer/Proxy

This [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) component can and should be re-used as-is. This does a lot of the heavy-lifting and will save you some time to implement a custm component and database schema that essentially can be reused as an off-the-shelf component. It is highly encouraged to use this component for reasons we will get into later on.

### Example Middleware Plugin

There is a generic scaffold implementation provided in this repo which sends an **Application Specific Messages** for each type Symbl insight (Topics, Trackers, Entities, etc). The intent of this Enterprise Reference Implementation is only to be just that... a reference. This [Example Middleware Plugin](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-middleware-plugin) contained in the repo is only a starting point for your own implementation of a plugin. This scaffold code should be modified to capture your business rules to fit your specific business needs.

### Example Simple Client

Instead of building out a full blown web client using a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, this repo provides a simple client that takes your local laptop's microphone input to provide the conversation. That's what this [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) does.

## How to Deploy the Example Application in this Repository

This section goes over how to deploy and run the components within this repo.

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

### Running the Example Middleware Plugin and Simulated Client

Make sure the Prerequisite Component instances are running before proceeding! Clone the [Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation) repo and the change your working directory to the root of that repo.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-reference-implementation.git
foo@bar:~$ cd enterprise-reference-implementation
```

Once you are at the root of the repo, you can start the [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) in the console by running the following commands:

```bash
foo@bar:~$ cd ./cmd/symbl-proxy-dataminer
foo@bar:~$ go run cmd.go
```

In a new console window when you are at the root of the repo, you can start the [Example Middleware Plugin](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-middleware-plugin) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-middleware-plugin
foo@bar:~$ go run cmd.go
```

Create a third/new console window when you are at the root of the repo, you can start the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-simulated-client-app
foo@bar:~$ go run cmd.go
```

Once the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) is running, you should be able to speak into your microphone to simulate a conversation that would take place on various platforms like a [Zoom](https://zoom.us/) meeting, on a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, etc.

As you start to speak into your microphone, you should see example application specific messages come through on the example client. Pretty simple!

## Running Plugins from Conversation Plugins Repo

Are you looking for a more meaningful demo or example of the real power of this architecture? We previously hinted at this implementation of a pluggable framework where you can start various Middleware Plugins to provide off-the-shelf capabilities. We have an [Enterprise Conversation Plugins](https://github.com/dvonthenen/enterprise-conversation-plugins) repo that serves as an App Store of pre-built functionality in the form individual plugins.

The first couple of plugins implemented in this repo are:
- the [Historical Insights Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/historical-insights) which triggers an Application Specific Message of the last 5 mentions of a Topic, Tracker, Entity, etc
- the [Statistical Insights Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/statistical-insights) which provides the number of times a topic, Tracker, Entity, etc has been mentioned in the past 30 mins, hour, 4 hours, day, 2 days, week and month.

### How Do I Launch These Plugins

To try these plugins out, make sure the Prerequisite Component instances are running before proceeding! Assuming you have already cloned the [Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation), clone the [Enterprise Conversation Plugins](hhttps://github.com/dvonthenen/enterprise-conversation-plugins) repo to your local laptop.

```bash
foo@bar:~$ git clone git@github.com:dvonthenen/enterprise-conversation-plugins.git
foo@bar:~$ cd enterprise-reference-implementation
```

Like before, start the [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) in the console by running the following commands:

In your first console windows, run:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-reference-implementation
foo@bar:~$ cd ./cmd/symbl-proxy-dataminer
foo@bar:~$ go run cmd.go
```

Instead of starting the Example Middleware Plugin, start the [Historical Insights Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/historical-insights) and the [Statistical Insights Plugin](https://github.com/dvonthenen/enterprise-conversation-plugins/tree/main/plugins/statistical-insights)  each in their own console window.

In a second console windows, run the `historical-insights` plugin by executing:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-plugins
foo@bar:~$ cd ./plugins/historical-insights
foo@bar:~$ go run cmd.go
```

In a third console windows, run the `statistical-insights` plugin by executing:
```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-conversation-plugins
foo@bar:~$ cd ./plugins/statistical-insights
foo@bar:~$ go run cmd.go
```

Finally, create a fourth console window to start the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) by running the following commands:

```bash
foo@bar:~$ cd ${REPLACE WITH YOUR ROOT DIR}/enterprise-reference-implementation
foo@bar:~$ cd ./cmd/example-simulated-client-app
foo@bar:~$ go run cmd.go
```

Then start talking into the microphone. Since both of these plugins are about aggregating conversations over time, close the `example-simulated-client-app` instance after having mentioned some topics, entities, trackers, etc. Then start up another instance of the `example-simulated-client-app` (which can be done in the same console window) and mention the same topics, entities, trackers, etc as in the previous conversation session. You should start to see some historical and statistical data flowing through to the client when those past insights are triggered.

## More Information

If you are looking for a detailed description and even a video that walks through this architecture diagram, please look at this blog called [Everything to Know About Enterprise Reference Implementation for Conversation Aggregation](https://symbl.ai/blog/everything-to-know-about-enterprise-reference-implementation-for-conversation-aggregation/).

Looking for more information about deployment strategies as they relate to storage, take a look at this blog post called [Databases and Persistent Storage for Conversation Data](https://symbl.ai/blog/databases-and-persistent-storage-for-conversation-data/).

### Contact Information

You can reach out to the Community via:

- [Google Group][google_group] for this Community Meeting and Office Hours
- Find us by using the [Community Calendar][google_calendar]
- Taking a look at the [Community Meeting and Office Hours Agenda Notes][agenda_doc]. Feel free to add any agenda items!
- Don't want to wait? Contact us through our [Community Slack][slack]
- If you want to do it the old fashion way, our email is community\[at\]symbl\[dot\]ai

[google_group]: https://bit.ly/3Cp5c9D
[google_calendar]: https://bit.ly/3jRGEj4
[agenda_doc]: https://bit.ly/3WH4hcO
[slack]: https://join.slack.com/t/symbldotai/shared_invite/zt-4sic2s11-D3x496pll8UHSJ89cm78CA
