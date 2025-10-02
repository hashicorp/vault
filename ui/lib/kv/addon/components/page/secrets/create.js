/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvSecretCreate is used for creating the initial version of a secret
 *
 * <Page::Secrets::Create
 *    @form={{this.model.form}}
 *    @path={{this.model.path}}
 *    @backend={{this.model.backend}}
 *    @breadcrumbs={{this.breadcrumbs}}
 *  />
 *
 * @param {Form} form - kv form
 * @param {string} path - secret path
 * @param {string} backend - secret mount path
 * @param {array} breadcrumbs - breadcrumb objects to render in page header
 */

export default class KvSecretCreate extends Component {
  @tracked showJsonView = false;
  @tracked showMetadata = false;
}
