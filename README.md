# CheckThisOut

Config:

```yaml
...
```

## Logging

### API

- POST /api/log/input - fluentd input
- GET /api/log/exact?key=...
- GET /api/log/range?query=...

### Internal

Store in badger in separate log database.

- Key: `enviroment service_name time flag randInt`
- Value: take raw `[]byte` request directly

### Thanks

https://github.com/antfu/vitesse
https://github.com/antfu/vitesse-lite
https://github.com/leezng/vue-json-pretty
