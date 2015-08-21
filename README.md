# draincat

Like netcat but for Heroku logs.

## Usage

Spin up a HTTP drain and dump the logs to a file.

```
bin/draincat > logs.txt
```

## Features

* JSON output (`DRAINCAT_JSON=1`)
* Serial processing of logplex frames

## TODO

- [X] Multiplex simutaneous lpx frames using Go channels
- [X] Print logs in the same frame together (no interleaving with other frames)
- [ ] Syslog drains (`draincat --type=syslog`)
- [ ] metrics?
- [X] JSON output  
