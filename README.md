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

Includes this subpackages as follows:

* github.com/WiFeng/go-sky/config
* github.com/WiFeng/go-sky/database
* github.com/WiFeng/go-sky/elasticsearch
* github.com/WiFeng/go-sky/helper
* github.com/WiFeng/go-sky/http
* github.com/WiFeng/go-sky/kafka
* github.com/WiFeng/go-sky/log
* github.com/WiFeng/go-sky/metrics
* github.com/WiFeng/go-sky/redis
* github.com/WiFeng/go-sky/trace

## Related projects

1. github.com/WiFeng/go-sky-example
2. github.com/WiFeng/go-sky-helloworld

## Features

1. Support config.toml, and load separated config file by different runtime enviroment (config_development.toml/config_production.toml)
2. Support many popular componets including sql/redis/kafka/elasticsearch.
3. Support tracing (include http server and http client / redis / sql / kafka / elasticsearch)
4. Support log rotating and include trace_id in all log items.
5. Support promethues metric (include http server by now)

![image](https://user-images.githubusercontent.com/2247568/107139748-82f40200-6958-11eb-856e-467afb1868c4.png)

![image](https://user-images.githubusercontent.com/2247568/140611536-04d28cd7-a0b0-4d0d-b76e-cddb2e521d6a.png)


## Usage

1. There is a demo project in sky-example.

## TODO

1. Support redis/sql operation metric
2. Support more custom config to make some function disabled.

## Contribution

We need more cute members to make this project more robust. Believe you can!
