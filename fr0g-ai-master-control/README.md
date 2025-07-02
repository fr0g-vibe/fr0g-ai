# fr0g-ai-master-control
## purpose
this repository contains the logic controller for AI automation tasks that leverages the fr0g-ai-bridge, fr0g-ai-aip gRPC interfaces.

https://github.com/fr0g-vibe/fr0g-ai-aip
https://github.com/fr0g-vibe/fr0g-ai-bridge

## features
* golang based
* continually monitors external inputs
  * inputs
    * webhook
    * cron
    * incoming email
    * phone calls
    * text messages
    * socials
* leverages gRPC for AI inference and persona management
* execution rules are established in a plain language system such as OPA to manage what master control can and can't do
## use-cases
* when an email is recieved consult a counsel of ai personas (aip) to determine what should be done with the email
  * master-control determines if the information (summarized) should be forwarded to a human, archived, or dropped
  * master-control provides regular updates to the user with batched email communications periodically
* user provides a list of RSS feeds to master-control which regularly pulls rss feed data and generates summaries of current events
  * master-control submits summaries to communities of AiP to get community "responses" and "opinions"
  * master-control leverages AiP to attempt to predict short-term, medium-term, and long-term implications of collections of current event news by "connecting the dots"
* user connects a discord websocket to a new future component fr0g-ai-socket that accepts various websocket and webhook functionality
  * master-control responds to various types of incoming webhook and websocket actions with a highly modular structure for modifying response behavior
* more to come
