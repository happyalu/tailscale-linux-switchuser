# Doesn't work anymore

This project does not work with latest changes in how tailscale maintains login data.

# Switch logged-in user of tailscale on Linux

I'm using two accounts with tailscale, wrote this to easily switch the network that should be used.

Usage:

```
go get github.com/happyalu/tailscale-linux-switchuser

sudo tailscale-linux-switchuser -login user@domain.com
```
