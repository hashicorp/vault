import { action } from '@ember/object';
import { bind } from '@ember/runloop';
import codemirror from 'codemirror';
import Modifier from 'ember-modifier';

import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/addon/lint/lint.js';
import 'codemirror/addon/lint/json-lint.js';
// right now we only use the ruby and javascript, if you use another mode you'll need to import it.
// https://codemirror.net/mode/
import 'codemirror/mode/ruby/ruby';
import 'codemirror/mode/javascript/javascript';

export default class CodeMirrorModifier extends Modifier {
  didInstall() {
    this._setup();
  }

  didUpdateArguments() {
    this._editor.setOption('readOnly', this.args.named.readOnly);
    if (!this.args.named.content) {
      return;
    }
    if (this._editor.getValue() !== this.args.named.content) {
      this._editor.setValue(this.args.named.content);
    }
  }

  @action
  _onChange(editor) {
    // avoid sending change event after initial setup when editor value is set to content
    if (this.args.named.content !== editor.getValue()) {
      this.args.named.onUpdate(editor.getValue(), this._editor);
    }
  }

  @action
  _onFocus(editor) {
    this.args.named.onFocus(editor.getValue());
  }

  _setup() {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }
    const editor = codemirror(this.element, {
      // IMPORTANT: `gutters` must come before `lint` since the presence of
      // `gutters` is cached internally when `lint` is toggled
      gutters: this.args.named.gutters || ['CodeMirror-lint-markers'],
      matchBrackets: true,
      lint: { lintOnChange: true },
      showCursorWhenSelecting: true,
      styleActiveLine: true,
      tabSize: 2,
      // all values we can pass into the JsonEditor
      extraKeys: this.args.named.extraKeys || '',
      lineNumbers: this.args.named.lineNumbers,
      mode: this.args.named.mode || 'application/json',
      readOnly: this.args.named.readOnly || false,
      theme: this.args.named.theme || 'hashi',
      value: this.args.named.content || '',
      viewportMargin: this.args.named.viewportMargin || '',
    });

    editor.on('change', bind(this, this._onChange));
    editor.on('focus', bind(this, this._onFocus));

    this._editor = editor;
  }
}
