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
 * @param {boolean} [obscureJson=false] - used to obfuscate json values in JsonEditor
 */

export default class KvDataFields extends Component {
  @tracked lintingErrors;
  @tracked codeMirrorString;

  constructor() {
    super(...arguments);
    this.codeMirrorString = this.args.secret?.secretData
      ? stringify([this.args.secret.secretData], {})
      : '{ "": "" }';
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.args.secret.secretData = JSON.parse(value);
    }
    this.codeMirrorString = value;
  }
}
