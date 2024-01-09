/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { guidFor } from '@ember/object/internals';
/**
 * @module TextFile
 * `TextFile` components render a file upload input with the option to toggle a "Enter as text" button
 *  that changes the input into a textarea
 *
 * @example
 * <TextFile
 *  @uploadOnly={{true}}
 *  @helpText="help text"
 *  @onChange={{this.handleChange}}
 *  @label="PEM Bundle"
 * />
 *
 * @param {function} onChange - Callback function to call when the value of the input changes, returns an object in the shape of { value: fileContents, filename: 'some-file.txt' }
 * @param {bool} [uploadOnly=false] - When true, renders a static file upload input and removes the option to toggle and input plain text
 * @param {string} [helpText] - Text underneath label.
 * @param {string} [label='File']  - Text to use as the label for the file input. If none, default of 'File' is rendered
 */

export default class TextFileComponent extends Component {
  @tracked content = '';
  @tracked filename = '';
  @tracked uploadError = '';
  @tracked showTextArea = false;
  elementId = guidFor(this);

  async readFile(file) {
    try {
      this.content = await file.text();
      this.filename = file.name;
      this.handleChange();
    } catch (error) {
      this.clearFile();
      this.uploadError = 'There was a problem uploading. Please try again.';
    }
  }

  @action
  handleFileUpload(e) {
    e.preventDefault();
    const { files } = e.target;
    if (!files.length) return;
    this.readFile(files[0]);
  }

  @action
  handleTextInput(e) {
    e.preventDefault();
    this.content = e.target.value;
    this.handleChange();
  }

  @action
  clearFile() {
    this.content = '';
    this.filename = '';
    this.handleChange();
  }

  handleChange() {
    this.args.onChange({ value: this.content, filename: this.filename });
    this.uploadError = '';
  }
}
