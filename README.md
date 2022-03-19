# ware-kv 

[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php) 
[![experimental](http://badges.github.io/stability-badges/dist/experimental.svg)](http://github.com/badges/stability-badges)

It's a kv in-memory database, based on **HTTP RESTful** type.

## Why *ware-kv*?

- Take [Redis](https://github.com/redis/redis), which is the main stream open-source kv in-memory database, for standard. But *ware-kv* got something different:
  - Just unified HTTP RESTful interfaces, no SDK is required, out of the box!
  - Redis's basic element can only store `string` , but ware-kv can store
    - string
    - integer number
    - float number
    - list (except `set`)
    - map-dict (except `set`)
  - Contains some lightweight middlewares, like:
    - Message Queue
    - Bloom Filter
    - Distributed Lock
  - Support web-based platform tools, including:
    - Operation usage
    - Performance monitoring
  - Provide double-write consistency plan for MySQL.
- Thread-safe for sure!
- Maybe there's no database like *ware-kv*?
- By the way, complete the graduation project! :)

## How *ware-kv*?

### Basic *Wares*

- string
- list
- sort-list
- object
- set
- bitmap

### Special *Wares*

- counter
- bloom filter
- distributed lock
- ~~message queue~~ (not yet...)
- ~~cache~~ (not yet...)

### Others

- Support consistency for *crash-safe*.
  - Tracker (Logic Log like `WAL`)
  - ~~Camera (Physics Log)~~(not yet...)

- Support *pub/sub* keys.
- Support set key's *expire time*.
- ~~Support *cache eviction*.~~ (not yet...)
- ~~Support *double-write consistency* plan for MySQL.~~ (not yet...)

## Ideas Came From? 🧠

- [Elastic Search](https://github.com/elastic/elasticsearch) Style boost me...

- Various middleware scattered in every corner...

So, I just want to build a **modern** **lightweight** No-SQL(maybe kv) database, which is integrated with **common middleware and common problem solutions**.

## Incremental Plan? 🎯

If time permits,  I'll add:(Now *ware-kv* is just stand-alone environment)

- Distributed Cluster

- Sentinel

- Data Sharding

- Master-Slave Replication

