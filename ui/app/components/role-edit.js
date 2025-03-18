/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';
import Component from '@ember/component';
import { set } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'core/utils/key-codes';
import { ROUTES } from 'vault/utils/routes';

/**
 * @type Class
 */
export default Component.extend(FocusOnInsertMixin, {
  router: service(),

  mode: null,
  emptyData: '{\n}',
  onDataChange() {},
  onRefresh() {},
  model: null,
  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  willDestroyElement() {
    if (this.model && this.model.isError && !this.model.isDestroyed && !this.model.isDestroying) {
      this.model.rollbackAttributes();
    }
    this._super(...arguments);
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
    if (e.keyCode !== keys.ESC || this.mode !== 'show') {
      return;
    }
    this.transitionToRoute(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT);
  },

  hasDataChanges() {
    this.onDataChange(this.model.hasDirtyAttributes);
  },

  persist(method, successCallback) {
    const model = this.model;
    return model[method]().then(() => {
      if (!model.isError) {
        successCallback(model);
      }
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      // all of the attributes with fieldValue:'id' are called `name`
      const modelId = this.model.id || this.model.name;
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && isBlank(modelId)) {
        return;
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_SHOW, modelId);
      });
    },

    setValue(key, event) {
      set(this.model, key, event.target.checked);
    },

    refresh() {
      this.onRefresh();
    },

    delete() {
      this.persist('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(ROUTES.VAULT_CLUSTER_SECRETS_BACKEND_LISTROOT);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.model, attr, JSON.parse(val));
      }
    },
  },
});
