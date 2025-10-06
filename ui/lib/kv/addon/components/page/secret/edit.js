/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { isAdvancedSecret } from 'core/utils/advanced-secret';

/**
 * @module KvSecretEdit is used for creating a new version of a secret
 *
 * <Page::Secret::Edit
 *  @form={{this.model.form}}
 *  @secret={{this.model.newVersion}}
 *  @metadata={{this.model.metadata}}
 *  @path={{this.model.path}}
 *  @backend={{this.model.backend}}
 *  @breadcrumbs={{this.breadcrumbs}
 * />
 *
 * @param {Form} form - kv form
 * @param {object} secret - secret data
 * @param {object} metadata - secret metadata
 * @param {string} path - secret path
 * @param {string} backend - secret mount path
 * @param {array} breadcrumbs - breadcrumb objects to render in page header
 */

/* eslint-disable no-undef */
export default class KvSecretEdit extends Component {
  @tracked showJsonView = false;
  @tracked showDiff = false;
  @tracked updatedSecret;

  constructor() {
    super(...arguments);
    this.originalSecret = JSON.stringify(this.args.form.data.secretData || {});
    this.updatedSecret = this.args.form.data.secretData || {};
    if (isAdvancedSecret(this.originalSecret)) {
      // Default to JSON view if advanced
      this.showJsonView = true;
    }
  }

  get showOldVersionAlert() {
    const { secret, metadata } = this.args;
    // isNew check prevents alert from flashing after save but before route transitions
    if (metadata?.current_version && secret?.version) {
      return metadata.current_version !== secret.version;
    }
    return false;
  }

  get diffDelta() {
    const oldData = JSON.parse(this.originalSecret);
    const diffpatcher = jsondiffpatch.create({});
    return diffpatcher.diff(oldData, this.updatedSecret);
  }

  get visualDiff() {
    if (this.showDiff) {
      return this.diffDelta
        ? jsondiffpatch.formatters.html.format(this.diffDelta, this.updatedSecret)
        : JSON.stringify(this.updatedSecret, undefined, 2);
    }
    return null;
  }

  @action
  onSecretDataUpdate(value) {
    this.updatedSecret = value;
  }
}
