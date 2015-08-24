# draincat

Like netcat but for Heroku logs.

## Usage

Spin up a HTTP drain and dump the logs to a file.

```
bin/draincat -p 5000 > logs.txt
```

## Features

* JSON output
* Serial processing of logplex frames

## TODO

- [X] Multiplex simutaneous lpx frames using Go channels
- [X] Print logs in the same frame together (no interleaving with other frames)
- [ ] Structured data field (from log-shuttle)
- [ ] Syslog drains (`draincat --type=syslog`)
- [ ] metrics?
- [X] JSON output  
