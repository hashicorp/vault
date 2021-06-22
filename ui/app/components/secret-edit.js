import Ember from 'ember';
import { isBlank, isNone } from '@ember/utils';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed, set } from '@ember/object';
import { alias, or, not } from '@ember/object/computed';
import { task, waitForEvent } from 'ember-concurrency';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import keys from 'vault/lib/keycodes';
import KVObject from 'vault/lib/kv-object';
import { maybeQueryRecord } from 'vault/macros/maybe-query-record';
import ControlGroupError from 'vault/lib/control-group-error';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend(FocusOnInsertMixin, WithNavToNearestAncestor, {
  wizard: service(),
  controlGroup: service(),
  router: service(),
  store: service(),
  flashMessages: service(),

  // a key model
  key: null,
  model: null,

  // a value to pre-fill the key input - this is populated by the corresponding
  // 'initialKey' queryParam
  initialKey: null,

  // set in the route's setupController hook
  mode: null,

  secretData: null,

  wrappedData: null,
  isWrapping: false,
  showWrapButton: not('wrappedData'),

  // called with a bool indicating if there's been a change in the secretData
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

  // cp-validation related properties
  validationMessages: null,
  validationErrorCount: 0,

  secretPaths: null,

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
    this.checkRows();
    if (this.wizard.featureState === 'details' && this.mode === 'create') {
      let engine = this.model.backend.includes('kv') ? 'kv' : this.model.backend;
      this.wizard.transitionFeatureMachine('details', 'CONTINUE', engine);
    }
    if (this.mode === 'edit') {
      this.send('addRow');
    }
    this.set('validationMessages', {
      path: '',
      key: '',
      maxVersions: '',
    });
    // for validation, return array of path names already assigned
    if (Ember.testing) {
      this.set('secretPaths', ['beep', 'bop', 'boop']);
    } else {
      let adapter = this.store.adapterFor('secret-v2');
      let type = { modelName: 'secret-v2' };
      let query = { backend: this.model.backend };
      adapter.query(this.store, type, query).then(result => {
        this.set('secretPaths', result.data.keys);
      });
    }
  },

  waitForKeyUp: task(function*(name, value) {
    this.checkValidation(name, value);
    while (true) {
      let event = yield waitForEvent(document.body, 'keyup');
      this.onEscape(event);
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement'),

  updatePath: maybeQueryRecord(
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
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),

  v2UpdatePath: maybeQueryRecord(
    'capabilities',
    context => {
      if (!context.model || context.mode === 'create' || context.isV2 === false) {
        return;
      }
      let backend = context.get('model.engine.id');
      let id = context.model.id;
      return {
        id: `${backend}/metadata/${id}`,
      };
    },
    'isV2',
    'model',
    'model.id',
    'mode'
  ),
  canEditV2Secret: alias('v2UpdatePath.canUpdate'),

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

  transitionToRoute() {
    return this.router.transitionTo(...arguments);
  },

  checkValidation(name, value) {
    // because path and key are not on the model performing custom validations instead of cp-validations
    if (name === 'path' || name === 'key') {
      // no value indicates missing presence
      !value
        ? set(this.validationMessages, name, `${name} can't be blank.`)
        : set(this.validationMessages, name, '');

      this.secretPaths.includes(value)
        ? set(this.validationMessages, name, `A secret with this ${name} already exists.`)
        : set(this.validationMessages, name, '');
    }
    if (name === 'maxVersions') {
      // checking for value because value which is blank on first load. No keyup event has occurred and default is 10.
      if (value) {
        let number = Number(value);
        this.model.set('maxVersions', number);
      }
      if (!this.model.validations.attrs.maxVersions.isValid) {
        set(this.validationMessages, name, this.model.validations.attrs.maxVersions.message);
      } else {
        set(this.validationMessages, name, '');
      }
    }

    let values = Object.values(this.validationMessages);

    this.set('validationErrorCount', values.filter(Boolean).length);
  },

  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.mode !== 'show') {
      return;
    }
    const parentKey = this.model.parentKey;
    if (parentKey) {
      this.transitionToRoute(LIST_ROUTE, parentKey);
    } else {
      this.transitionToRoute(LIST_ROOT_ROUTE);
    }
  },

  // successCallback is called in the context of the component
  persistKey(successCallback) {
    let secret = this.model;
    let secretData = this.modelForData;
    let isV2 = this.isV2;
    let key = secretData.get('path') || secret.id;

    if (key.startsWith('/')) {
      key = key.replace(/^\/+/g, '');
      secretData.set(secretData.pathAttr, key);
    }

    if (this.mode === 'create') {
      key = JSON.stringify({
        backend: secret.backend,
        id: key,
      });
    }

    return secretData
      .save()
      .then(() => {
        if (!secretData.isError) {
          if (isV2) {
            secret.set('id', key);
          }
          if (isV2 && Object.keys(secret.changedAttributes()).length) {
            // save secret metadata
            secret
              .save()
              .then(() => {
                this.saveComplete(successCallback, key);
              })
              .catch(e => {
                this.set(e, e.errors.join(' '));
              });
          } else {
            this.saveComplete(successCallback, key);
          }
        }
      })
      .catch(error => {
        if (error instanceof ControlGroupError) {
          let errorMessage = this.controlGroup.logFromError(error);
          this.set('error', errorMessage.content);
        }
        throw error;
      });
  },
  saveComplete(callback, key) {
    if (this.wizard.featureState === 'secret') {
      this.wizard.transitionFeatureMachine('secret', 'CONTINUE');
    }
    callback(key);
  },

  checkRows() {
    if (this.secretData.length === 0) {
      this.send('addRow');
    }
  },

  actions: {
    //submit on shift + enter
    handleKeyDown(e) {
      e.stopPropagation();
      if (!(e.keyCode === keys.ENTER && e.metaKey)) {
        return;
      }
      let $form = this.element.querySelector('form');
      if ($form.length) {
        $form.submit();
      }
    },

    handleChange() {
      this.set('codemirrorString', this.secretData.toJSONString(true));
      set(this.modelForData, 'secretData', this.secretData.toJSON());
    },

    handleWrapClick() {
      this.set('isWrapping', true);
      if (this.isV2) {
        this.store
          .adapterFor('secret-v2-version')
          .queryRecord(this.modelForData.id, { wrapTTL: 1800 })
          .then(resp => {
            this.set('wrappedData', resp.wrap_info.token);
            this.flashMessages.success('Secret Successfully Wrapped!');
          })
          .catch(() => {
            this.flashMessages.danger('Could Not Wrap Secret');
          })
          .finally(() => {
            this.set('isWrapping', false);
          });
      } else {
        this.store
          .adapterFor('secret')
          .queryRecord(null, null, { backend: this.model.backend, id: this.modelForData.id, wrapTTL: 1800 })
          .then(resp => {
            this.set('wrappedData', resp.wrap_info.token);
            this.flashMessages.success('Secret Successfully Wrapped!');
          })
          .catch(() => {
            this.flashMessages.danger('Could Not Wrap Secret');
          })
          .finally(() => {
            this.set('isWrapping', false);
          });
      }
    },

    clearWrappedData() {
      this.set('wrappedData', null);
    },

    handleCopySuccess() {
      this.flashMessages.success('Copied Wrapped Data!');
      this.send('clearWrappedData');
    },

    handleCopyError() {
      this.flashMessages.danger('Could Not Copy Wrapped Data');
      this.send('clearWrappedData');
    },

    createOrUpdateKey(type, event) {
      event.preventDefault();
      let model = this.modelForData;
      let arraySecretKeys = Object.keys(model.secretData);
      if (type === 'create' && isBlank(model.path || model.id)) {
        this.checkValidation('path', '');
        return;
      }
      if (arraySecretKeys.includes('')) {
        this.checkValidation('key', '');
        return;
      }
      this.persistKey(key => {
        // ARG TODO issue with saving no value but key and transition.
        // I'm at a loss. Need better handler of JSON.parse and catching the error.
        console.log(key, 'KEY!!!');
        let secretKey;
        try {
          secretKey = typeof key === 'object' ? JSON.parse(key).id : key;
          typeof secretKey === 'number' ? secretKey.toString() : secretKey;
        } catch {
          secretKey = key;
        }
        this.transitionToRoute(SHOW_ROUTE, secretKey);
      });
    },

    deleteKey() {
      let { id } = this.model;
      this.model.destroyRecord().then(() => {
        this.navToNearestAncestor.perform(id);
      });
    },

    refresh() {
      this.onRefresh();
    },

    addRow() {
      const data = this.secretData;
      if (isNone(data.findBy('name', ''))) {
        data.pushObject({ name: '', value: '' });
        this.send('handleChange');
      }
      this.checkRows();
    },

    deleteRow(name) {
      const data = this.secretData;
      const item = data.findBy('name', name);
      if (isBlank(item.name)) {
        return;
      }
      data.removeObject(item);
      this.checkRows();
      this.send('handleChange');
    },

    toggleAdvanced(bool) {
      this.onToggleAdvancedEdit(bool);
    },

    codemirrorUpdated(val, codemirror) {
      this.set('error', null);
      codemirror.performLint();
      const noErrors = codemirror.state.lint.marked.length === 0;
      if (noErrors) {
        try {
          this.secretData.fromJSONString(val);
          set(this.modelForData, 'secretData', this.secretData.toJSON());
        } catch (e) {
          this.set('error', e.message);
        }
      }
      this.set('hasLintError', !noErrors);
      this.set('codemirrorString', val);
    },

    formatJSON() {
      this.set('codemirrorString', this.secretData.toJSONString(true));
    },
  },
});
