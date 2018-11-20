import { inject as service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';
import Component from '@ember/component';
import { set, get } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend(FocusOnInsertMixin, {
  router: service(),
  wizard: service(),

  mode: null,
  emptyData: '{\n}',
  onDataChange() {},
  onRefresh() {},
  model: null,
  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  didReceiveAttrs() {
    this._super(...arguments);
    if (
      (this.get('wizard.featureState') === 'details' && this.get('mode') === 'create') ||
      (this.get('wizard.featureState') === 'role' && this.get('mode') === 'show')
    ) {
      this.get('wizard').transitionFeatureMachine(
        this.get('wizard.featureState'),
        'CONTINUE',
        this.get('backendType')
      );
    }
    if (this.get('wizard.featureState') === 'displayRole') {
      this.get('wizard').transitionFeatureMachine(
        this.get('wizard.featureState'),
        'NOOP',
        this.get('backendType')
      );
    }
  },

  willDestroyElement() {
    this._super(...arguments);
    const model = this.get('model');
    if (get(model, 'isError')) {
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

  transitionToRoute() {
    this.get('router').transitionTo(...arguments);
  },

  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.get('mode') !== 'show') {
      return;
    }
    this.transitionToRoute(LIST_ROOT_ROUTE);
  },

  hasDataChanges() {
    get(this, 'onDataChange')(get(this, 'model.hasDirtyAttributes'));
  },

  persist(method, successCallback) {
    const model = get(this, 'model');
    return model[method]().then(() => {
      if (!get(model, 'isError')) {
        if (this.get('wizard.featureState') === 'role') {
          this.get('wizard').transitionFeatureMachine('role', 'CONTINUE', this.get('backendType'));
        }
        successCallback(model);
      }
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      const modelId = this.get('model.id');
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && isBlank(modelId)) {
        return;
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    },

    setValue(key, event) {
      set(get(this, 'model'), key, event.target.checked);
    },

    refresh() {
      this.get('onRefresh')();
    },

    delete() {
      this.persist('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.get('model'), attr, JSON.parse(val));
      }
    },
  },
});
