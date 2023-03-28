/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { set } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

const BASE_64_REGEX = /^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/gi;

export default Component.extend({
  'data-test-pgp-file': true,
  attributeBindings: ['data-test-pgp-file'],
  classNames: ['box', 'is-fullwidth', 'is-marginless', 'is-shadowless'],
  key: null,
  index: null,
  onChange: () => {},

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

  /*
   * @public
   * @param String
   * Text to use as help under the textarea in text-input mode
   * If null, a default will be rendered
   */
  textareaHelpText: null,

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => this.setPGPKey.perform(reader.result, file.name);
    // this gives us a base64-encoded string which is important in the onload
    reader.readAsDataURL(file);
  },

  setPGPKey: task(
    waitFor(function* (dataURL, filename) {
      const b64File = dataURL.split(',')[1].trim();
      const decoded = atob(b64File).trim();

      // If a b64-encoded file was uploaded, then after decoding, it
      // will still be b64.
      // If after decoding it's not b64, we want
      // the original as it was only encoded when we used `readAsDataURL`.
      const fileData = decoded.match(BASE_64_REGEX) ? decoded : b64File;
      yield this.onChange(this.index, { value: fileData, filename: filename });
    })
  ),

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
      const key = this.key;
      set(key, 'value', e.target.value);
      this.onChange(this.index, this.key);
    },
    clearKey() {
      this.onChange(this.index, { value: '' });
    },
  },
});
