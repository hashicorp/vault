/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import type { Breadcrumb } from 'vault/vault/app-types';

export const ldapBreadcrumbs = (
  fullPath: string | undefined, // i.e. path/to/item
  roleType: string,
  mountPath: string,
  lastItemCurrent = false // this array of objects can be spread anywhere within the crumbs array
): Breadcrumb[] => {
  if (!fullPath) return [];
  const ancestry = fullPath.split('/').filter((path) => path !== '');
  const isDirectory = fullPath.endsWith('/');

  return ancestry.map((name: string, idx: number) => {
    const isLast = ancestry.length === idx + 1;
    // if the end of the path is the current route, don't return a route link
    if (isLast && lastItemCurrent) return { label: name };

    // each segment is a continued concatenation of ancestral paths.
    // for example, if the full path to an item is "prod/admin/west"
    // the segments will be: prod/, prod/admin/, prod/admin/west.
    // LIST or GET requests can then be made for each crumb accordingly.
    const segment = ancestry.slice(0, idx + 1).join('/');

    const itemPath = isLast && !isDirectory ? segment : `${segment}/`;
    const routeParams = [mountPath, roleType, itemPath];
    return {
      label: name,
      route: isLast && !isDirectory ? 'roles.role.details' : 'roles.subdirectory',
      models: routeParams,
    };
  });
};
