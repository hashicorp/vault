/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import timestamp from 'vault/utils/timestamp';
import { tracked } from '@glimmer/tracking';
import { assert } from '@ember/debug';
/**
 * @module DownloadButton
 * DownloadButton components are an action button used to download data. Both the action text and icon are yielded.
 * * NOTE: when using in an engine, remember to add the 'download' service to its dependencies (in /engine.js) and map to it in /app.js
 * [ember-docs](https://ember-engines.com/docs/services)
 * @example
 * ```js
 *   <DownloadButton
 *     class="button"
 *     @data={{this.data}}
 *     @filename={{this.filename}}
 *     @mime={{this.mime}}
 *     @extension={{this.extension}}
 *     @stringify={{true}}
 *   >
 *    <Icon @name="download" />
 *      Download
 *   </DownloadButton>
 * ```
 * @param {string} [filename] - name of file that prefixes the ISO timestamp generated at download
 * @param {string} [data] - data to download
 * @param {function} [fetchData] - function that fetches data and returns download content
 * @param {string} [extension='txt'] - file extension, the download service uses this to determine the mimetype
 * @param {boolean} [stringify=false] - argument to stringify the data before passing to the File constructor
 * @param {callback} [onSuccess] - callback from parent to invoke if download is successful
 */

export default class DownloadButton extends Component {
  @service download;
  @service flashMessages;
  @tracked fetchedData;

  constructor() {
    super(...arguments);
    const hasConflictingArgs = this.args.data && this.args.fetchData;
    assert(
      'Only pass either @data or @fetchData, passing both means @data will be overwritten by the return value of @fetchData',
      !hasConflictingArgs
    );
  }
  get filename() {
    const ts = timestamp.now().toISOString();
    return this.args.filename ? this.args.filename + '-' + ts : ts;
  }

  get content() {
    if (this.args.stringify) {
      return JSON.stringify(this.args.data, null, 2);
    }
    return this.fetchedData || this.args.data;
  }

  get extension() {
    return this.args.extension || 'txt';
  }

  @action
  async handleDownload() {
    if (this.args.fetchData) {
      this.fetchedData = await this.args.fetchData();
    }
    try {
      this.download.miscExtension(this.filename, this.content, this.extension);
      this.flashMessages.info(`Downloading ${this.filename}`);
      if (this.args.onSuccess) {
        this.args.onSuccess();
      }
    } catch (error) {
      this.flashMessages.danger(errorMessage(error, 'There was a problem downloading. Please try again.'));
    }
  }
}
