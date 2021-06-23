import Component from '@glimmer/component';
import { set } from '@ember/object';
import { action } from '@ember/object';
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
 * @param [inputOnly] {bool} - When true, only the file input will be rendered
 * @param [helpText] {string} - Text underneath label.
 * @param file {object} - * Object in the shape of:
 * {
 *   value: 'file contents here',
 *   fileName: 'nameOfFile.txt',
 *   enterAsText: boolean ability to enter as text
 * }
 * @param [onChange=Function.prototype] {Function|action} - A function to call when the value of the input changes.
 * @param [label=null] {string} - Text to use as the label for the file input. If null, a default will be rendered.
 */

export default class TextFile extends Component {
  fileHelpText = 'Select a file from your computer';
  textareaHelpText = 'Enter the value as text';
  elementId = guidFor(this);
  index = '';

  @tracked file = null;
  @tracked showValue = false;

  get inputOnly() {
    return this.args.inputOnly || false;
  }
  get label() {
    return this.args.label || null;
  }

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
    let file = this.args.file;
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
