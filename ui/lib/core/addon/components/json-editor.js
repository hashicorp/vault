/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { assert } from '@ember/debug';

/**
 * @module JsonEditor
 *
 * @example
 * <JsonEditor @title="Policy" @value={{hash foo="bar"}} @viewportMargin={{100}} />
 *
 * @param {string} [title] - Name above codemirror view
 * @param {string} [value] - a specific string the comes from codemirror. It's the value inside the codemirror display
 * @param {Function} [valueUpdated] - action to preform when you edit the codemirror value.
 * @param {Function} [onBlur] - action to preform when you focus out of codemirror.
 * @param {string} [helpText] - helper text.
 * @param {Object} [extraKeys] - Provides keyboard shortcut methods for things like saving on shift + enter.
 * @param {string} [mode] - The mode defined for styling
 * @param {Boolean} [readOnly] - Sets the view to readOnly, allowing for copying but no editing. It also hides the cursor. Defaults to false.
 * @param {String} [value] - Value within the display. Generally, a json string.
 * @param {string} [example] - Example to show when value is null -- when example is provided a restore action will render in the toolbar to clear the current value and show the example after input
 * @param {Function} [onSetup] - action to preform when the codemirror editor is setup.
 * @param {Function} [onRestoreExample] - override callback to customize  "Restore example" behavior. Default behavior is to reset @value to `null`
 *
 */

export default class JsonEditorComponent extends Component {
  _codemirrorEditor = null;

  constructor() {
    super(...arguments);

    const hasValueUpdated = !this.args.readOnly ? !!this.args.valueUpdated : true;
    assert('@valueUpdated callback is required when component is not @readOnly', hasValueUpdated);
  }

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
  restoreExample() {
    if (this.args.onRestoreExample) {
      // Override to reset the @value of the code editor to something other than `null`
      this.args.onRestoreExample();
    } else {
      // Display @example in the editor but reset @value to `null` because
      // sometimes @example is not valid to set and submit as the actual input value.
      this._codemirrorEditor.dispatch({
        changes: [
          {
            from: 0,
            to: this._codemirrorEditor.state.doc.length,
            insert: this.args.example,
          },
        ],
      });
      this.args.valueUpdated(null);
    }
  }
}
