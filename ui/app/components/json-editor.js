/* eslint-disable no-undef */
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

// ARG TODO fill out
/**
 * @module JsonEditor
 *
 * @example
 * ```js
 * <JsonEditor @activeClusterName={{cluster.name}} @onLinkClick={{action "onLinkClick"}} />
 * ```
 *
 * @param {string} activeClusterName - name of the current cluster, passed from the parent.
 * @param {Function} onLinkClick - parent action which determines the behavior on link click
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

  get visualDiff() {
    let diffpatcher = jsondiffpatch.create({});
    let delta = diffpatcher.diff(this.args.leftSideVersionData, this.args.rightSideVersionData);
    if (delta == undefined) {
      return 'No changes, previous state matches.';
    }
    return jsondiffpatch.formatters.html.format(delta, this.args.leftSideVersionData);
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
