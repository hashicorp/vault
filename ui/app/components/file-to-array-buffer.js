import Component from '@ember/component';

export default Component.extend({
  classNames: ['box', 'is-fullwidth', 'is-marginless', 'is-shadowless'],
  onChange: () => {},
  file: null,
  fileName: null,

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
  fileHelpText: null,

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.send('onChange', reader.result, file.name);
    reader.readAsArrayBuffer(file);
  },

  actions: {
    pickedFile(e) {
      let { files } = e.target;
      if (!files.length) {
        return;
      }
      for (let i = 0, len = files.length; i < len; i++) {
        this.readFile(files[i]);
      }
    },
    clearFile() {
      this.send('onChange');
    },
    onChange(file, filename) {
      this.set('file', file);
      this.set('fileName', filename);
      this.onChange(file, filename);
    },
  },
});
