# Repo which we are cloning and executing npm run build:deploy-preview within
REPO_TO_CLONE=dev-portal
# Set the subdirectory name for the base project
PREVIEW_DIR=website-preview
# The directory we want to clone the project into
CLONE_DIR=website-preview
# The product for which we are building the deploy preview
PRODUCT=vault
# Preview mode, controls the UI rendered (either the product site or developer). Can be `io` or `developer`
PREVIEW_MODE=io

from_cache=false

if [ -d "$PREVIEW_DIR" ]; then
  echo "$PREVIEW_DIR found"
  CLONE_DIR="$PREVIEW_DIR-tmp"
  from_cache=true
fi

# Clone the base project, if needed
echo "⏳ Cloning the $REPO_TO_CLONE repo, this might take a while..."
git clone --depth=1 "https://github.com/hashicorp/$REPO_TO_CLONE.git" "$CLONE_DIR"

if [ "$from_cache" = true ]; then
  echo "Setting up $PREVIEW_DIR"
  cp -R "./$CLONE_DIR/." "./$PREVIEW_DIR"
fi

# cd into the preview directory project
cd "$PREVIEW_DIR"

# Run the build:deploy-preview start script
PREVIEW_MODE=$PREVIEW_MODE REPO=$PRODUCT HASHI_ENV=project-preview npm run build:deploy-preview
