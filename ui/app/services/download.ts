import Service from '@ember/service';

export default class DownloadService extends Service {
  // some browsers (ex. Firefox) require the filename has an explicit extension, always include it in the filename

  download(filename: string, mimetype: string, content: string) {
    const { document, URL } = window;
    const downloadElement = document.createElement('a');
    const data = new File([content], filename, { type: mimetype });
    downloadElement.download = filename;
    downloadElement.href = URL.createObjectURL(data);
    document.body.appendChild(downloadElement);
    downloadElement.click();
    URL.revokeObjectURL(downloadElement.href);
    downloadElement.remove();
  }

  // SAMPLE CSV FORMAT ('content' argument)
  // Must be a string with each row \n separated and each column comma separated
  // 'Namespace path,Authentication method,Total clients,Entity clients,Non-entity clients\n
  //  namespacelonglonglong4/,,191,171,20\n
  //  namespacelonglonglong4/,auth/method/uMGBU,35,20,15\n'
  csv(filename: string, content: string) {
    const formattedFilename = `${filename?.replace(/\s+/g, '-')}.csv` || 'vault-data.csv';
    this.download(formattedFilename, 'text/csv', content);
    return formattedFilename;
  }

  pem(filename: string, content: string) {
    const formattedFilename = `${filename?.replace(/\s+/g, '-')}.pem` || 'vault-cert.pem';
    this.download(formattedFilename, 'application/x-pem-file', content);
    return formattedFilename;
  }
}
