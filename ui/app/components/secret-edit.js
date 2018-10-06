import { or } from '@ember/object/computed';
import { isBlank, isNone } from '@ember/utils';
import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed, get } from '@ember/object';
import { task, waitForEvent } from 'ember-concurrency';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';
import KVObject from 'vault/lib/kv-object';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend(FocusOnInsertMixin, {
  wizard: service(),
  router: service(),

  // a key model
  key: null,
  model: null,

  // a value to pre-fill the key input - this is populated by the corresponding
  // 'initialKey' queryParam
  initialKey: null,

  // set in the route's setupController hook
  mode: null,

  secretData: null,

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
    if (this.get('wizard.featureState') === 'details' && this.get('mode') === 'create') {
      let engine = this.get('model').backend.includes('kv') ? 'kv' : this.get('model').backend;
      this.get('wizard').transitionFeatureMachine('details', 'CONTINUE', engine);
    }

    if (this.mode === 'edit') {
      this.send('addRow');
    }
  },

  willDestroyElement() {
    this._super(...arguments);
    if (this.model.isError && !this.model.isDestroyed) {
      model.rollbackAttributes();
    }
  },

  waitForKeyUp: task(function*() {
    while (true) {
      let event = yield waitForEvent(document.body, 'keyup');
      this.onEscape(event);
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement'),

  partialName: computed('mode', function() {
    return `partials/secret-form-${this.get('mode')}`;
  }),

  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  buttonDisabled: or(
    'requestInFlight',
    'model.isFolder',
    'model.isError',
    'model.flagsIsInvalid',
    'hasLintError',
    'error'
  ),

  modelForData: computed('isV2', 'model', function() {
    return this.isV2 ? this.model.belongsTo('selectedVersion').value() : this.model;
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

  showAdvancedMode: computed('preferAdvancedEdit', 'secretDataIsAdvanced', 'lastChange', function() {
    return this.secretDataIsAdvanced || this.preferAdvancedEdit;
  }),

  transitionToRoute() {
    this.router.transitionTo(...arguments);
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
    let model = this.modelForData;
    let key = model.get('path') || model.id;

    if (key.startsWith('/')) {
      key = key.replace(/^\/+/g, '');
      model.set(model.pathAttr, key);
    }

    return model.save().then(() => {
      if (!model.isError) {
        if (this.wizard.featureState === 'secret') {
          this.wizard.transitionFeatureMachine('secret', 'CONTINUE');
        }
        successCallback(key);
      }
    });
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
      console.log('form is: ', $form);
      if ($form.length) {
        $form.submit();
      }
    },

    handleChange() {
      this.set('codemirrorString', this.secretData.toJSONString(true));
    },

    createOrUpdateKey(type, event) {
      event.preventDefault();
      const newData = this.secretData.toJSON();
      let model = this.modelForData;
      model.set('secretData', newData);

      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && isBlank(model.get('path') || model.id)) {
        return;
      }

      this.persistKey(key => {
        this.transitionToRoute(SHOW_ROUTE, key);
      });
    },

    deleteKey() {
      this.model.destroyRecord().then(() => {
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },

    refresh() {
      this.onRefresh();
    },

    addRow() {
      const data = this.secretData;
      if (isNone(data.findBy('name', ''))) {
        data.pushObject({ name: '', value: '' });
        this.set('codemirrorString', data.toJSONString(true));
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
      this.set('codemirrorString', data.toJSONString(true));
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
