import Component from '@ember/component';

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

export default Component.extend({
  showToolbar: true,
  title: null,
  subTitle: null,
  helpText: null,
  value: null,
  options: null,
  valueUpdated: null,
  onFocusOut: null,
  readOnly: false,

  init() {
    this._super(...arguments);
    this.options = { ...JSON_EDITOR_DEFAULTS, ...this.options };
    if (this.options.autoHeight) {
      this.options.viewportMargin = Infinity;
      delete this.options.autoHeight;
    }
    if (this.options.readOnly) {
      this.options.readOnly = 'nocursor';
      this.options.lineNumbers = false;
      delete this.options.gutters;
    }
  },

  actions: {
    updateValue(...args) {
      if (this.valueUpdated) {
        this.valueUpdated(...args);
      }
    },
    onFocus(...args) {
      if (this.onFocusOut) {
        this.onFocusOut(...args);
      }
    },
  },
});
