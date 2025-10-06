/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

import type { PathInfo } from 'vault/utils/openapi-helpers';

export default class PathHelpService extends Service {
  getPaths(apiPath: string, backend: string, itemType?: string, itemID?: string): Promise<PathInfo>;
  hydrateModel(modelType: string, backend: string): Promise<void>;
  getNewModel(modelType: string, backend: string, apiPath: string, itemType?: string): Promise<void>;
}
