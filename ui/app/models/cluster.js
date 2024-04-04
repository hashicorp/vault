/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr, belongsTo, hasMany } from '@ember-data/model';
import { service } from '@ember/service';
import { get } from '@ember/object';

export default class ClusterModel extends Model {
  @service version;

  @hasMany('nodes', { async: false, inverse: null }) nodes;
  @attr('string') name;
  @attr('string') status;
  @attr('boolean') standby;
  @attr('string') type;
  @attr('object') license;

  // manually set on response in cluster adapter
  @attr('boolean') hasChrootNamespace;
  @attr('boolean') replicationRedacted;

  /* Licensing concerns */
  get licenseExpiry() {
    return this.license?.expiry_time;
  }
  get licenseState() {
    return this.license?.state;
  }

  get needsInit() {
    return this.nodes.every((node) => {
      return node.initialized === false;
    });
  }

  get unsealed() {
    return !!this.nodes.find((node) => {
      return node.sealed === false;
    });
  }

  get sealed() {
    return !this.unsealed;
  }

  get leaderNode() {
    const nodes = this.nodes;
    if (nodes.length === 1) {
      return nodes[0];
    } else {
      return nodes.find((node) => node.isLeader === true);
    }
  }

  get sealThreshold() {
    return this.leaderNode?.sealThreshold;
  }
  get sealProgress() {
    return this.leaderNode?.progress;
  }
  get sealType() {
    return this.leaderNode?.type;
  }
  get storageType() {
    return this.leaderNode?.storageType;
  }
  get hcpLinkStatus() {
    return this.leaderNode?.hcpLinkStatus;
  }
  get hasProgress() {
    return this.sealProgress >= 1;
  }
  get usingRaft() {
    return this.storageType === 'raft';
  }

  //replication mode - will only ever be 'unsupported'
  //otherwise the particular mode will have the relevant mode attr through replication-attributes
  // eg dr.mode or performance.mode
  @attr('string')
  mode;
  get allReplicationDisabled() {
    return this.dr?.replicationDisabled && this.performance?.replicationDisabled;
  }
  get anyReplicationEnabled() {
    return this.dr?.replicationEnabled || this.performance?.replicationEnabled;
  }

  @belongsTo('replication-attributes', { async: false, inverse: null }) dr;
  @belongsTo('replication-attributes', { async: false, inverse: null }) performance;
  // this service exposes what mode the UI is currently viewing
  // replicationAttrs will then return the relevant `replication-attributes` model
  @service('replication-mode') rm;
  get drMode() {
    return this.dr.mode;
  }
  get replicationMode() {
    return this.rm.mode;
  }
  get replicationModeForDisplay() {
    return this.replicationMode === 'dr' ? 'Disaster Recovery' : 'Performance';
  }
  get replicationIsInitializing() {
    // a mode of null only happens when a cluster is being initialized
    // otherwise the mode will be 'disabled', 'primary', 'secondary'
    return !this.dr?.mode || !this.performance?.mode;
  }
  get replicationAttrs() {
    const replicationMode = this.replicationMode;
    return replicationMode ? get(this, replicationMode) : null;
  }
}
