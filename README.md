# CTO

## DockerHub

[cto-core](https://hub.docker.com/repository/docker/barklan/cto-core) | [cto-porter](https://hub.docker.com/repository/docker/barklan/cto-porter) | [cto-explorer](https://hub.docker.com/repository/docker/barklan/cto-explorer) | [fluentd-cto](https://hub.docker.com/repository/docker/barklan/fluentd-cto) | [cto-loginput](https://hub.docker.com/repository/docker/barklan/cto-loginput)

## Routing

#### Rest

[Docs](https://docs.ctopanel.com/)

```
/api
    /loginput (:8900 internal)
    /core (:8888 internal)
        /debug/{projectID}?key | GET
        /setproject{projectID} | POST
    /porter (:9010 internal)
```

#### gRPC

`porter` on 50051

## Sanity


```
       │   │        │
  front│   │bot     │fluentd
       ▼   ▼        ▼
     ┌───────┐ ┌──────────┐   r
┌───►│porter │ │ loginput ├────┐
│    └┬──┬──┬┘ └┬─────────┘    ▼
│     │  │  │   │            ┌───┐
│ ┌───┘  └──┼───┼───────────►│pg │
│ │         │   │        rw  └─┬─┘
│ │         ▼   ▼              │
│ │       ┌───────┐            │
│ │       │  mq   │           r│
│ │       └─┬───┬─┘            │
│ │   query │   │loginput      │
│ │   fanout│   │fanout        │
│ │         │   │              │
│ │choose ┌─┴───┴─┐ replicated │
│ │one    │ core  │ stateful  ─┘
│ └──────►│       │ cores
│         └┬──────┘
│ callback │
└──────────┘
```
