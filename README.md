# draincat

Like netcat but for Heroku logs.

## Usage

Spin up a HTTP drain and dump the logs to a file.

```
bin/draincat > logs.txt
```

## Features

* JSON output `DRAINCAT_JSON=1`


## TODO

- [X] Multiplex simutaneous lpx frames using Go channels
- [ ] Option to process multiple logs in single frame together?
- [ ] Syslog drains (`draincat --type=syslog`)
- [ ] metrics?
- [X] JSON output  
