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
 * @param {object} [options] - option object that overrides codemirror default options such as the styling.
 */

export default class JsonEditorComponent extends Component {
  get getShowToolbar() {
    return this.args.showToolbar === false ? false : true;
  }

  @action
  update(value) {
    if (!this.args.readOnly) {
      // ARG TODO
      this.args.valueUpdated(value);
    }
  }

  @action
  focus(...args) {
    if (this.args.onFocusOut) {
      this.args.onFocusOut(...args);
    }
  }
}
