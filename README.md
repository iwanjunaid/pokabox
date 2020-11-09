# Pokabox

Pokabox is [Go](https://golang.org/) transactional outbox pattern implementation for Postgres and Kafka.

For another implementation for MongoDB and Kafka, please see [Mokabox](https://github.com/iwanjunaid/mokabox).

## Table of Contents

- [How It Works?](#how-it-works)
- [Getting Started](#getting-started)
- [Events Handling](#events-handling)

## How It Works?

TODO

## Getting Started

1. Create ```outbox``` table

```sql
CREATE TABLE outbox(
  id            UUID PRIMARY KEY,
  group_id      UUID NOT NULL,
  kafka_topic   VARCHAR(255) NOT NULL,
  kafka_key     VARCHAR(255) DEFAULT NULL,
  kafka_value   TEXT NOT NULL,
  priority      INT DEFAULT 1000,
  status        VARCHAR(4) DEFAULT 'NEW',
  version       INTEGER DEFAULT 1,
  created_at    TIMESTAMP NOT NULL,
  sent_at       TIMESTAMP DEFAULT NULL
);
```

## Events Handling

TODO
