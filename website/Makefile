.DEFAULT_GOAL := website

PWD=$$(pwd)
DOCKER_IMAGE="hashicorp/dev-portal"
DOCKER_IMAGE_LOCAL="dev-portal-local"
DOCKER_RUN_FLAGS=-it \
		--publish "3000:3000" \
		--rm \
		--tty \
		--volume "$(PWD)/content:/app/content" \
		--volume "$(PWD)/public:/app/public" \
		--volume "$(PWD)/data:/app/data" \
		--volume "$(PWD)/redirects.js:/app/redirects.js" \
		--volume "next-dir:/app/website-preview/.next" \
		--volume "$(PWD)/.env:/app/.env" \
		-e "REPO=vault" \
		-e "PREVIEW_MODE=io"

# Default: run this if working on the website locally to run in watch mode.
.PHONY: website
website:
	@echo "==> Downloading latest Docker image..."
	@docker pull $(DOCKER_IMAGE)
	@echo "==> Starting website..."
	@docker run $(DOCKER_RUN_FLAGS) $(DOCKER_IMAGE)

# Use this if you have run `website/build-local` to use the locally built image.
.PHONY: website/local
website/local:
	@echo "==> Starting website from local image..."
	@docker run $(DOCKER_RUN_FLAGS) $(DOCKER_IMAGE_LOCAL)

# Run this to generate a new local Docker image.
.PHONY: website/build-local
website/build-local:
	@echo "==> Building local Docker image"
	@docker build https://github.com/hashicorp/dev-portal.git\#main \
		-t $(DOCKER_IMAGE_LOCAL)

