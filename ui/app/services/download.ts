import Service from '@ember/service';
import timestamp from 'core/utils/timestamp';

interface Extensions {
  csv: string;
  hcl: string;
  sentinel: string;
  json: string;
  pem: string;
  txt: string;
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
const EXTENSION_TO_MIME: Extensions = {
  csv: 'txt/csv',
  hcl: 'text/plain',
  sentinel: 'text/plain',
  json: 'application/json',
  pem: 'application/x-pem-file',
  txt: 'text/plain',
};

export default class DownloadService extends Service {
  download(filename: string, content: string, extension: string) {
    // replace spaces with hyphens, append extension to filename
    const formattedFilename =
      `${filename?.replace(/\s+/g, '-')}.${extension}` ||
      `vault-data-${timestamp.now().toISOString()}.${extension}`;

    // map extension to MIME type or use default
    const mimetype = EXTENSION_TO_MIME[extension as keyof Extensions] || 'text/plain';

    // commence download
    const { document, URL } = window;
    const downloadElement = document.createElement('a');
    const data = new File([content], formattedFilename, { type: mimetype });
    downloadElement.download = formattedFilename;
    downloadElement.href = URL.createObjectURL(data);
    document.body.appendChild(downloadElement);
    downloadElement.click();
    URL.revokeObjectURL(downloadElement.href);
    downloadElement.remove();
    return formattedFilename;
  }

  // SAMPLE CSV FORMAT ('content' argument)
  // Must be a string with each row \n separated and each column comma separated
  // 'Namespace path,Authentication method,Total clients,Entity clients,Non-entity clients\n
  //  namespacelonglonglong4/,,191,171,20\n
  //  namespacelonglonglong4/,auth/method/uMGBU,35,20,15\n'
  csv(filename: string, content: string) {
    this.download(filename, content, 'csv');
  }

  pem(filename: string, content: string) {
    this.download(filename, content, 'pem');
  }

  miscExtension(filename: string, content: string, extension: string) {
    this.download(filename, content, extension);
  }
}
