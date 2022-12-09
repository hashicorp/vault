import { action } from '@ember/object';
import Component from '@glimmer/component';
/**
 * @module DownloadButton
 * DownloadButton components are an action button used to download data. Both the action text and icon are yielded.
 *
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

  get content() {
    const content = new File([this.data], this.filename, { type: this.mime });
    return content;
  }

  // TODO refactor and call service instead
  @action
  handleDownload() {
    const { document, URL } = window;
    const downloadElement = document.createElement('a');
    downloadElement.download = this.filename;
    downloadElement.href = URL.createObjectURL(this.content);
    document.body.appendChild(downloadElement);
    downloadElement.click();
    URL.revokeObjectURL(downloadElement.href);
    downloadElement.remove();
  }
}
