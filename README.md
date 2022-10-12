
## Register the plugin factory
```go
import (
    "github.com/hosie/firefly-nats/pkg/nats"
    "github.com/hyperledger/firefly/pkg/events/eifactory"
) 
eifactory.RegisterFactory(&nats.Factory{})
```
NOTE: you must do the above before calling `firefly.Run()`


## Start the NATS server
Add the following to `docker-compose.override.yml` for your stack

```yaml
   nats-server:
       image: nats:latest
       ports:
           - 8222:8222
           - 6222:6222
           - 4222:4222
```

## Enable the plugin
Add the following to `firefly_core.yml`
```
event:
 dbevents:
   bufferSize: 10000
 transports:
   enabled: ["nats","websockets","webhooks"]
```
Restart the stack

## Start the NATS subscriber
From root of the repo
```
go run test/sub.go
```



## Create a firefly subscription
```bash
curl -X 'POST' \
 'http://127.0.0.1:5000/api/v1/subscriptions' \
 -H 'accept: application/json' \
 -H 'Request-Timeout: 2m0s' \
 -H 'Content-Type: application/json' \
 -d '
{
   "name": "nats1",
   "namespace": "default",
   "transport": "nats"
 }
'
```