/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Model } from 'vault/app-types';

export default class ClusterModel extends Model {
  id: string;
  version: any;
  nodes: any;
  name: any;
  status: any;
  standby: any;
  type: any;
  license: any;
  hasChrootNamespace: any;
  replicationRedacted: any;
  get licenseExpiry(): any;
  get licenseState(): any;
  get needsInit(): any;
  get unsealed(): boolean;
  get sealed(): boolean;
  get leaderNode(): any;
  get sealThreshold(): any;
  get sealProgress(): any;
  get sealType(): any;
  get storageType(): any;
  get hcpLinkStatus(): any;
  get hasProgress(): boolean;
  get usingRaft(): boolean;
  mode: any;
  get allReplicationDisabled(): any;
  get anyReplicationEnabled(): any;
  dr: any;
  performance: any;
  rm: any;
  get drMode(): any;
  get replicationMode(): any;
  get replicationModeForDisplay(): 'Disaster Recovery' | 'Performance';
  get replicationIsInitializing(): boolean;
}
