import Component from '@glimmer/component';
import { action } from '@ember/object';

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
 * @param {Function} valueUpdated - action to preform when you edit the codemirror value.
 * @param {Function} [onFocusOut] - action to preform when you focus out of codemirror.
 * @param {string} [helpText] - helper text.
 * @param {Object} [extraKeys] - to provide keyboard shortcut methods for things like saving on shit + enter.
 * @param {Array} [gutters] - An array of CSS class names or class name / CSS string pairs, each of which defines a width (and optionally a background), and which will be used to draw the background of the gutters.
 * @param {string} [mode] - right now we only import ruby so must be ruby or defaults to effectively null.
 * @param {Boolean} [readOnly] - defaults to false.
 * @param {String} [theme] - specify or customize the look.
 * @param {String} [value] - value within the display.
 * @param {String} [viewportMargin] - Sized of viewport. Often set to "Infinity" to show full amount always.
 */

export default class JsonEditorComponent extends Component {
  get getShowToolbar() {
    return this.args.showToolbar === false ? false : true;
  }

  @action
  update(...args) {
    if (!this.args.readOnly) {
      // catching a situation in which the user is not readOnly and has not provided a valueUpdated function to the instance
      this.args.valueUpdated(...args);
    }
  }

  @action
  focus(...args) {
    if (this.args.onFocusOut) {
      this.args.onFocusOut(...args);
    }
  }
}
