import Ember from 'ember';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';
import KVObject from 'vault/lib/kv-object';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
const { get, computed } = Ember;

export default Ember.Component.extend(FocusOnInsertMixin, {
  // a key model
  key: null,

  // a value to pre-fill the key input - this is populated by the corresponding
  // 'initialKey' queryParam
  initialKey: null,

  // set in the route's setupController hook
  mode: null,

  secretData: null,

  // called with a bool indicating if there's been a change in the secretData
  onDataChange: () => {},

  // did user request advanced mode
  preferAdvancedEdit: false,

  // use a named action here so we don't have to pass one in
  // this will bubble to the route
  toggleAdvancedEdit: 'toggleAdvancedEdit',

  codemirrorString: null,

  hasLintError: false,

  init() {
    this._super(...arguments);
    const secrets = this.get('key.secretData');
    const data = KVObject.create({ content: [] }).fromJSON(secrets);
    this.set('secretData', data);
    this.set('codemirrorString', data.toJSONString());
    if (data.isAdvanced()) {
      this.set('preferAdvancedEdit', true);
    }
    this.checkRows();
    if (this.get('mode') === 'edit') {
      this.send('addRow');
    }
  },

  willDestroyElement() {
    const key = this.get('key');
    if (get(key, 'isError') && !key.isDestroyed) {
      key.rollbackAttributes();
    }
  },

  partialName: Ember.computed('mode', function() {
    return `partials/secret-form-${this.get('mode')}`;
  }),

  routing: Ember.inject.service('-routing'),

  showPrefix: computed.or('key.initialParentKey', 'key.parentKey'),

  requestInFlight: computed.or('key.isLoading', 'key.isReloading', 'key.isSaving'),

  buttonDisabled: computed.or(
    'requestInFlight',
    'key.isFolder',
    'key.isError',
    'key.flagsIsInvalid',
    'hasLintError'
  ),

  basicModeDisabled: computed('secretDataIsAdvanced', 'showAdvancedMode', function() {
    return this.get('secretDataIsAdvanced') || this.get('showAdvancedMode') === false;
  }),

  secretDataAsJSON: computed('secretData', 'secretData.[]', function() {
    return this.get('secretData').toJSON();
  }),

  secretDataIsAdvanced: computed('secretData', 'secretData.[]', function() {
    return this.get('secretData').isAdvanced();
  }),

  hasDataChanges() {
    const keyDataString = this.get('key.dataAsJSONString');
    const sameData = this.get('secretData').toJSONString() === keyDataString;
    if (sameData === false) {
      this.set('lastChange', Date.now());
    }

    this.get('onDataChange')(!sameData);
  },

  showAdvancedMode: computed('preferAdvancedEdit', 'secretDataIsAdvanced', 'lastChange', function() {
    return this.get('secretDataIsAdvanced') || this.get('preferAdvancedEdit');
  }),

  transitionToRoute() {
    const router = this.get('routing.router');
    router.transitionTo.apply(router, arguments);
  },

  bindKeys: Ember.on('didInsertElement', function() {
    Ember.$(document).on('keyup.keyEdit', this.onEscape.bind(this));
  }),

  unbindKeys: Ember.on('willDestroyElement', function() {
    Ember.$(document).off('keyup.keyEdit');
  }),

  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.get('mode') !== 'show') {
      return;
    }
    const parentKey = this.get('key.parentKey');
    if (parentKey) {
      this.transitionToRoute(LIST_ROUTE, parentKey);
    } else {
      this.transitionToRoute(LIST_ROOT_ROUTE);
    }
  },

  // successCallback is called in the context of the component
  persistKey(method, successCallback, isCreate) {
    let model = this.get('key');
    let key = model.get('id');

    if (key.startsWith('/')) {
      key = key.replace(/^\/+/g, '');
      model.set('id', key);
    }

    if (isCreate && typeof model.createRecord === 'function') {
      // create an ember data model from the proxy
      model = model.createRecord(model.get('backend'));
      this.set('key', model);
    }

    return model[method]().then(() => {
      if (!Ember.get(model, 'isError')) {
        successCallback(key);
      }
    });
  },

  checkRows() {
    if (this.get('secretData').get('length') === 0) {
      this.send('addRow');
    }
  },

  actions: {
    handleKeyDown(e) {
      e.stopPropagation();
      if (!(e.keyCode === keys.ENTER && e.metaKey)) {
        return;
      }
      let $form = this.$('form');
      if ($form.length) {
        $form.submit();
      }
      $form = null;
    },

    handleChange() {
      this.set('codemirrorString', this.get('secretData').toJSONString(true));
      this.hasDataChanges();
    },

    createOrUpdateKey(type, event) {
      event.preventDefault();
      const newData = this.get('secretData').toJSON();
      this.get('key').set('secretData', newData);

      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && Ember.isBlank(this.get('key.id'))) {
        return;
      }

      this.persistKey(
        'save',
        key => {
          this.hasDataChanges();
          this.transitionToRoute(SHOW_ROUTE, key);
        },
        type === 'create'
      );
    },

    deleteKey() {
      this.persistKey('destroyRecord', () => {
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },

    refresh() {
      this.attrs.onRefresh();
    },

    addRow() {
      const data = this.get('secretData');
      if (Ember.isNone(data.findBy('name', ''))) {
        data.pushObject({ name: '', value: '' });
        this.set('codemirrorString', data.toJSONString(true));
      }
      this.checkRows();
      this.hasDataChanges();
    },

    deleteRow(name) {
      const data = this.get('secretData');
      const item = data.findBy('name', name);
      if (Ember.isBlank(item.name)) {
        return;
      }
      data.removeObject(item);
      this.checkRows();
      this.hasDataChanges();
      this.set('codemirrorString', data.toJSONString(true));
      this.rerender();
    },

    toggleAdvanced(bool) {
      this.sendAction('toggleAdvancedEdit', bool);
    },

    codemirrorUpdated(val, codemirror) {
      codemirror.performLint();
      const noErrors = codemirror.state.lint.marked.length === 0;
      if (noErrors) {
        this.get('secretData').fromJSONString(val);
      }
      this.set('hasLintError', !noErrors);
      this.set('codemirrorString', val);
    },

    formatJSON() {
      this.set('codemirrorString', this.get('secretData').toJSONString(true));
    },
  },
});
