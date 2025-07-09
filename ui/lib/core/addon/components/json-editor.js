/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module JsonEditor
 *
 * @example
 * <JsonEditor @title="Policy" @value={{hash foo="bar"}} @viewportMargin={{100}} />
 *
 * @param {string} [title] - Name above codemirror view
 * @param {boolean} [showToolbar=true] - If false, toolbar and title are hidden
 * @param {string} [value] - a specific string the comes from codemirror. It's the value inside the codemirror display
 * @param {Function} [valueUpdated] - action to preform when you edit the codemirror value.
 * @param {Function} [onBlur] - action to preform when you focus out of codemirror.
 * @param {string} [helpText] - helper text.
 * @param {Object} [extraKeys] - Provides keyboard shortcut methods for things like saving on shift + enter.
 * @param {string} [mode] - The mode defined for styling
 * @param {Boolean} [readOnly] - Sets the view to readOnly, allowing for copying but no editing. It also hides the cursor. Defaults to false.
 * @param {String} [value] - Value within the display. Generally, a json string.
 * @param {string} [example] - Example to show when value is null -- when example is provided a restore action will render in the toolbar to clear the current value and show the example after input
 * @param {string} [container] - **REQUIRED if rendering within a modal** Selector string or element object of containing element, set the focused element as the container value. This is for the Hds::Copy::Button and to set `autoRefresh=true` so content renders https://hds-website-hashicorp.vercel.app/components/copy/button?tab=code
 * @param {Function} [onSetup] - action to preform when the codemirror editor is setup.
 *
 */

export default class JsonEditorComponent extends Component {
  _codemirrorEditor = null;

  get mode() {
    return this.args.mode ?? 'json';
  }

  get getShowToolbar() {
    return this.args.showToolbar ?? true;
  }

  get ariaLabel() {
    return this.args.title ?? 'JSON Editor';
  }

  @action
  onSetup(editor) {
    this._codemirrorEditor = editor;

    this.args.onSetup?.(editor);
  }

  @action
  onUpdate(...args) {
    if (!this.args.readOnly) {
      // catching a situation in which the user is not readOnly and has not provided a valueUpdated function to the instance
      this.args.valueUpdated(...args);
    }
  }

  @action
  restoreExample() {
    this._codemirrorEditor.dispatch({
      changes: [
        {
          from: 0,
          to: this._codemirrorEditor.state.doc.length,
          insert: this.args.example,
        },
      ],
    });
    this.args.valueUpdated(null, this._codemirrorEditor);
  }
}
