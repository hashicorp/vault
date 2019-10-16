configure-cache:
	@mkdir -p tmp/cache

build: configure-cache
	@echo "==> Starting build in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd):/opt/buildhome/repo" \
		--volume "$(shell pwd)/tmp/cache:/opt/buildhome/cache" \
		--env "ENV=production" \
		netlify/build \
		build "sh bootstrap.sh && middleman build --verbose"

website: configure-cache
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
		build "sh bootstrap.sh && middleman"

website-linux: configure-cache
	@echo "==> Starting website in Docker..."
	@docker run \
	    --name "netlify" \
	    -d \
		--rm \
		--tty \
		--publish "4567:4567" \
		--publish "35729:35729" \
		--env "ENV=production" \
		netlify/build
	@docker exec -it netlify mkdir -p /opt/buildhome/cache
	@docker exec -it netlify mkdir -p /opt/buildhome/repo
	@docker cp /home/tbex/go/src/github.com/hashicorp/vault/website/. netlify:/opt/buildhome/repo/
	@docker exec -it netlify build "sh bootstrap.sh && middleman build --verbose"

.PHONY: configure-cache build website
