import Component from '@glimmer/component';
import layout from '../templates/components/download-csv';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';

/**
 * @module DownloadCsv
 * Download csv component is used to display a link which initiates a csv file download of the data provided by it's parent component.
 *
 * @example
 * ```js
 * <DownloadCsv @label={{'Export all namespace data'}} @csvData={{"Namespace path,Active clients /n nsTest5/,2725"}} @fileName={{'client-count.csv'}} />
 * ```
 *
 * @param {string} label - Label for the download link button
 * @param {string} csvData - Data in csv format
 * @param {string} fileName - Custom name for the downloaded file
 *
 */
class DownloadCsvComponent extends Component {
  @action
  downloadCsv() {
    let hiddenElement = document.createElement('a');
    hiddenElement.setAttribute('href', 'data:text/csv;charset=utf-8,' + encodeURI(this.args.csvData));
    hiddenElement.setAttribute('target', '_blank');
    hiddenElement.setAttribute('download', this.args.fileName || 'vault-data.csv');
    hiddenElement.click();
  }
}

export default setComponentTemplate(layout, DownloadCsvComponent);
