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

This [Dataminer/Proxy component](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cli/cmd/symbl-dataminer) aims to ingest the data, create any additional metadata to associate with these insights, and then save this context to recall later. There happens to be an excellent platform that does all of the heavy lifting for us. We can extract these conversation insights without having to train models, have expertise in Artificial Intelligence or Machine Learning, or require a team of data scientists. Of course, I'm talking about using the real-time streaming capabilities on the Symbl.ai Platform.

Preserving insights represents the first of two significant pieces of work in this design. In order to aggregate conversation insights from external conversation sources and through historical data, we need to have a method for persisting this data to recall and make associations to conversations happening now. This component bites off this aspect of the design.

#### Example Middleware Component

The last but very significant feature required in this Enterprise Reference Implementation for Conversation Analysis, which is making associations or defining the relationships between contextual insights, happens in this [Middleware component](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cli/cmd/example-your-middleware). This Middleware component is deeply tied to what your business cares about. This component, either in code or interfacing with another external system, captures your company's specific business rules. These business rules can then be used to notify others within the company to take action, create events that you might want to pass along to other software systems, or trigger actions you want to perform directly in this component.

Although there is a generic implementation provided in this Middleware component, the intent of this Enterprise Reference Implementation is only to be just that… a reference. This Middleware component contained in the repo should either, at minimum, be modified to capture your business rules or in practice, be reimplemented to fit your specific business needs.

#### Simple Client

Instead of building out a web client using a [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) platform, an easier to use interface into this system should just be taking your local laptop's microphone input and speaking into it. That what this [Simulated Web Client App](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cli/cmd/simulate-web-client-app) does.

## How to Deploy this Enterprise Reference Implementation

TODO: Coming soon...

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
