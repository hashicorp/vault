/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

interface PathsResponse {
  [key: string]: {
    capabilities: string[];
  };
}
export default class PermissionsService extends Service {
  exactPaths: PathsResponse | null;
  globPaths: PathsResponse | null;
  canViewAll: boolean | null;
  permissionsBanner: string | null;
  chrootNamespace: string | null | undefined;
  hasNavPermission: (string) => boolean;
}
