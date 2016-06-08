# RabbitMQ Backend

## Testing

There are unit and integration RabbitMQ backend tests. Unit tests can be run by `go test`. Integration tests require setting the following environment variables:
```
RABBITMQ_CONNECTION_URI=
RABBITMQ_USERNAME=
RABBITMQ_PASSWORD=
```