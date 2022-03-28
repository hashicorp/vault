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
 * @param {Function} [valueUpdated] - action to preform when you edit the codemirror value.
 * @param {Function} [onFocusOut] - action to preform when you focus out of codemirror.
 * @param {string} [helpText] - helper text.
 * @param {Object} [extraKeys] - to provide keyboard shortcut methods for things like saving on shift + enter.
 * @param {Array} [gutters] - An array of CSS class names or class name / CSS string pairs, each of which defines a width (and optionally a background), and which will be used to draw the background of the gutters.
 * @param {string} [mode] - right now we only import ruby so must mode but be ruby or defaults to javascript. If you wanted another language you need to import it into the modifier.
 * @param {Boolean} [readOnly] - defaults to false.
 * @param {String} [theme] - specify or customize the look via css.
 * @param {String} [value] - value within the display.
 * @param {String} [viewportMargin] - Size of viewport. Often set to "Infinity" to load/show all text regardless of length.
 */

export default class JsonEditorComponent extends Component {
  get getShowToolbar() {
    return this.args.showToolbar === false ? false : true;
  }

  @action
  onUpdate(...args) {
    if (!this.args.readOnly) {
      // catching a situation in which the user is not readOnly and has not provided a valueUpdated function to the instance
      this.args.valueUpdated(...args);
    }
  }

  @action
  OnFocus(...args) {
    if (this.args.onFocusOut) {
      this.args.onFocusOut(...args);
    }
  }
}
