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
//  ARG TODO Fill in
 */

export default class JsonEditorComponent extends Component {
  value = null;
  valueUpdated = null;
  onFocusOut = null;

  get getShowToolbar() {
    return this.args.showToolbar === false ? false : true;
  }

  @action
  updateValue(...args) {
    if (this.args.valueUpdated) {
      this.args.valueUpdated(...args);
    }
  }

  @action
  onFocus(...args) {
    if (this.args.onFocusOut) {
      this.args.onFocusOut(...args);
    }
  }
}
