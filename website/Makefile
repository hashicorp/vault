build:
	@echo "==> Starting build in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd):/website" \
		-e "ENV=production" \
		netlify/build \
		build cd website && middleman build --verbose

website:
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--volume "$(shell pwd):/website" \
		--publish "4567:4567" \
		--publish "35729:35729" \
		-e "ENV=production" \
		netlify/build \
		build cd website && middleman

.PHONY: build website
