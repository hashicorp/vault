import { assign } from '@ember/polyfills';
import IvyCodemirrorComponent from './ivy-codemirror';
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

export default IvyCodemirrorComponent.extend({
  'data-test-component': 'json-editor',
  updateCodeMirrorOptions() {
    const options = assign({}, JSON_EDITOR_DEFAULTS, this.get('options'));
    if (options.autoHeight) {
      options.viewportMargin = Infinity;
      delete options.autoHeight;
    }

    if (options) {
      Object.keys(options).forEach(function(option) {
        this.updateCodeMirrorOption(option, options[option]);
      }, this);
    }
  },
});
