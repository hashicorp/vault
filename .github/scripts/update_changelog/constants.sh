# Directories
HOME=$(pwd)
OUTPUT_DIR="${HOME}/output"
TEMPLATE_DIR="${HOME}/templates"
DOWNLOAD_DIR="${HOME}/downloads"

# Files
#CHANGELOG_URL="https://raw.githubusercontent.com/hashicorp/vault/main/CHANGELOG.md"
CHANGELOG_URL="https://github.com/hashicorp/vault/tree/${REPO}/${BRANCH}/CHANGELOG.md"
LOCAL_CHANGELOG="${OUTPUT_DIR}/changelog.md"
MD_TOC="${TEMPLATE_DIR}/toc.mdx"
MD_TAB_BODY="${TEMPLATE_DIR}/tab-body.mdx"
MD_TAB_OPEN="${TEMPLATE_DIR}/tab-open.mdx"
MD_TAB_CLOSE="${TEMPLATE_DIR}/tab-close.mdx"

# Constants
MAJOR_VERSION_DELTA="0.01"
