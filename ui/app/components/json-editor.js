/* eslint-disable no-undef */
import Component from '@ember/component';
import { computed } from '@ember/object';

// ARG TODO glimmerize this component and add documentation
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
  diffView: false,

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
  // Computed diff
  visualDiff: computed('leftSideVersionData', 'rightSideVersionData', function() {
    let diffpatcher = jsondiffpatch.create({});

    let delta = diffpatcher.diff(this.leftSideVersionData, this.rightSideVersionData);
    if (delta == undefined) {
      return 'No changes, previous state matches.';
    }
    // beautiful html diff
    return jsondiffpatch.formatters.html.format(delta, this.leftSideVersionData);
  }),
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
