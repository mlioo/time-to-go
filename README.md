# Time to go! (Mass Slack channel leaver)
Sick of endless slack channels? Here is a a quick and dirty way to go through and mass leave those pesky slack channels that keep building it. Think of this like when you close chrome when you have too many tabs open to start fresh. 

***warning be careful when leaving private channels you can't rejoin unless someone invites you back!***

# How to setup
Create a new slack bot with the following user scopes and install it to your workspace.

```
channels:read
channels:write
groups:read
groups:write
```

Then set the environment variable `SLACK_TOKEN` with the User OAuth Token and follow through the prompts.
