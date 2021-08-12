/**
 * @module SecretCreateOrUpdate
 * SecretCreateOrUpdate component displays either the form for creating a new secret or creating a new version of the secret
 *
 * @example
 * ```js
 * <SecretCreateOrUpdate @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
import ControlGroupError from 'vault/lib/control-group-error';
import Ember from 'ember';
import keys from 'vault/lib/keycodes';
import { inject as service } from '@ember/service';
import { isBlank, isNone } from '@ember/utils';
import { set } from '@ember/object';
import { task, waitForEvent } from 'ember-concurrency';

const LIST_ROUTE = 'vault.cluster.secrets.backend.list';
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend({
  secretPaths: null,
  validationErrorCount: 0,
  validationMessages: null,

  controlGroup: service(),
  router: service(),
  store: service(),
  wizard: service(),

  init() {
    this._super(...arguments);
    this.set('validationMessages', {
      path: '',
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
    this.checkRows();
    if (this.mode === 'edit') {
      this.send('addRow');
    }
  },

  checkRows() {
    if (this.secretData.length === 0) {
      this.send('addRow');
    }
  },
  checkValidation(name, value) {
    if (name === 'path') {
      !value
        ? set(this.validationMessages, name, `${name} can't be blank.`)
        : set(this.validationMessages, name, '');
    }
    // check duplicate on path
    if (name === 'path' && value) {
      this.secretPaths?.includes(value)
        ? set(this.validationMessages, name, `A secret with this ${name} already exists.`)
        : set(this.validationMessages, name, '');
    }
    // check maxVersions is a number
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
  transitionToRoute() {
    return this.router.transitionTo(...arguments);
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

  actions: {
    addRow() {
      const data = this.secretData;
      // fired off on init
      if (isNone(data.findBy('name', ''))) {
        data.pushObject({ name: '', value: '' });
        this.send('handleChange');
      }
      this.checkRows();
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
    createOrUpdateKey(type, event) {
      event.preventDefault();
      if (type === 'create' && isBlank(this.modelForData.path || this.modelForData.id)) {
        this.checkValidation('path', '');
        return;
      }

      this.persistKey(() => {
        this.transitionToRoute(SHOW_ROUTE, this.model.path || this.model.id);
      });
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
    formatJSON() {
      this.set('codemirrorString', this.secretData.toJSONString(true));
    },
    handleChange() {
      this.set('codemirrorString', this.secretData.toJSONString(true));
      set(this.modelForData, 'secretData', this.secretData.toJSON());
    },

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
  },
});
