import Service from '@ember/service';

// SAMPLE CSV FORMAT ('content' argument)
// Must be a string with each row \n separated and each column comma separated
// 'Namespace path,Authentication method,Total clients,Entity clients,Non-entity clients\n
//  namespacelonglonglong4/,,191,171,20\n
//  namespacelonglonglong4/,auth/method/uMGBU,35,20,15\n'

export default class DownloadCsvService extends Service {
  download(filename, content) {
    // even though Blob type 'text/csv' is specified below, some browsers (ex. Firefox) require the filename has an explicit extension
    let formattedFilename = `${filename?.replace(/\s+/g, '-')}.csv` || 'vault-data.csv';
    let { document, URL } = window;
    let downloadElement = document.createElement('a');
    downloadElement.download = formattedFilename;
    downloadElement.href = URL.createObjectURL(
      new Blob([content], {
        type: 'text/csv',
      })
    );
    document.body.appendChild(downloadElement);
    downloadElement.click();
    URL.revokeObjectURL(downloadElement.href);
    downloadElement.remove();
  }
}
