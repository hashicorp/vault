import Component from '@glimmer/component';
import { set } from '@ember/object';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class TextFile extends Component {
  'data-test-component' = 'text-file';
  classNames = ['box', 'is-fullwidth', 'is-marginless', 'is-shadowless'];
  classNameBindings = ['inputOnly:is-paddingless'];

  /*
   * @public
   * @param Object
   * Object in the shape of:
   * {
   *   value: 'file contents here',
   *   fileName: 'nameOfFile.txt',
   *   enterAsText: bool
   * }
   */
  @tracked
  file = null;
  @tracked
  index = null;
  @tracked
  showValue = false;

  /*
   * @public
   * @param Boolean
   * When true, only the file input will be rendered
   */
  inputOnly = false;

  /*
   * @public
   * @param String
   * Text to use as the label for the file input
   * If null, a default will be rendered
   */
  label = null;

  /*
   * @public
   * @param String
   * Text to use as help under the file input
   * If null, a default will be rendered
   */
  fileHelpText = 'Select a file from your computer';

  /*
   * @public
   * @param String
   * Text to use as help under the textarea in text-input mode
   * If null, a default will be rendered
   */
  textareaHelpText = 'Enter the value as text';

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.setFile(reader.result, file.name);
    reader.readAsText(file);
  }

  setFile(contents, filename) {
    let index = this.args.index || this.index; // ARG Todo understand defaults and args.
    console.log(index, 'index');
    console.log(Object.keys(this.args));
    this.args.onChange(index, { value: contents, fileName: filename });
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
    console.log(file, 'file');
    this.args.onChange(this.args.index, file);
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
