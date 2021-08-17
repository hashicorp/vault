/**
 * @module SecretEdit
 * SecretEdit component manages the secret and model data, and displays either the create, update, empty state or show view of a KV secret.
 *
 * @example
 * ```js
 * <SecretEdit @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
//  * ARG TODO FINISH THIS
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { alias, or } from '@ember/object/computed';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import KVObject from 'vault/lib/kv-object';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';

export default Component.extend(FocusOnInsertMixin, WithNavToNearestAncestor, {
  wizard: service(),
  store: service(),

  // a key model
  key: null,
  model: null,

  // a value to pre-fill the key input - this is populated by the corresponding
  // 'initialKey' queryParam
  initialKey: null,

  // set in the route's setupController hook
  mode: null,

  secretData: null,

  // called with a bool indicating if there's been a change in the secretData and customMetadata
  onDataChange() {},
  onRefresh() {},
  onToggleAdvancedEdit() {},

  // did user request advanced mode
  preferAdvancedEdit: false,

  // use a named action here so we don't have to pass one in
  // this will bubble to the route
  toggleAdvancedEdit: 'toggleAdvancedEdit',
  error: null,

  codemirrorString: null,

  hasLintError: false,
  isV2: false,

  init() {
    this._super(...arguments);
    let secrets = this.model.secretData;
    if (!secrets && this.model.selectedVersion) {
      this.set('isV2', true);
      secrets = this.model.belongsTo('selectedVersion').value().secretData;
    }
    const data = KVObject.create({ content: [] }).fromJSON(secrets);
    this.set('secretData', data);

    this.set('codemirrorString', data.toJSONString());
    if (data.isAdvanced()) {
      this.set('preferAdvancedEdit', true);
    }
    if (this.wizard.featureState === 'details' && this.mode === 'create') {
      let engine = this.model.backend.includes('kv') ? 'kv' : this.model.backend;
      this.wizard.transitionFeatureMachine('details', 'CONTINUE', engine);
    }
  },

  checkSecretCapabilities: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.model || context.mode === 'create') {
        return;
      }
      let backend = context.isV2 ? context.get('model.engine.id') : context.model.backend;
      let id = context.model.id;
      let path = context.isV2 ? `${backend}/data/${id}` : `${backend}/${id}`;
      return {
        id: path,
      };
    },
    'isV2',
    'model',
    'model.id',
    'mode'
  ),
  canDeleteSecretData: alias('checkSecretCapabilities.canDelete'),
  canUpdateSecretData: alias('checkSecretCapabilities.canUpdate'),

  checkMetadataCapabilities: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.model || !context.isV2) {
        return;
      }
      let backend = context.model.backend;
      let path = `${backend}/metadata/`;
      return {
        id: path,
      };
    },
    'isV2',
    'model',
    'model.id',
    'mode'
  ),
  canDeleteSecretMetadata: alias('checkMetadataCapabilities.canDelete'),
  canUpdateSecretMetadata: alias('checkMetadataCapabilities.canUpdate'),
  canCreateSecretMetadata: alias('checkMetadataCapabilities.canCreate'),

  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  buttonDisabled: or('requestInFlight', 'model.isFolder', 'model.flagsIsInvalid', 'hasLintError', 'error'),

  modelForData: computed('isV2', 'model', function() {
    let { model } = this;
    if (!model) return null;
    return this.isV2 ? model.belongsTo('selectedVersion').value() : model;
  }),

  basicModeDisabled: computed('secretDataIsAdvanced', 'showAdvancedMode', function() {
    return this.secretDataIsAdvanced || this.showAdvancedMode === false;
  }),

  secretDataAsJSON: computed('secretData', 'secretData.[]', function() {
    return this.secretData.toJSON();
  }),

  secretDataIsAdvanced: computed('secretData', 'secretData.[]', function() {
    return this.secretData.isAdvanced();
  }),

  showAdvancedMode: or('secretDataIsAdvanced', 'preferAdvancedEdit'),

  isWriteWithoutRead: computed('model.failedServerRead', 'modelForData.failedServerRead', 'isV2', function() {
    if (!this.model) return;
    // if the version couldn't be read from the server
    if (this.isV2 && this.modelForData.failedServerRead) {
      return true;
    }
    // if the model couldn't be read from the server
    if (!this.isV2 && this.model.failedServerRead) {
      return true;
    }
    return false;
  }),

  actions: {
    // ARG TODO couldn't find this being used anywhere
    // deleteKey() {
    //   let { id } = this.model;
    //   this.model.destroyRecord().then(() => {
    //     this.navToNearestAncestor.perform(id);
    //   });
    // },

    refresh() {
      this.onRefresh();
    },

    toggleAdvanced(bool) {
      this.onToggleAdvancedEdit(bool);
    },
  },
});
