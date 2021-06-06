# Default: run this if working on the website locally to run in watch mode.
website:
	@echo "==> Downloading latest Docker image..."
	@docker pull hashicorp/vault-website
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--workdir "/website" \
		--volume "$(shell pwd):/website" \
		--volume "/website/node_modules" \
		--publish "3000:3000" \
		hashicorp/vault-website \
		npm start

# This command will generate a static version of the website to the "out" folder.
build:
	@echo "==> Downloading latest Docker image..."
	@docker pull hashicorp/vault-website
	@echo "==> Starting build in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--workdir "/website" \
		--volume "$(shell pwd):/website" \
		--volume "/website/node_modules" \
		hashicorp/vault-website \
		npm run static

# If you are changing node dependencies locally, run this to generate a new
# local Docker image with the dependency changes included.
build-image:
	@echo "==> Building Docker image..."
	@docker build -t hashicorp-vault-website-local .

# Use this if you have run `build-image` to use the locally built image
# rather than our CI-generated image to test dependency changes.
website-local:
	@echo "==> Downloading latest Docker image..."
	@docker pull hashicorp/vault-website
	@echo "==> Starting website in Docker..."
	@docker run \
		--interactive \
		--rm \
		--tty \
		--workdir "/website" \
		--volume "$(shell pwd):/website" \
		--volume "/website/node_modules" \
		--publish "3000:3000" \
		hashicorp-vault-website-local \
		npm start

.DEFAULT_GOAL := website
.PHONY: build build-image website website-local
