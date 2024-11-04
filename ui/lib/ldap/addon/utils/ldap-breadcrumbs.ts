/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const ldapBreadcrumbs = (
  path: string, // i.e. path/to/item
  roleType: string,
  mountPath: string,
  isCurrent = false
) => {
  const ancestry = path.split('/').filter(Boolean); // remove last empty string

  return ancestry.map((name: string, idx: number) => {
    const isLast = ancestry.length === idx + 1;
    // if the end of the path is the current route, don't return a route link
    if (isLast && isCurrent) return { label: name };
    // each segment is a continued concatenation of ancestral paths.
    // for example, if the full path to an item is "prod/admin/west"
    // the segments will be: prod/, prod/admin/, prod/admin/west.
    // LIST or GET requests can then be made for each crumb accordingly.
    let segment = ancestry.slice(0, idx + 1).join('/');
    // although the trailing slash in the URL isn't ideal, it allows us to keep the path name consistent
    // otherwise we have to manually keep track of whether we're in a directory or at the end of the path.
    segment = isLast ? segment : `${segment}/`;
    const models = [mountPath, roleType, segment];
    return {
      label: name,
      // if the last crumb is a subdirectory path, it's also the current route so we return above
      route: isLast ? 'roles.role.details' : 'roles.subdirectory',
      models,
    };
  });
};
