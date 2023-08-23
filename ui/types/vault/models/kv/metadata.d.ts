/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model from '@ember-data/model';

export default class KvSecretDataModel extends Model {
  backend: string;
  path: string;
  fullSecretPath: string;
  maxVersions: number;
  casRequired: boolean;
  deleteVersionAfter: string;
  customMetadata: object;
  createdTime: string;
  currentVersion: number;
  oldestVersion: number;
  updatedTime: string;
  versions: object;
  // apiPaths for capabilities
  dataPath: Promise<CapabilitiesModel>;
  metadataPath: Promise<CapabilitiesModel>;

  // Capabilities
  get canDeleteMetadata(): boolean;
  get canReadMetadata(): boolean;
  get canUpdateMetadata(): boolean;
  get canCreateVersionData(): boolean;
}
