/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvPatchJsonForm
 * @description
 * This component renders one of two ways to patch a KV v2 secret (the other is using the KvPatch::Editor::Form).
 *
 * @example
 * <KvPatch::JsonForm @onSubmit={{perform this.save}} @onCancel={{this.onCancel}} @isSaving={{this.save.isRunning}} />
 *
 * @param {boolean} isSaving - if true, disables the save and cancel buttons. useful if the onSubmit callback is a concurrency task
 * @param {function} onCancel - called when form is canceled
 * @param {function} onSubmit - called when form is saved, called with with the key value object containing patch data
 * @param {object} subkeys - leaf keys of a kv v2 secret, all values (unless a nested object with more keys) return null. used for toggle that reveals codeblock of subkey structure
 * @param {string} submitError - error message string from parent if submit failed
 */

export default class KvPatchJsonForm extends Component {
  @tracked jsonObject;
  @tracked lintingErrors;

  constructor() {
    super(...arguments);
    // prefill JSON editor with an empty object
    this.jsonObject = JSON.stringify({ '': '' }, null, 2);
  }

  @action
  handleJson(value, codemirror) {
    codemirror.performLint();
    this.lintingErrors = codemirror.state.lint.marked.length > 0;
    if (!this.lintingErrors) {
      this.jsonObject = value;
    }
  }

  @action
  submit(event) {
    event.preventDefault();
    const patchData = JSON.parse(this.jsonObject);
    this.args.onSubmit(patchData);
  }
}
