/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { dropTask } from 'ember-concurrency';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import sortObjects from 'vault/utils/sort-objects';

export default class VaultClusterAccessMethodsController extends Controller {
  @service flashMessages;

  @tracked authMethodOptions = [];
  @tracked selectedAuthType = null;
  @tracked selectedAuthName = null;
  @tracked methodToDisable = null;

  queryParams = ['page, pageFilter'];

  page = 1;
  pageFilter = null;
  filter = null;

  // list returned by getter is sorted in template
  get authMethodList() {
    // return an options list to filter by engine type, ex: 'kv'
    if (this.selectedAuthType) {
      // check first if the user has also filtered by name.
      // names are individualized across type so you can't have the same name for an aws auth method and userpass.
      // this means it's fine to filter by first type and then name or just name.
      if (this.selectedAuthName) {
        return this.model.filter((method) => this.selectedAuthName === method.id);
      }
      // otherwise filter by auth type
      return this.model.filter((method) => this.selectedAuthType === method.type);
    }
    // return an options list to filter by auth name, ex: 'my-userpass'
    if (this.selectedAuthName) {
      return this.model.filter((method) => this.selectedAuthName === method.id);
    }
    // no filters, return full list
    return this.model;
  }

  get authMethodArrayByType() {
    const arrayOfAllAuthTypes = this.authMethodList.map((modelObject) => modelObject.type);
    // filter out repeated auth types (e.g. [userpass, userpass] => [userpass])
    const arrayOfUniqueAuthTypes = [...new Set(arrayOfAllAuthTypes)];

    return arrayOfUniqueAuthTypes.map((authType) => ({
      name: authType,
      id: authType,
    }));
  }

  get authMethodArrayByName() {
    return this.authMethodList.map((modelObject) => ({
      name: modelObject.id,
      id: modelObject.id,
    }));
  }

  @action
  filterAuthType([type]) {
    this.selectedAuthType = type;
  }

  @action
  filterAuthName([name]) {
    this.selectedAuthName = name;
  }

  @dropTask
  *disableMethod(method) {
    const { type, path } = method;
    try {
      yield method.destroyRecord();
      this.flashMessages.success(`The ${type} Auth Method at ${path} has been disabled.`);
    } catch (err) {
      this.flashMessages.danger(
        `There was an error disabling Auth Method at ${path}: ${err.errors.join(' ')}.`
      );
    } finally {
      this.methodToDisable = null;
    }
  }

  // template helper
  sortMethods = (methods) => sortObjects(methods.slice(), 'path');
}
