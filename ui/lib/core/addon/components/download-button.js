import { action } from '@ember/object';
import Component from '@glimmer/component';

export default class DownloadButton extends Component {
  get extension() {
    return this.args.extension || 'txt';
  }

  get mime() {
    return this.args.mime || 'text/plain';
  }

  get download() {
    return `${this.args.filename}-${new Date().toISOString()}.${this.extension}`;
  }

  get fileLike() {
    let file;
    let data = this.args.data;
    const filename = this.download;
    const mime = this.mime;
    if (this.args.stringify) {
      data = JSON.stringify(data, null, 2);
    }
    if (window.navigator.msSaveOrOpenBlob) {
      file = new Blob([data], { type: mime });
      file.name = filename;
    } else {
      file = new File([data], filename, { type: mime });
    }
    return file;
  }

  get href() {
    return window.URL.createObjectURL(this.fileLike);
  }

  @action
  handleDownload(event) {
    if (!window.navigator.msSaveOrOpenBlob) {
      return;
    }
    event.preventDefault();
    const file = this.fileLike;
    window.navigator.msSaveOrOpenBlob(file, file.name);
  }
}
