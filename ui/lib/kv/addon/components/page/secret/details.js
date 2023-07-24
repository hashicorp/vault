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
 * @param {boolean} noMetadataPermission - True if we received a 403 from the kv/metadata network request
 * @param {array} breadcrumbs - Array to generate breadcrumbs, passed to the page header component
 */

export default class KvSecretDetails extends Component {
  @tracked showJsonView = false;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }

  get emptyState() {
    const { version, canReadData, destroyed, deletionTime } = this.args.secret;
    if (!canReadData) {
      return {
        title: 'You do not have permission to read this secret',
        message:
          'Your policies may permit you to write a new version of this secret, but do not allow you to read its current contents.',
      };
    }
    if (destroyed) {
      return {
        title: `Version ${version} of this secret has been permanently destroyed`,
        message:
          'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.',
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    if (deletionTime) {
      return {
        title: `Version ${version} of this secret has been deleted`,
        message:
          'This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above.',
        link: '/vault/docs/secrets/kv/kv-v2',
      };
    }
    return false;
  }
}
