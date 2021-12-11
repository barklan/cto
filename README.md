# cto

Preconfigured fluentd in docker swarm cluster:

```yaml
version: "3.7"
services:

  fluentd-cto:
    image: 'barklan/fluentd-cto:1.0.0'
    environment:
      FLUENTD_HOSTNAME: 'example-env.com'
      FLUENTD_HTTP_DUMP_ENDPOINT: '...'
    ports:
      - 24224:24224
      - 24224:24224/udp
    user: root
    networks:
      - traefik-public
    deploy:
      mode: global
      restart_policy:
        condition: on-failure

networks:
  traefik-public:
    external: ${TRAEFIK_PUBLIC_NETWORK_IS_EXTERNAL-true}
```

Minimal project config (logger only):

```yaml
project_id: "example"

envs:
  - 'example-env.com'

tg:
  chat_id: -##############
```
