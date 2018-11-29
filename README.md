# SimpleGoChat

A minimal client/server chat application written in Go.

### Usage

To build the application, type `make` or `make go`. You must have Go installed and in your path. The program will be compiled to a static binary at `./app`.

To run, use the following command:

```
./app -host_ip <address> -host_port <port> -node_type ['client'|'server']
```

The `host_ip` field must be set to `localhost` for the server. When joining as a client, you will be prompted for a name. Clients are given the option of either creating a new (password-protected) room or joining an existing one.