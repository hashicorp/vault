// https://www.npmjs.com/package/ember-modifier
// import { modifier } from 'ember-modifier';

// export default modifier((element, [eventName, handler]) => {
//   element.addEventListener(eventName, handler);

//   return () => {
//     element.removeEventListener(eventName, handler);
//   };
// });

import { action } from '@ember/object';
import { bind } from '@ember/runloop';
import codemirror from 'codemirror';
import Modifier from 'ember-modifier';

import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/mode/clike/clike';
import 'codemirror/mode/go/go';
import 'codemirror/mode/javascript/javascript';
import 'codemirror/mode/python/python';

const EXTENSION_REGEXP = /(?:\.([^.]+))?$/;

/**
 * Maps file extensions to loaded, CodeMirror-compatible language modes.
 *
 * **Important:** These CodeMirror modes must be loaded to be useable. See
 * the above imports which load the supported language modes.
 */
// const modeMap: Record<string, string> = {
const modeMap = {
  go: 'text/x-go',
  java: 'text/x-java',
  js: 'javascript',
  py: 'python',
};

/**
 * This is a magic CodeMirror mode string to indicate that no highlighting
 * should be used.
 *
 * See https://codemirror.net/doc/manual.html#option_mode.
 */
const DoNotHighlight = 'null';

// interface Args {
//   named: {
//     content: string;
//     path: string;
//     readOnly: boolean;
//     onUpdate: (content: string) => void;
//     [key: string]: unknown;
//   };
//   positional: never;
// }

export default class CodeMirrorModifier extends Modifier {
  didInstall() {
    this._setup();
  }

  didUpdateArguments() {
    if (this._editor.getValue() !== this.args.named.content) {
      this._editor.setValue(this.args.named.content);
    }

    this._editor.setOption('readOnly', this.args.named.readOnly);
    this._editor.setOption('mode', this.mode);
  }

  // private _editor!: CodeMirror.Editor;

  /**
   * Transforms the given path into an equivalent CodeMirror compatible
   * language mode string by inspecting the extension.
   *
   * If no matching language modes are supported or the file extension cannot be
   * determined, this will return the magic CodeMirror "null" string mode value.
   * The value "null" indicates no highlighting should be applied.
   */
  get mode() {
    if (!this.args.named.path) {
      return DoNotHighlight;
    }

    const extension = EXTENSION_REGEXP.exec(this.args.named.path);

    if (!extension || !extension[1]) {
      return DoNotHighlight;
    }

    return modeMap[extension[1].toLowerCase()] || DoNotHighlight;
  }

  @action
  _onChange(editor) {
    this.args.named.onUpdate(editor.getValue());
  }

  _setup() {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    // const JSON_EDITOR_DEFAULTS = {
    //   // IMPORTANT: `gutters` must come before `lint` since the presence of
    //   // `gutters` is cached internally when `lint` is toggled
    //   gutters: ['CodeMirror-lint-markers'],
    //   tabSize: 2,
    //   mode: 'application/json',
    //   lineNumbers: true,
    //   lint: { lintOnChange: false },
    //   theme: 'hashi',
    //   readOnly: false,
    //   showCursorWhenSelecting: true,
    // };

    const editor = codemirror(this.element, {
      lineNumbers: false,
      matchBrackets: true,
      mode: this.mode,
      readOnly: this.args.named.readOnly,
      styleActiveLine: true,
      theme: 'hashi',
      value: this.args.named.content || '',
      viewportMargin: Infinity,
    });

    editor.on('change', bind(this, this._onChange));

    this._editor = editor;
  }
}
