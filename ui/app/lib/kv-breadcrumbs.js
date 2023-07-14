/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

function pathIsFolder(secretPrefix) {
  // This regex only checks for / at the end of the string. ex: boop/ === true, boop/bop === false;
  return secretPrefix ? !!secretPrefix.match(/\/$/) : false;
}

function pathIsFromNested(path) {
  // This regex just looks for a / anywhere in the path. ex: boop/ === true, boop/bop === true;
  return path ? !!path.match(/\//) : false;
}

function breadcrumbsForNestedSecret(path) {
  // path === "meep/moop/"
  const pathAsArray = path.split('/').filter((path) => path);
  const modelIdArray = pathAsArray.map((key, index) => {
    return `${pathAsArray.slice(0, index + 1).join('/')}/`; // ex: ['meep/', 'meep/moop/']. We need these model Ids to tell the LinkTo on the breadcrumb what to put into the dynamic *secretPrefix on the breadcrumb: ex/kv/meep/moop/
  });

  return pathAsArray.map((key, index) => {
    // we do not want to return "route or model" on the last item otherwise it will add link to the current page.
    if (pathAsArray.length - 1 === index) {
      return { label: key };
    }
    return { label: key, route: 'list-nested-secret', model: modelIdArray[index] };
  });
}

export { breadcrumbsForNestedSecret, pathIsFolder, pathIsFromNested };
