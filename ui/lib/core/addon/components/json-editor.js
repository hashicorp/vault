/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { stringify } from 'core/helpers/stringify';
import { obfuscateData } from 'core/utils/advanced-secret';

/**
 * @module JsonEditor
 *
 * @example
 * ```js
 * <JsonEditor @title="Policy" @value={{codemirror.string}} @valueUpdated={{ action "codemirrorUpdate"}} />
 * ```
 *
 * @param {string} [title] - Name above codemirror view
 * @param {string} value - a specific string the comes from codemirror. It's the value inside the codemirror display
 * @param {Function} [valueUpdated] - action to preform when you edit the codemirror value.
 * @param {Function} [onFocusOut] - action to preform when you focus out of codemirror.
 * @param {string} [helpText] - helper text.
 * @param {Object} [extraKeys] - Provides keyboard shortcut methods for things like saving on shift + enter.
 * @param {Array} [gutters] - An array of CSS class names or class name / CSS string pairs, each of which defines a width (and optionally a background), and which will be used to draw the background of the gutters.
 * @param {string} [mode] - The mode defined for styling. Right now we only import ruby so mode must but be ruby or defaults to javascript. If you wanted another language you need to import it into the modifier.
 * @param {Boolean} [readOnly] - Sets the view to readOnly, allowing for copying but no editing. It also hides the cursor. Defaults to false.
 * @param {String} [theme] - Specify or customize the look via a named "theme" class in scss.
 * @param {String} [value] - Value within the display. Generally, a json string.
 * @param {String} [viewportMargin] - Size of viewport. Often set to "Infinity" to load/show all text regardless of length.
 * @param {string} [example] - Example to show when value is null -- when example is provided a restore action will render in the toolbar to clear the current value and show the example after input
 * * REQUIRED if rendering within a modal *
 * @container gives context for the <Hd::Copy::Button> and sets autoRefresh=true so JsonEditor renders content (without this property @value only renders if editor is focused)
 * @param {string} [container] - Selector string or element object of containing element, set the focused element as the container value. This is for the Hds::Copy::Button and to set autoRefresh=true so content renders https://hds-website-hashicorp.vercel.app/components/copy/button?tab=code
 *
 */

export default class JsonEditorComponent extends Component {
  @tracked revealValues = false;
  get getShowToolbar() {
    return this.args.showToolbar === false ? false : true;
  }

  get showObfuscatedData() {
    return this.args.readOnly && this.args.obscure && !this.revealValues;
  }
  get obfuscatedData() {
    return stringify([obfuscateData(JSON.parse(this.args.value))], {});
  }

  @action
  onSetup(editor) {
    // store reference to codemirror editor so that it can be passed to valueUpdated when restoring example
    this._codemirrorEditor = editor;
  }

  @action
  onUpdate(...args) {
    if (!this.args.readOnly) {
      // catching a situation in which the user is not readOnly and has not provided a valueUpdated function to the instance
      this.args.valueUpdated(...args);
    }
  }

  @action
  onFocus(...args) {
    if (this.args.onFocusOut) {
      this.args.onFocusOut(...args);
    }
  }

  @action
  restoreExample() {
    // set value to null which will cause the example value to be passed into the editor
    this.args.valueUpdated(null, this._codemirrorEditor);
  }
}
