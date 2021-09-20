import Component from '@glimmer/component';
import layout from '../templates/components/download-csv';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';

class DownloadCsvComponent extends Component {
  @action
  downloadCsv() {
    let hiddenElement = document.createElement('a');
    hiddenElement.setAttribute('href', 'data:text/csv;charset=utf-8,' + encodeURI(this.args.csvData));
    hiddenElement.setAttribute('target', '_blank');
    hiddenElement.setAttribute('download', this.args.fileName);
    hiddenElement.click();
  }
}

export default setComponentTemplate(layout, DownloadCsvComponent);
