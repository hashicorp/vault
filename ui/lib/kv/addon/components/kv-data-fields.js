/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import KVObject from 'vault/lib/kv-object';

/**
 * @module KvDataFields is used for rendering the fields associated with kv secret data, it hides/shows a json editor and renders validation errors for the json editor
 *
 * <KvDataFields
 *  @showJson={{true}}
 *  @secret={{@secret}}
 *  @isEdit={{true}}
 *  @modelValidations={{this.modelValidations}}
 *  @pathValidations={{this.pathValidations}}
 * />
 *
 * @param {model} secret - Ember data model: 'kv/data', the new record saved by the form
 * @param {boolean} showJson - boolean passed from parent to hide/show json editor
 * @param {object} [modelValidations] - object of errors.  If attr.name is in object and has error message display in AlertInline.
 * @param {callback} [pathValidations] - callback function fired for the path input on key up
 * @param {boolean} [isEdit=false] - if true, this is a new secret version rather than a new secret. Used to change text for some form labels
 */

export default class KvDataFields extends Component {
  @tracked lintingErrors;

  get emptyJson() {
    // if secretData is null, this specially formats a blank object and renders a nice initial state for the json editor
    return KVObject.create({ content: [{ name: '', value: '' }] }).toJSONString(true);
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.args.secret.secretData = JSON.parse(value);
    }
  }
}
