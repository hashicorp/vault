import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
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
 * @param {string} data - data to download
 * @param {boolean} [stringify=false] - argument to stringify the data before passing to the File constructor
 * @param {string} [filename] - name of file that prefixes the ISO timestamp generated when download
 * @param {string} [mime='text/plain'] - media type to be downloaded
 * @param {string} [extension='txt'] - file extension
 */

export default class DownloadButton extends Component {
  @service download;

  get extension() {
    return this.args.extension || 'txt';
  }

  get mime() {
    return this.args.mime || 'text/plain';
  }

  get filename() {
    const defaultFilename = `${new Date().toISOString()}.${this.extension}`;
    return this.args.filename ? this.args.filename + '-' + defaultFilename : defaultFilename;
  }

  get data() {
    if (this.args.stringify) {
      return JSON.stringify(this.args.data, null, 2);
    }
    return this.args.data;
  }

  @action
  handleDownload() {
    this.download.download(this.filename, this.mime, this.data);
  }
}
