/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSecretDetails renders the key/value data of a KV secret. 
 * It also renders a dropdown to display different versions of the secret.
 * <Page::Secret::Details
 *  @secretPath={{this.model.path}}
 *  @secret={{this.model.secret}}
 *  @metadata={{this.model.metadata}}
 *  @breadcrumbs={{this.breadcrumbs}}
  /> 
 *
 * @param {string} secretPath - path of kv secret 'my/secret' used as the title for the KV page header 
 * @param {model} secret - Ember data model: 'kv/data'  
 * @param {model} metadata - Ember data model: 'kv/metadata'
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 * 
 * sample: 
  {
    backend: 'my-kv-engine',
    path: 'full/secret/path',
    secret: KvDataEmberModel,
    metadata: KvMetadataEmberModel,
  }
 *
 */

export default class KvSecretDetails extends Component {
  @tracked showJsonView = false;
  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }
}
