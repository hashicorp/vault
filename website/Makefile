configure_cache:
	mkdir -p tmp/cache

build: configure_cache
	@echo "==> Starting build in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd):/opt/buildhome/repo" \
		--volume "$(shell pwd)/tmp/cache:/opt/buildhome/cache" \
		--env "ENV=production" \
		netlify/build \
		build middleman build --verbose

website: configure_cache
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd):/opt/buildhome/repo" \
		--volume "$(shell pwd)/tmp/cache:/opt/buildhome/cache" \
		--publish "4567:4567" \
		--publish "35729:35729" \
		--env "ENV=production" \
		netlify/build \
		build middleman

.PHONY: build website
