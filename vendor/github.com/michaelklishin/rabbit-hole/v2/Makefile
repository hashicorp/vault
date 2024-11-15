.DEFAULT_GOAL := test

.PHONY: all
all: test

.PHONY: test
test:
	go test -v
.PHONY: tests
tests: test

COVER_FILE := coverage
.PHONY: cover
cover:
	go test -v -test.coverprofile="$(COVER_FILE).prof"
	sed -i.bak 's|_'$(GOPATH)'|.|g' $(COVER_FILE).prof
	go tool cover -html=$(COVER_FILE).prof -o $(COVER_FILE).html
	rm $(COVER_FILE).prof*

.PHONY: ginkgo
ginkgo:
	command -v ginkgo || go get -u github.com/onsi/ginkgo/ginkgo
	ginkgo -v

.PHONY: docker
docker:
	docker run --rm \
	  --interactive --tty --entrypoint /bin/bash \
	  --volume $(CURDIR):/usr/src/app --workdir /usr/src/app \
	  golang:1.12

.PHONY: docker.rabbitmq
docker.rabbitmq:
	docker run -p 15672:15672 -p 5672:5672 -p 4639:4639 --name rabbithole_rabbitmq -dit rabbitmq:3.12-management
	sleep 10
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmqctl set_cluster_name rabbitmq@localhost"
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmqctl add_vhost /"
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmqctl set_permissions -p / guest \".*\" \".*\" \".*\""
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmqctl add_vhost rabbit/hole"
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmqctl set_permissions -p rabbit/hole guest \".*\" \".*\" \".*\""
	docker exec -ti rabbithole_rabbitmq /bin/bash -c "rabbitmq-plugins enable rabbitmq_federation rabbitmq_federation_management rabbitmq_shovel rabbitmq_shovel_management"
