import Component from '@glimmer/component';
import { set, action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { guidFor } from '@ember/object/internals';

/**
 * @module TextFile
 * `TextFile` components are file upload components where you can either toggle to upload a file or enter text.
 *
 * @example
 * <TextFile
 *  @inputOnly={{true}}
 *  @helpText="help text"
 *  @file={{object}}
 *  @onChange={{action "someOnChangeFunction"}}
 *  @label={{"string"}}
 * />
 *
 * @param {object} file - Object in the shape of:
 * {
 *   value: 'file contents here',
 *   fileName: 'nameOfFile.txt',
 *   enterAsText: boolean ability to enter as text
 * }
 * @param {Function|action} [onChange=Function.prototype] - A function to call when the value of the input changes.
 * @param {bool} [inputOnly] - When true, only the file input will be rendered
 * @param {string} [helpText] - Text underneath label.
 * @param {string} [label=null]  - Text to use as the label for the file input. If none, a default of 'File' is rendered
 */

export default class TextFileComponent extends Component {
  fileHelpText = '';
  textareaHelpText = '';
  elementId = guidFor(this);
  index = '';

  @tracked file = null;
  @tracked showValue = false;

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.setFile(reader.result, file.name);
    reader.readAsText(file);
  }

  setFile(contents, filename) {
    this.args.onChange(this.index, { value: contents, fileName: filename });
  }

  @action
  pickedFile(e) {
    e.preventDefault();
    const { files } = e.target;
    if (!files.length) {
      return;
    }
    for (let i = 0, len = files.length; i < len; i++) {
      this.readFile(files[i]);
    }
  }
  @action
  updateData(e) {
    e.preventDefault();
    const file = this.args.file;
    set(file, 'value', e.target.value);
    this.args.onChange(this.index, file);
  }
  @action
  clearFile() {
    this.args.onChange(this.index, { value: '' });
  }
  @action
  toggleMask() {
    this.showValue = !this.showValue;
  }
}
