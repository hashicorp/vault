/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model from '@ember-data/model';

export default class KvSecretMetadataModel extends Model {
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

  get pathIsDirectory(): boolean;
  get isSecretDeleted(): boolean;
  get sortedVersions(): number[];
  get currentSecret(): { state: string; isDeactivated: boolean };

  // Capabilities
  get canDeleteMetadata(): boolean;
  get canReadMetadata(): boolean;
  get canUpdateMetadata(): boolean;
  get canCreateVersionData(): boolean;
}
