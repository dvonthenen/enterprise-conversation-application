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

![Enterprise Reference Architecture](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/images/enterprise-architecture.png)

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

Given the architecture above, the code in this repository implements the following using the described software platform or components:

![Enterprise Reference Implementation](https://github.com/dvonthenen/enterprise-reference-implementation/tree/main/docs/images/enterprise-architecture-implementation.png)

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
