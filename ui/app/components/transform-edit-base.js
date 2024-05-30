/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import Component from '@ember/component';
import { set } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export const addToList = (list, itemToAdd) => {
  if (!list || !Array.isArray(list)) return list;
  list.push(itemToAdd);
  return list.uniq();
};

export const removeFromList = (list, itemToRemove) => {
  if (!list) return list;
  const index = list.indexOf(itemToRemove);
  if (index < 0) return list;
  const newList = list.removeAt(index, 1);
  return newList.uniq();
};

export default Component.extend(FocusOnInsertMixin, {
  store: service(),
  flashMessages: service(),
  router: service(),

  mode: null,
  onDataChange() {},
  onRefresh() {},
  model: null,
  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  init() {
    this._super(...arguments);
    this.set('backendType', 'transform');
  },

  willDestroyElement() {
    if (this.model && this.model.isError && !this.model.isDestroyed && !this.model.isDestroying) {
      this.model.rollbackAttributes();
    }
    this._super(...arguments);
  },

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  },

  modelPrefixFromType(modelType) {
    let modelPrefix = '';
    if (modelType && modelType.startsWith('transform/')) {
      modelPrefix = `${modelType.replace('transform/', '')}/`;
    }
    return modelPrefix;
  },

  listTabFromType(modelType) {
    let tab;
    if (modelType && modelType.startsWith('transform/')) {
      tab = `${modelType.replace('transform/', '')}`;
    }
    return tab;
  },

  persist(method, successCallback) {
    const model = this.model;
    return model[method]()
      .then(() => {
        successCallback(model);
      })
      .catch((e) => {
        model.set('displayErrors', e.errors);
        throw e;
      });
  },

  applyDelete(callback = () => {}) {
    const tab = this.listTabFromType(this.model.constructor.modelName);
    this.persist('destroyRecord', () => {
      this.hasDataChanges();
      callback();
      this.transitionToRoute(LIST_ROOT_ROUTE, { queryParams: { tab } });
    });
  },

  applyChanges(type, callback = () => {}) {
    const modelId = this.model.id || this.model.name; // transform comes in as model.name
    const modelPrefix = this.modelPrefixFromType(this.model.constructor.modelName);
    // prevent from submitting if there's no key
    // maybe do something fancier later
    if (type === 'create' && isBlank(modelId)) {
      return;
    }

    this.persist('save', () => {
      this.hasDataChanges();
      callback();
      this.transitionToRoute(SHOW_ROUTE, `${modelPrefix}${modelId}`);
    });
  },

  hasDataChanges() {
    this.onDataChange(this.model?.hasDirtyAttributes);
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      this.applyChanges(type);
    },

    setValue(key, event) {
      set(this.model, key, event.target.checked);
    },

    refresh() {
      this.onRefresh();
    },

    delete() {
      this.applyDelete();
    },
  },
});
