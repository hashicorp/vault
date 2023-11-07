/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import { tracked } from '@glimmer/tracking';
import { allSupportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';
import { service } from '@ember/service';

class AuthState {
  @tracked type = '';
  @tracked token = '';
  @tracked username = '';
  @tracked password = '';
  @tracked role = '';
  @tracked jwt = '';

  resetFields() {
    this.token = '';
    this.userame = '';
    this.password = '';
    this.role = '';
    this.jwt = '';
  }

  constructor(type) {
    this.type = type || 'token';
  }
}

export function withAuthForm(mountType) {
  return function decorator(SuperClass) {
    if (!Object.prototype.isPrototypeOf.call(Component, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withAuthForm decorator must be used on instance of a glimmer component class. Decorator not applied to returned class'
      );
      return SuperClass;
    }
    return class AuthFormComponent extends SuperClass {
      @service session;
      @tracked namespace = '';
      @tracked mountPath = '';
      @tracked error = '';
      @tracked state = new AuthState();
      static _type;

      constructor() {
        super(...arguments);
        if (!mountType || typeof mountType !== 'string') {
          throw new Error('must pass mount type as string');
        }
        this._type = mountType;
      }

      get showFields() {
        const backend = allSupportedAuthBackends().findBy('type', this._type);
        return backend.formAttributes;
      }

      maybeMask = (field) => {
        if (field === 'token' || field === 'password') {
          return 'password';
        }
        return 'text';
      };

      @action
      handleFormChange(evt) {
        this.state[evt.target.name] = evt.target.value;
      }

      @action
      async handleFormLogin(evt) {
        evt.preventDefault();
        const authenticator = `authenticator:${this._type}`;
        const fields = this.showFields.reduce((obj, field) => {
          obj[field] = this.state[field];
          return obj;
        }, {});

        try {
          await this.session.authenticate(authenticator, fields, {
            backend: this.mountPath,
            namespace: this.namespace,
          });
        } catch (e) {
          this.error = errorMessage(e);
        }

        if (this.session.isAuthenticated && this.args.onSuccess) {
          this.args.onSuccess();
        }
      }
    };
  };
}
