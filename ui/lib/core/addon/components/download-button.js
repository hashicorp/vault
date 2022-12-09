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
