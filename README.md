# go-sky

go-sky is a Golang micro service framework integrating several very popular packages and tools. As follows:

* github.com/go-kit/kit
* github.com/uber-go/zap
* gopkg.in/natefinch/lumberjack.v2
* github.com/opentracing/opentracing-go
* github.com/go-redis/redis/v8
* github.com/elastic/go-elasticsearch/v7
* github.com/Shopify/sarama
* github.com/prometheus

go-sky includes this subpackages as follows:

* github.com/WiFeng/go-sky/sky/config
* github.com/WiFeng/go-sky/sky/database
* github.com/WiFeng/go-sky/sky/elasticsearch
* github.com/WiFeng/go-sky/sky/helper
* github.com/WiFeng/go-sky/sky/http
* github.com/WiFeng/go-sky/sky/kafka
* github.com/WiFeng/go-sky/sky/log
* github.com/WiFeng/go-sky/sky/metrics
* github.com/WiFeng/go-sky/sky/redis
* github.com/WiFeng/go-sky/sky/trace

## Features

1. Support config.toml, and load separated config file by different runtime enviroment (config_development.toml/config_production.toml)
2. Support tracing (include http server and http client / redis / sql / kafka / elasticsearch)
3. Support log rotating and include trace_id in all log items.
4. Support promethues metric (include http server by now)

## Usage

1. There is a demo project in sky-example.

## TODO

1. Support redis/sql operation metric
2. Support more custom config to make some function disabled.

## Feedback

We need more cute members to make this project more robust. Believe you can!
