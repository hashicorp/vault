/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const ldapBreadcrumbs = (pathToItem: string, roleType: string, mountPath: string) => {
  const ancestry = pathToItem.split('/').filter(Boolean); // remove last empty string

  return ancestry.map((path: string, idx: number) => {
    const isLast = ancestry.length === idx + 1;
    if (isLast) return { label: path }; // last crumb is current view, so no route link

    // each segment is concatenation of ancestral paths
    // for example, if the full path to an item is "prod/admin/west"
    // the segments will be: prod, prod/admin, prod/admin/west
    // (without trailing forward slashes so the URL is sanitized).
    // LIST requests can then be made for each subdirectory crumb
    const segment = `${ancestry.slice(0, idx + 1).join('/')}`;
    const models = [mountPath, roleType, segment];
    return {
      label: path,
      route: 'roles.subdirectory',
      models,
    };
  });
};
