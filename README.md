# draincat

Like netcat but for Heroku logs.

## Usage

Spin up a HTTP drain and dump the logs to a file.

```
bin/draincat -p 5000 > logs.txt
```

## Features

* JSON output


## TODO

- [X] Multiplex simutaneous lpx frames using Go channels
- [ ] Structured data field (from log-shuttle)
- [ ] Syslog drains (`draincat --type=syslog`)
- [ ] metrics?
- [X] JSON output  
