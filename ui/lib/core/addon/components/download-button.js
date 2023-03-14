import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import timestamp from 'vault/utils/timestamp';
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
 * @param {string} [extension='txt'] - file extension, the download service uses this to determine the mimetype
 * @param {boolean} [stringify=false] - argument to stringify the data before passing to the File constructor
 */

export default class DownloadButton extends Component {
  @service download;
  @service flashMessages;

  get filename() {
    const ts = timestamp.now().toISOString();
    return this.args.filename ? this.args.filename + '-' + ts : ts;
  }

  get content() {
    if (this.args.stringify) {
      return JSON.stringify(this.args.data, null, 2);
    }
    return this.args.data;
  }

  get extension() {
    return this.args.extension || 'txt';
  }

  @action
  handleDownload() {
    try {
      this.download.miscExtension(this.filename, this.content, this.extension);
      this.flashMessages.info(`Downloading ${this.filename}`);
    } catch (error) {
      this.flashMessages.danger(errorMessage(error, 'There was a problem downloading. Please try again.'));
    }
  }
}
