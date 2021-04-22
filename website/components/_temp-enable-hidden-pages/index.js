// Imports below are used in server-side only
import fs from 'fs'
import path from 'path'
import {
  generateStaticPaths as docsPageStaticPaths,
  generateStaticProps as docsPageStaticProps,
} from '@hashicorp/react-docs-page/server'

/**
 * DEBT
 * This is a short term hotfix for "hidden" docs-sidenav items.
 *
 * We likely do NOT want to support this in docs-page/server,
 * instead, a simple "hidden" attribute supported on docs-sidenav
 * nodes would do the trick, ensuring the "path" is "registered"
 * in the appropriate nav-data file, and located in the correct spot
 * in the nav-data tree, while also hiding that item in the sidebar.
 *
 * We can remove this hack with once support lands for "hidden" items,
 * currently this is somewhat blocked by branding rollout:
 * Asana task that will resolve this debt:
 * https://app.asana.com/0/1100423001970639/1200197752405255/f
 * Draft PR to support "hidden" nav items:
 * https://github.com/hashicorp/react-components/pull/220
 **/

const DEFAULT_PARAM_ID = 'page'

export async function generateStaticPaths({
  navDataFile,
  navDataFileHidden,
  localContentDir,
}) {
  const visiblePaths = await docsPageStaticPaths({
    navDataFile,
    localContentDir,
  })
  const hiddenPaths = await docsPageStaticPaths({
    navDataFile: navDataFileHidden,
    localContentDir,
  })
  return visiblePaths.concat(hiddenPaths)
}

export async function generateStaticProps({
  navDataFile,
  navDataFileHidden,
  localContentDir,
  product,
  params,
  paramId = DEFAULT_PARAM_ID,
  additionalComponents,
}) {
  // Read in the "hidden" nav data, and flatten it
  const navDataVisible = readNavData(navDataFile)
  const navDataHidden = readNavData(navDataFileHidden)
  // Check if this is a "hidden" page, if so, use the navDataHidden
  // to generate static props.
  const currentPath = params[paramId] ? params[paramId].join('/') : ''
  const hiddenPaths = flattenNavData(navDataHidden).map((n) => n.path)
  const isHiddenPage = hiddenPaths.filter((p) => p == currentPath).length > 0
  // Return the static props, but always pass the navDataVisible
  // as the navData to be displayed.
  const staticProps = await docsPageStaticProps({
    navDataFile: isHiddenPage ? navDataFileHidden : navDataFile,
    localContentDir,
    product,
    params,
    paramId,
    additionalComponents,
  })
  return { ...staticProps, navData: navDataVisible }
}

function readNavData(navDataFile) {
  const filePath = path.join(process.cwd(), navDataFile)
  return JSON.parse(fs.readFileSync(filePath))
}

function flattenNavData(nodes) {
  return nodes.reduce((acc, n) => {
    if (!n.routes) return acc.concat(n)
    return acc.concat(flattenNavData(n.routes))
  }, [])
}
