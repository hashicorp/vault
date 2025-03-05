/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type RouterService from '@ember/routing/router-service';
import type LicenseModel from 'vault/models/license';
import type Controller from '@ember/controller';

export default class ClusterLicense extends Route {
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;
  @service declare readonly router: RouterService;

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster');
    }
  }

  model() {
    return hash({
      automatedSnapshot: this.fetchRaftSnapshotConfig(),
      license: this.store.queryRecord('license', {}),
      replication: this.fetchReplication(),
      namespaces: this.fetchNamespaceInternal(),
    });
  }

  fetchReplication() {
    // ARG TODO unsure what happens if not a feature included on model
    const clusterModel = this.modelFor('vault.cluster') as {
      hasChrootNamespace?: boolean;
      replicationRedacted?: boolean;
      dr?: object;
      performance?: object;
    };
    const hasChroot = clusterModel?.hasChrootNamespace;
    const replication =
      hasChroot || clusterModel.replicationRedacted
        ? null
        : {
            dr: clusterModel.dr,
            performance: clusterModel.performance,
          };
    return replication;
  }

  async fetchRaftSnapshotConfig() {
    // no adapter/serializer/model for raft snapshot config
    // if error, raft storage is not in use, otherwise permissions (must be root and sudo permissions)
    try {
      const adapter = this.store.adapterFor('application');
      await adapter.ajax('/v1/sys/storage/raft/snapshot-auto/config', 'GET', { data: { list: true } });
      // ARG TODO test with raftConfig
      // probably want to return some kind of data count...
      return { color: 'success', text: 'Enabled' };
    } catch (e: any) {
      // todo types on error
      if (e.httpStatus === 400 && e.errors[0] === `raft storage is not in use`) {
        return { color: 'warning', text: 'Not in use' };
      }
      return { color: 'critical', text: 'Permission error.' };
    }
  }

  async fetchNamespaceInternal() {
    // doing outside of model/adapter because of potentional future changes
    // if error, raft storage is not in use, otherwise permissions (must be root and sudo permissions)
    try {
      const adapter = this.store.adapterFor('application');
      const response = await adapter.ajax('/v1/sys/internal/ui/namespaces', 'GET');
      // probably want to return some kind of data count...
      return response;
    } catch (e: any) {
      // todo handle error
      return e;
    }
  }
}
