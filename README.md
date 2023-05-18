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

This architecture consists of several components which differ based on the type of conversation application you are building. They are:

- Symbl.ai Proxy/Dataminer Service
- Symbl.ai Asynchronous Service
- Middleware Plugin
- Client Interface (Web, CLI, App, etc)

#### **Symbl.ai Proxy/Dataminer Service**

This [Symbl.ai Proxy/Dataminer Service](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-proxy-dataminer) aims to ingest the data, create any additional metadata to associate with these insights, and then save this context to recall later. There happens to be an excellent platform that does all of the heavy lifting for us, **cough cough** the [Symbl Platform](https://platform.symbl.ai). Using the Symbl Platform, we can extract these conversation insights without having to train models, have expertise in Artificial Intelligence or Machine Learning, or require a team of data scientists.

Preserving insights represents the first of two significant pieces of work in this design. In order to aggregate conversation insights from external conversation sources and through historical data, we need to have a method for persisting this data to recall and make associations to conversations happening now. This component bites off this aspect of the design.

This Proxy/Dataminer Service makes use of the [Realtime/Streaming API](https://docs.symbl.ai/docs/streaming-api) to receive conversation insights from the Symbl Platform. For each realtime insight emitted by the platform, this component will persist the conversation insight into the graph database and then send a corresponding RabbitMQ message for that insight. This message will then be consumed by your middleware plugin to perform some action of your chosing.

#### **Symbl.ai Asynchronous Service**

This [Symbl Asynchronous REST/Dataminer Service](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/cmd/symbl-rest-dataminer) component can and should be re-used as-is. This does a lot of the heavy-lifting and will save you some time to implement a custom component and database schema that essentially can be reused as an off-the-shelf component. It is highly encouraged to use this component for reasons we will get into later on.

This REST/Dataminer Service makes use of the [Asynchronous API](https://docs.symbl.ai/docs/async-api) to receive conversation insights from the Symbl Platform. This component exposes a REST endpoint which takes in a conversation ID and then retreives all of the conversation insights for that pre-processed conversation ID. For each insight retrieved from the platform, this component will also persist conversation insights into a graph database and then send a corresponding RabbitMQ message for each type of insight. This message will then be consumed by your middleware plugin to perform some action of your chosing.

#### **Middleware Plugin**

Middleware Plugins are a significant component required in this Enterprise Reference Implementation for Conversation Analysis. It provides associations or defining the relationships between contextual insights and your business.

A Middleware Plugin is deeply tied to what your business cares about. This plugin, either in code or interfacing with another external system, captures your company's specific business rules, goals, or what you hope to achieve. For example, these business rules can then be used to notify others within the company to take action, create events that you might want to pass along to other software systems, or trigger actions you want to perform directly in this component.

There are two flavors of Middleware Plugins:
- Realtime for processing conversation currently happening (like in a meeting)
- Asynchronous for processing past/recorded conversations

To reiterate, since these plugins are codified business rules or actions for **YOUR** business, these plugins will most likely require to be built by **YOU**. Having said that, we decided to create a pluggable framework for your Middleware needs which means that there is a great deal of code reuse that can happen if you fully leverage this framework.

#### **Client Application aka User Interface**

If you are using this project to process Realtime Conversations, the last piece in this architecture is the Client-side Application. This is effectively your User Interface (UI). This interface can be an Angular Web UI, Golang Command Line Interface (CLI), another REST service, [CPaaS](https://www.gartner.com/en/information-technology/glossary/communications-platform-service-cpaas) application, etc.

This interfaces does not and probably will not be built from scratch. It's more likely you already have an interface that your customers or clients already use. You will augment that User Interface to receive these **Application Specific Messages** to display within your interface.

For example, your company might provide a sales enablement platform that uses conversation insights to help increase and drives sales for your customers. Within your customer conversation on your platform, you might have received a Symbl event that a customer mentioned the topic surrounding a product you have within your profolio. An **Application Specific Messages** can be sent to your sales person within your UI to discuss that particular product to your customer.

If you are processing **ONLY** Asynchronous conversations, you will not have an UI component to notify users and thus not need to send **Application Specific Messages**. Your UI will most likely consist of some kind of dashboard for your users to consume conversation insights.

## How to Deploy the Example Application in this Repository

There are 3 main configurations for the implementation contained in this repo.

**Realtime Conversation Processing**
To deploy this configuration, follow this setup guide: [https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/realtime-setup.md](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/realtime-setup.md).

**Asynchronous Conversation Processing**
To deploy this configuration, follow this setup guide: [https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/asynchronous-setup.md](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/asynchronous-setup.md).

**Realtime and Asynchronous Conversation Process**
To deploy this configuration, follow each setup guide above. The setup is independent of each other and the only share components between these two configurations are the [Neo4J](https://neo4j.com/) Database and [RabbitMQ](https://rabbitmq.com/) Server.

## More Information

If you are looking for a detailed description and even a video that walks through this architecture diagram, please look at this blog called [Everything to Know About Enterprise Reference Implementation for Conversation Aggregation](https://symbl.ai/blog/everything-to-know-about-enterprise-reference-implementation-for-conversation-aggregation/).

Looking for more information about deployment strategies as they relate to storage, take a look at this blog post called [Databases and Persistent Storage for Conversation Data](https://symbl.ai/blog/databases-and-persistent-storage-for-conversation-data/).

### Contact Information

You can reach out to the Community via:

- [Google Group][google_group] for this Community Meeting and Office Hours
- Find us by using the [Community Calendar][google_calendar]
- Taking a look at the [Community Meeting][community_meeting]. Feel free to add any topics to [agenda doc][agenda_doc]!
- Bring all questions to the [Office Hours][office_hours].
- Don't want to wait? Contact us through our [Community Slack][slack]
- If you want to do it the old fashion way, our email is symblai-community-meeting\[at\]symbl\[dot\]ai

[google_group]: https://bit.ly/3Cp5c9D
[google_calendar]: https://bit.ly/3jRGEj4
[agenda_doc]: https://bit.ly/3WH4hcO
[community_meeting]: bit.ly/3M13vDg
[office_hours]: bit.ly/3LTbELg
[slack]: https://join.slack.com/t/symbldotai/shared_invite/zt-4sic2s11-D3x496pll8UHSJ89cm78CA
