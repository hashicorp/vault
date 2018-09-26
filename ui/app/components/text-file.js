import Component from '@ember/component';
import { set } from '@ember/object';

export default Component.extend({
  'data-test-component': 'text-file',
  classNames: ['box', 'is-fullwidth', 'is-marginless', 'is-shadowless'],
  classNameBindings: ['inputOnly:is-paddingless'],

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
  file: null,

  index: null,
  onChange: () => {},

  /*
   * @public
   * @param Boolean
   * When true, only the file input will be rendered
   */
  inputOnly: false,

  /*
   * @public
   * @param String
   * Text to use as the label for the file input
   * If null, a default will be rendered
   */
  label: null,

  /*
   * @public
   * @param String
   * Text to use as help under the file input
   * If null, a default will be rendered
   */
  fileHelpText: 'Select a file from your computer',

  /*
   * @public
   * @param String
   * Text to use as help under the textarea in text-input mode
   * If null, a default will be rendered
   */
  textareaHelpText: 'Enter the value as text',

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.setFile(reader.result, file.name);
    reader.readAsText(file);
  },

  setFile(contents, filename) {
    this.get('onChange')(this.get('index'), { value: contents, fileName: filename });
  },

  actions: {
    pickedFile(e) {
      const { files } = e.target;
      if (!files.length) {
        return;
      }
      for (let i = 0, len = files.length; i < len; i++) {
        this.readFile(files[i]);
      }
    },
    updateData(e) {
      const file = this.get('file');
      set(file, 'value', e.target.value);
      this.get('onChange')(this.get('index'), this.get('file'));
    },
    clearFile() {
      this.get('onChange')(this.get('index'), { value: '' });
    },
  },
});
