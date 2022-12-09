import { action } from '@ember/object';
import Component from '@glimmer/component';

export default class DownloadButton extends Component {
  get extension() {
    return this.args.extension || 'txt';
  }

  get mime() {
    return this.args.mime || 'text/plain';
  }

  get filename() {
    return `${this.args.filename}-${new Date().toISOString()}.${this.extension}`;
  }

  get content() {
    let data = this.args.data;
    if (this.args.stringify) {
      data = JSON.stringify(data, null, 2);
    }
    const content = new File([data], this.filename, { type: this.mime });
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
