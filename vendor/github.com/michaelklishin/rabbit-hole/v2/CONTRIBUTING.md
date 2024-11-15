## Contributing

The workflow is pretty standard:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Run integration tests (see below)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push -u origin my-new-feature`)
6. Submit a pull request

## Running Tests

### Required Plugins

The test suite assumes you have a RabbitMQ node running on localhost with `rabbitmq_management` and
`rabbitmq_shovel_management` plugins enabled and that
`rabbitmqctl` is available in `PATH` (or `RABBITHOLE_RABBITMQCTL` points to it).

To enable the plugins:

``` shell
./bin/ci/before_build.sh
```

That will enable dependencies and reduce node's stats emission interval.

### Setting Up Virtual Hosts and Permissions

Before running the tests, make sure to run `bin/ci/before_build.sh` that will create a vhost and user(s) needed
by the test suite.

### Running Tests

The project uses [Ginkgo](http://onsi.github.io/ginkgo/) and [Gomega](https://github.com/onsi/gomega).

To clone dependencies and run tests, use `make`. It is also possible
to use the brilliant [Ginkgo CLI runner](http://onsi.github.io/ginkgo/#the-ginkgo-cli) e.g.
to only run a subset of tests.
