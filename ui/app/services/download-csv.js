import Service from '@ember/service';

// SAMPLE CSV FORMAT ('content' argument)
// Must be a string with each row \n separated and each column comma separated
// 'Namespace path,Authentication method,Total clients,Entity clients,Non-entity clients\n
//  namespacelonglonglong4/,,191,171,20\n
//  namespacelonglonglong4/,auth/method/uMGBU,35,20,15\n'

export default class DownloadCsvService extends Service {
  download(filename, content) {
    let formattedFilename = filename?.replace(/\s+/g, '-') || 'vault-data.csv';
    let { document, URL } = window;
    let downloadLink = document.createElement('a');
    downloadLink.download = formattedFilename;
    downloadLink.href = URL.createObjectURL(
      new Blob([content], {
        type: 'text/csv',
      })
    );

    document.body.appendChild(downloadLink);
    downloadLink.click();
    downloadLink.remove();
  }
}
