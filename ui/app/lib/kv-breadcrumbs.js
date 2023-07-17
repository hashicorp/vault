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
  /* 
    takes a path nested secret path "beep/boop" and returns an array of objects used for breadcrumbs: 
    [
    { label: 'beep', route: 'list-directory', model: 'beep'},
    { label: 'boop', route: 'list-directory', model: 'beep/boop'},
    { label: 'bop' ] } // last item is current route, only return label so breadcrumb has no link 
    ]
  */
  const pathAsArray = path.split('/').filter((path) => path);
  const modelIdArray = pathAsArray.map((key, index) => {
    return `${pathAsArray.slice(0, index + 1).join('/')}/`;
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
