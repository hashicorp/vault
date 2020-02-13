# MongoDB Tests
The test `TestInit_clientTLS` cannot be run within CircleCI in its current form. This is because [it's not
possible to use volume mounting with the docker executor](https://support.circleci.com/hc/en-us/articles/360007324514-How-can-I-mount-volumes-to-docker-containers-).

Because of this, the test is skipped. Running this locally shouldn't present any issues as long as you have
docker set up to allow volume mounting from this directory:

```sh
go test -v -run Init_clientTLS
```

This may be able to be fixed if we mess with the entrypoint or the command arguments.
