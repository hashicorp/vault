/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

function pathIsDirectory(pathToSecret) {
  // This regex only checks for / at the end of the string. ex: boop/ === true, boop/bop === false;
  return pathToSecret ? !!pathToSecret.match(/\/$/) : false;
}

function pathIsFromDirectory(path) {
  // This regex just looks for a / anywhere in the path. ex: boop/ === true, boop/bop === true;
  return path ? !!path.match(/\//) : false;
}

function breadcrumbsForDirectory(path) {
  // path === "beep/boop/"
  const pathAsArray = path.split('/').filter((path) => path);
  const modelIdArray = pathAsArray.map((key, index) => {
    return `${pathAsArray.slice(0, index + 1).join('/')}/`; // ex: ['beep/', 'beep/boop/']. We need these model Ids to tell the LinkTo on the breadcrumb what to put into the dynamic *pathToSecret on the breadcrumb: ex/kv/beep/boop/
  });

  return pathAsArray.map((key, index) => {
    // we do not want to return "route or model" on the last item otherwise it will add link to the current page.
    if (pathAsArray.length - 1 === index) {
      return { label: key };
    }
    return { label: key, route: 'list-directory', model: modelIdArray[index] };
  });
}

export { breadcrumbsForDirectory, pathIsDirectory, pathIsFromDirectory };
