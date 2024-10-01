/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
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
 * @param {function} onChange - The function to call when the file read is complete. Receives the file as a JS ArrayBuffer
 * @param {string} [label] - Text to use as the label for the file input
 * @param {string} [fileHelpText] - Text to use as help under the file input
 *
 */
export default class FileToArrayBufferComponent extends Component {
  @tracked filename = null;
  @tracked fileSize = null;
  @tracked fileLastModified = null;

  readFile(file) {
    const reader = new FileReader();
    reader.onload = () => {
      // raft-snapshot-restore test was failing on CI trying to invoke fileChange on destroyed object
      // ensure that the component has not been torn down
      if (!this.isDestroyed && !this.isDestroying) {
        this.fileChange(reader.result, file);
      }
    };
    reader.readAsArrayBuffer(file);
  }

  @action
  pickedFile(e) {
    const { files } = e.target;
    if (!files.length) {
      return;
    }
    for (let i = 0, len = files.length; i < len; i++) {
      this.readFile(files[i]);
    }
  }

  @action
  fileChange(fileAsBytes, fileMeta) {
    const { name, size, lastModifiedDate } = fileMeta || {};
    this.filename = name;
    this.fileSize = size ? filesize(size) : null;
    this.fileLastModified = lastModifiedDate;
    this.args.onChange(fileAsBytes, name);
  }
}
