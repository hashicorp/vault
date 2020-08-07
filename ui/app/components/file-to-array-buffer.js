import Component from '@ember/component';
import filesize from 'filesize';

/**
 * @module FileToArrayBuffer
 * `FileToArrayBuffer` is a component that will allow you to pick a file from the local file system. Once
 * loaded, this file will be emitted as a JS ArrayBuffer to the passed `onChange` callback.
 *
 * @example
 * ```js
 *   <FileToArrayBuffer @onChange={{action (mut file)}} />
 * ```
 * @param onChange=null {Function} - The function to call when the file read is complete. This function
 * recieves the file as a JS ArrayBuffer
 * @param [label=null {String}] - Text to use as the label for the file input
 * @param [fileHelpText=null {String} - Text to use as help under the file input
 *
 */
export default Component.extend({
  classNames: ['box', 'is-fullwidth', 'is-marginless', 'is-shadowless'],
  onChange: () => {},
  label: null,
  fileHelpText: null,

  file: null,
  fileName: null,
  fileSize: null,
  fileLastModified: null,

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.send('onChange', reader.result, file);
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
    onChange(fileAsBytes, fileMeta) {
      let { name, size, lastModifiedDate } = fileMeta || {};
      let fileSize = size ? filesize(size) : null;
      this.set('file', fileAsBytes);
      this.set('fileName', name);
      this.set('fileSize', fileSize);
      this.set('fileLastModified', lastModifiedDate);
      this.onChange(fileAsBytes, name);
    },
  },
});
