# TuteScrew

TuteScrew is a Discord chat bot written in Go for teamRPI's Discord. The goal of TuteScrew is to be fast, efficient, and easily expandable in the future.

Featuring...
- RPI Student Verification via CAS

With support to come for...
- Pickup game (PUG) scheduling and stat tracking
- Additional moderation features

## Getting Started
### Installation

1. Install Go
2. Run `go get github.com/albshin/tutescrew`
3. Fill out the configuration file *see below*

### Configuration

The configuration file must be filled out before the bot can run.

- Token: Can be obtained from [here](https://discordapp.com/developers/applications/me).
- Prefix: The prefix needed to call upon the bot.
- CASAuthURL: URL to RPI's CAS authentication server. `currently: https://cas-auth.rpi.edu/cas/login`
- CASRedirectURL: URL to redirect to upon authentication success.