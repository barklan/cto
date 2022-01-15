# CTO

## DockerHub

[cto-core](https://hub.docker.com/repository/docker/barklan/cto-core) | [cto-porter](https://hub.docker.com/repository/docker/barklan/cto-porter) | [cto-explorer](https://hub.docker.com/repository/docker/barklan/cto-explorer) | [fluentd-cto](https://hub.docker.com/repository/docker/barklan/fluentd-cto) | [cto-loginput](https://hub.docker.com/repository/docker/barklan/cto-loginput)

## Sanity check


```
         │   │        │
    front│   │bot     │fluentd
         │   │        ▼
       ┌─┴───┴─┐ ┌──────────┐
┌─────►│       │ │          │
│      │porter │ │ loginput ├────┐
│ ┌────┤       │ │          │    │
│ │    └┬──┬──┬┘ └┬─────────┘   r│
│ │     │  │  │   │            ┌─┴─┐
│ │     │  └──┼───┼───────────►│pg │
│ │   rw│     │   │        rw  └─┬─┘
│ │     ▼     ▼   ▼              │
│ │  w┌───┐ ┌───────┐           r│
├─┼──►│ c │ │  mq   │            │
│ │   └───┘ └─┬───┬─┘            │
│ │     query │   │loginput      │
│ │     fanout│   │fanout        │
│ │           │   │              │
│ │ choose  ┌─┴───┴─┐ replicated │
│ │ one     │ core  │ stateful  ─┘
│ └────────►│       │ cores
│           └┬──────┘
│   callback │
└◄───────────┘
```

- core - meant to be replicated (one replica per node)
- loginput - can be replicated
- porter - not sure...
- mq - can be replicated through [quorum queues](https://www.rabbitmq.com/quorum-queues.html)
- pg - can be replicated cockroachdb
- c - sure, why not
