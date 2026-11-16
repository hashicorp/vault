/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@ember/component';
import { task, waitForEvent } from 'ember-concurrency';

import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'core/utils/keys';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';

export default Component.extend(FocusOnInsertMixin, {
  router: service(),
  mode: null,
  onRefresh() {},

  get keyData() {
    return this.mode === 'show' ? this.key : this.form.data;
  },

  get breadcrumbs() {
    const baseCrumbs = [
      {
        label: 'Vault',
        icon: 'vault',
        route: 'vault.cluster.dashboard',
      },
      {
        label: 'Secrets',
        route: 'vault.cluster.secrets',
      },
      {
        label: this.keyData.backend,
        route: 'vault.cluster.secrets.backend.list-root',
        model: this.keyData.backend,
      },
    ];
    if (this.mode === 'show') {
      return [
        ...baseCrumbs,
        {
          label: this.keyData.id,
        },
      ];
    } else if (this.mode === 'edit') {
      return [
        ...baseCrumbs,
        {
          label: this.keyData.id,
          route: 'vault.cluster.secrets.backend.show',
          models: [this.keyData.backend, this.keyData.id],
          query: { tab: 'details' },
        },
        { label: 'Edit' },
      ];
    } else if (this.mode === 'create') {
      return [...baseCrumbs, { label: 'Create' }];
    }
    return baseCrumbs;
  },

  get title() {
    if (this.mode === 'create') {
      return 'Create Key';
    } else if (this.mode === 'edit') {
      return 'Edit Key';
    } else {
      return 'Key';
    }
  },

  get subtitle() {
    if (this.mode === 'create' || this.mode === 'edit') return '';

    return this.keyData.id;
  },

  waitForKeyUp: task(function* () {
    while (true) {
      const event = yield waitForEvent(document.body, 'keyup');
      this.onEscape(event);
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement'),

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  },

  onEscape(e) {
    const isEscKeyPressed = keys.ESC.includes(e.key);
    if (isEscKeyPressed || this.mode !== 'show') {
      return;
    }
    this.transitionToRoute(LIST_ROOT_ROUTE);
  },

  actions: {
    refresh() {
      this.onRefresh();
    },
  },
});
