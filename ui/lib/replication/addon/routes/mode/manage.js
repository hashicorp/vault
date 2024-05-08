/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { camelize } from '@ember/string';
import { all } from 'rsvp';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { replicationActionForMode } from 'core/helpers/replication-action-for-mode';

const pathForAction = (action, replicationMode, clusterMode) => {
  let path;
  if (action === 'reindex' || action === 'recover') {
    path = `sys/replication/${action}`;
  } else {
    path = `sys/replication/${replicationMode}/${clusterMode}/${action}`;
  }
  return path;
};

export default Route.extend({
  router: service(),
  store: service(),
  model() {
    const store = this.store;
    const model = this.modelFor('mode');

    const replicationMode = this.paramsFor('mode').replication_mode;
    const clusterMode = model[replicationMode].modeForUrl;
    const actions = replicationActionForMode([replicationMode, clusterMode]);
    return all(
      actions.map((action) => {
        return store.findRecord('capabilities', pathForAction(action)).then((capability) => {
          model.set(`can${camelize(action)}`, capability.canUpdate);
        });
      })
    ).then(() => {
      return model;
    });
  },

  beforeModel() {
    const model = this.modelFor('mode');
    const replicationMode = this.paramsFor('mode').replication_mode;
    const modeModel = model[replicationMode];
    if (modeModel.replicationDisabled || modeModel.replicationUnsupported) {
      this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
