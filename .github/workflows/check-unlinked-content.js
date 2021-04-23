var fs = require("fs");
var path = require("path");

const COLOR_RESET = "\x1b[0m";
const COLOR_GREEN = "\x1b[32m";
const COLOR_RED = "\x1b[31m";

runCheck([
  {
    contentDir: "website/content/docs",
    navDataFiles: [
      "website/data/docs-nav-data.json",
      "website/data/docs-nav-data-hidden.json",
    ],
  },
  {
    contentDir: "website/content/api-docs",
    navDataFiles: [
      "website/data/api-docs-nav-data.json",
      "website/data/api-docs-nav-data-hidden.json",
    ],
  },
  {
    contentDir: "website/content/guides",
    navDataFiles: ["website/data/guides-nav-data.json"],
  },
  {
    contentDir: "website/content/intro",
    navDataFiles: ["website/data/intro-nav-data.json"],
  },
]);

async function runCheck(baseRoutes) {
  const validatedBaseRoutes = await Promise.all(
    baseRoutes.map(async ({ contentDir, navDataFiles }) => {
      const missingRoutes = await validateMissingRoutes(
        contentDir,
        navDataFiles
      );
      return { contentDir, navDataFiles, missingRoutes };
    })
  );
  const allMissingRoutes = validatedBaseRoutes.reduce((acc, baseRoute) => {
    return acc.concat(baseRoute.missingRoutes);
  }, []);
  if (allMissingRoutes.length == 0) {
    console.log(
      `\n${COLOR_GREEN}✓ All content files have routes, and are included in navigation data.${COLOR_RESET}\n`
    );
  } else {
    validatedBaseRoutes.forEach(
      ({ contentDir, navDataFiles, missingRoutes }) => {
        if (missingRoutes.length == 0) return true;
        console.log(
          `\n${COLOR_RED}Error: Missing pages found in the ${contentDir} directory.\n\nPlease add these paths to ${navDataFiles.join(
            " or "
          )}, or remove the .mdx files.\n\n${JSON.stringify(
            missingRoutes,
            null,
            2
          )}${COLOR_RESET}\n\n`
        );
      }
    );
    process.exit(1);
  }
}

async function validateMissingRoutes(contentDir, navDataFiles) {
  // Read in nav-data.json, and make a flattened array of nodes
  const navDataFlat = navDataFiles.reduce((acc, navDataFile) => {
    const navDataPath = path.join(process.cwd(), navDataFile);
    const navData = JSON.parse(fs.readFileSync(navDataPath));
    return acc.concat(flattenNodes(navData));
  }, []);
  // Read all files in the content directory
  const files = await walkAsync(contentDir);
  // Filter out content files that are already
  // included in nav-data.json
  const missingPages = files
    // Ignore non-.mdx files
    .filter((filePath) => {
      return path.extname(filePath) == ".mdx";
    })
    // Transform the filePath into an expected route
    .map((filePath) => {
      // Get the relative filepath, that's what we'll see in the route
      const contentDirPath = path.join(process.cwd(), contentDir);
      const relativePath = path.relative(contentDirPath, filePath);
      // Remove extensions, these will not be in routes
      const pathNoExt = relativePath.replace(/\.mdx$/, "");
      // Resolve /index routes, these will not have /index in their path
      const routePath = pathNoExt.replace(/\/?index$/, "");
      return routePath;
    })
    // Determine if there is a match in nav-data.
    // If there is no match, then this is an unlinked content file.
    .filter((pathToMatch) => {
      // If it's the root path index page, we know
      // it'll be rendered (hard-coded into docs-page/server.js)
      const isIndexPage = pathToMatch === "";
      if (isIndexPage) return false;
      // Otherwise, needs a path match in nav-data
      const matches = navDataFlat.filter(({ path }) => path == pathToMatch);
      return matches.length == 0;
    });
  return missingPages;
}

function flattenNodes(nodes) {
  return nodes.reduce((acc, n) => {
    if (!n.routes) return acc.concat(n);
    return acc.concat(flattenNodes(n.routes));
  }, []);
}

function walkAsync(relativeDir) {
  const dirPath = path.join(process.cwd(), relativeDir);
  return new Promise((resolve, reject) => {
    walk(dirPath, function (err, result) {
      if (err) reject(err);
      resolve(result);
    });
  });
}

function walk(dir, done) {
  var results = [];
  fs.readdir(dir, function (err, list) {
    if (err) return done(err);
    var pending = list.length;
    if (!pending) return done(null, results);
    list.forEach(function (file) {
      file = path.resolve(dir, file);
      fs.stat(file, function (err, stat) {
        if (stat && stat.isDirectory()) {
          walk(file, function (err, res) {
            results = results.concat(res);
            if (!--pending) done(null, results);
          });
        } else {
          results.push(file);
          if (!--pending) done(null, results);
        }
      });
    });
  });
}
