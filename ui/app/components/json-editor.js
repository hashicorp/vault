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
 * @param {object} [options] - option object that overrides codemirror default options such as the styling.
 */

const JSON_EDITOR_DEFAULTS = {
  // IMPORTANT: `gutters` must come before `lint` since the presence of
  // `gutters` is cached internally when `lint` is toggled
  gutters: ['CodeMirror-lint-markers'],
  tabSize: 2,
  mode: 'application/json',
  lineNumbers: true,
  lint: { lintOnChange: false },
  theme: 'hashi',
  readOnly: false,
  showCursorWhenSelecting: true,
};

export default class JsonEditorComponent extends Component {
  value = null;
  valueUpdated = null;
  onFocusOut = null;
  readOnly = false;
  options = null;

  constructor() {
    super(...arguments);
    this.options = { ...JSON_EDITOR_DEFAULTS, ...this.args.options };
    if (this.options.autoHeight) {
      this.options.viewportMargin = Infinity;
      delete this.options.autoHeight;
    }
    if (this.options.readOnly) {
      this.options.readOnly = 'nocursor';
      this.options.lineNumbers = false;
      delete this.options.gutters;
    }
  }

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
