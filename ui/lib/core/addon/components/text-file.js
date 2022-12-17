import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module TextFile
 * `TextFile` components render a file upload input with the option to toggle a "Enter as text" button
 *  that changes the input into a textarea
 *
 * @example
 * <TextFile
 *  @hideTextAreaToggle={{true}}
 *  @helpText="help text"
 *  @onChange={{this.handleChange}}
 *  @label="PEM Bundle"
 * />
 *
 * @param {function} onChange - Callback function to call when the value of the input changes, returns an object in the shape of { value: fileContents, filename: 'some-file.txt' }
 * @param {bool} [hideTextAreaToggle=false] - When true, renders a static file upload input and removes the option to toggle and input plain text
 * @param {string} [helpText] - Text underneath label.
 * @param {string} [label=null]  - Text to use as the label for the file input. If none, default of 'File' is rendered
 */

export default class TextFileComponent extends Component {
  @tracked file = null;
  @tracked showValue = false;
  @tracked showTextInput = false;

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.handleChange(reader.result, file.name);
    reader.readAsText(file);
  }

  handleChange(contents, filename) {
    this.args.onChange({ value: contents, filename });
  }

  @action
  handleFileUpload(e) {
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
  handleTextInput(e) {
    e.preventDefault();
    const file = this.args.file;
    file.value = e.target.value;
    this.args.onChange(file);
  }

  @action
  clearFile() {
    this.args.onChange({ value: '' });
  }

  @action
  toggleMask() {
    this.showValue = !this.showValue;
  }
}
