/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { stringify } from 'core/helpers/stringify';

/**
 * @module KvDataFields is used for rendering the fields associated with kv secret data, it hides/shows a json editor and renders validation errors for the json editor
 *
 * <KvDataFields
 *  @showJson={{true}}
 *  @secret={{@secret}}
 *  @type="edit"
 *  @modelValidations={{this.modelValidations}}
 *  @pathValidations={{this.pathValidations}}
 * />
 *
 * @param {model} secret - Ember data model: 'kv/data', the new record saved by the form
 * @param {boolean} showJson - boolean passed from parent to hide/show json editor
 * @param {object} [modelValidations] - object of errors.  If attr.name is in object and has error message display in AlertInline.
 * @param {callback} [pathValidations] - callback function fired for the path input on key up
 * @param {boolean} [type=null] - can be edit, create, or details. Used to change text for some form labels
 */

export default class KvDataFields extends Component {
  @tracked lintingErrors;
  @tracked codeMirrorSecretData;

  get startingValue() {
    // the value tracked by codemirror is always a JSON string.
    // The API will return a json blob which is stringified and passed to the codemirror editor.
    // Otherwise the default is a stringified object with an empty key value pair.
    return this.args.secret?.secretData
      ? stringify([this.args.secret.secretData], {})
      : JSON.stringify({ '': '' }, null, 2);
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.codeMirrorSecretData = value;
    }
  }

  @action
  onBlur() {
    // we parse the secretData and save it to the model only after the codemirror editor looses focus. To prevent issues caused by parsing differences between the codemirror editor and the namedArgs param, we do not parse the value on a keyPress, but only when the model.secretData needs to be saved. Examples: toggling between json editor and kv data fields, comparing a diff in create view, or saving the secretData.
    if (!this.args.secret) return;
    this.args.secret.secretData = JSON.parse(this.codeMirrorSecretData);
  }
}
