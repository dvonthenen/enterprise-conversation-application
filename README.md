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

## Generic Architecture Diagram

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

### Enterprise Reference Implementation

Given the architecture above, the **code in this repository** implements the following using the described software platform or components:

![Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation/blob/main/docs/images/enterprise-architecture-implementation.png?raw=true)

The 3 major components in this reference implementation are:

- Symbl.ai Dataminer/Proxy
- Example Middleware Plugin
- Simple Client (Web, App, etc)

And *currently* leverages these external software platforms:

- [Neo4J](https://neo4j.com/) for Conversation Insight Persistence
  - [What is Neo4J?](https://neo4j.com/docs/getting-started/current/)
  - [Feature Basics](https://graphacademy.neo4j.com/courses/neo4j-fundamentals/)
- [RabbitMQ](https://rabbitmq.com/) for Message Bus/Exchange
  - [What is RabbitMQ?](https://blog.iron.io/what-is-rabbitmq/)
  - [Platform Support](https://www.rabbitmq.com/devtools.html)

#### Symbl.ai Dataminer/Proxy

This [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) aims to ingest the data, create any additional metadata to associate with these insights, and then save this context to recall later. There happens to be an excellent platform that does all of the heavy lifting for us. We can extract these conversation insights without having to train models, have expertise in Artificial Intelligence or Machine Learning, or require a team of data scientists. Of course, I'm talking about using the real-time streaming capabilities on the Symbl.ai Platform.

Preserving insights represents the first of two significant pieces of work in this design. In order to aggregate conversation insights from external conversation sources and through historical data, we need to have a method for persisting this data to recall and make associations to conversations happening now. This component bites off this aspect of the design.

#### Example Middleware Component

The last but very significant feature required in this Enterprise Reference Implementation for Conversation Analysis, which is making associations or defining the relationships between contextual insights, happens in this [Example Middleware/Analyzer component](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-middleware-analyzer). This Middleware component is deeply tied to what your business cares about. This component, either in code or interfacing with another external system, captures your company's specific business rules. These business rules can then be used to notify others within the company to take action, create events that you might want to pass along to other software systems, or trigger actions you want to perform directly in this component.

Although there is a generic implementation provided in this Middleware component, the intent of this Enterprise Reference Implementation is only to be just that… a reference. This Middleware component contained in the repo should either, at minimum, be modified to capture your business rules or in practice, be reimplemented to fit your specific business needs.

#### Simple Client

Instead of building out a web client using a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, an easier to use interface into this system should just be taking your local laptop's microphone input and speaking into it. That what this [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) does.

## How to Deploy the Example Application in this Repository

The example application contained within this repo provides the last 5 mentions of Trackers, Entity, and Topics in **PREVIOUS** conversations. An implementation of this Reference Enterprise Architecture will consist of at least 3 components, they are:

- [Symbl Proxy/Dataminer](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer)
- [Example Middleware/Analyzer component](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-middleware-analyzer)
- [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app)

### Prerequisite Software

For the easiest way to run this locally, you are going to need to install:

- [Docker Desktop](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/doc/install)

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

To clean up both of these instance, you need to first `stop` then `rm` the instances in Docker.

```bash
foo@bar:~$ docker container ls --all
foo@bar:~$ docker stop <container id>
foo@bar:~$ docker rm <container id>
```

### Running the Example Application

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

In a new console window when you are at the root of the repo, you can start the [Example Middleware/Analyzer component](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-middleware-analyzer) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-middleware-analyzer
foo@bar:~$ go run cmd.go
```

Create a third/new console window when you are at the root of the repo, you can start the [Example Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) by running the following commands:

```bash
foo@bar:~$ cd ./cmd/example-simulated-client-app
foo@bar:~$ go run cmd.go
```

Once the [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) is running, you should be able to speak into your microphone to simulate a conversation that would take place on various platforms like a [Zoom](https://zoom.us/) meeting, on a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, etc.

To really get the effectiveness out of this Enterprise Reference Implementation example, you should create  at least 1 [custom Tracker](https://docs.symbl.ai/docs/custom-trackers) that you want to trigger in your first conversation (ie when you first start and speaking into the [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app)) and then closing [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) by hitting any key. This will simulate the completion of a single conversation.

Since the Example in this repo is exploring the use case of aggregating conversation over time, you want to start another conversation session by launching the [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app) again and then trigger a second instance of the same [custom Tracker](https://docs.symbl.ai/docs/custom-trackers). In the console of the [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app), you should see an Application-level message showing you the last reference(s) made to that [custom Tracker](https://docs.symbl.ai/docs/custom-trackers). You can do this repeatly and you will see the past 5 Trackers, Entities, and Topics mentions in previous conversations in the console of your [Simulated Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/example-simulated-client-app).

## More Information

If you are looking for a detailed description and even a video that walks through this architecture diagram, please look at this [blog post](https://symbl.ai/blog/everything-to-know-about-enterprise-reference-implementation-for-conversation-aggregation/).

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
