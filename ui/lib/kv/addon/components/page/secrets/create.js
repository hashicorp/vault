/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSecretsCreate renders the form for creating a new secret. 
 * 
 * <Page::Secrets::Create
 *  @secret={{this.model.secret}}
 *  @breadcrumbs={{this.breadcrumbs}}
  /> 
 *
 * @param {model} secret - Ember data model: 'kv/data'  
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretsCreate extends Component {
  @tracked showJsonView = false;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }
}
