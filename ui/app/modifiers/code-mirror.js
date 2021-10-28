import { getOwner } from '@ember/application';
import { guidFor } from '@ember/object/internals';
import Modifier from 'ember-modifier';
import codemirror from 'codemirror';
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/selection/active-line';
import 'codemirror/mode/javascript/javascript';
import 'codemirror/mode/ruby/ruby';
// import '@hashicorp/sentinel-codemirror/sentinel';
import 'codemirror/keymap/sublime';
import 'codemirror/addon/search/search';
import 'codemirror/addon/search/searchcursor';
import 'codemirror/addon/dialog/dialog';

export default class CodeMirrorModifier extends Modifier {
  get cmService() {
    return getOwner(this).lookup('service:code-mirror');
  }

  didInstall() {
    this._setup();
  }

  willRemove() {
    this._cleanup();
  }

  didUpdateArguments() {
    if (this._editor.getValue() !== this.args.named.value) {
      this._editor.setValue(this.args.named.value);
    }
  }

  _onChange(editor) {
    if (this.args.named.valueUpdated) {
      this.args.named.valueUpdated(editor.getValue());
    }
  }

  _setup() {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    // Assign an ID to this element if there is none. This is to
    // ensure that there are unique IDs in the code-mirror service
    // registry.
    if (!this.element.id) {
      this.element.id = guidFor(this.element);
    }

    let editor = codemirror(
      this.element,
      Object.assign(
        {
          value: this.args.named.value ? this.args.named.value : '',
          inputStyle: 'contenteditable',
        },
        this.args.named.options
      )
    );

    editor.on('change', editor => {
      this._onChange(editor);
    });

    if (this.cmService) {
      this.cmService.registerInstance(this.element.id, editor);
    }

    this._editor = editor;
  }

  _cleanup() {
    if (this.cmService) {
      this.cmService.unregisterInstance(this.element.id);
    }
  }
}
