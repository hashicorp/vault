import Service from '@ember/service';

// Sample of CSV format: each row needs to be \n separated, and each column separated by a comma
// 'Namespace path,Authentication method,Total clients,Entity clients,Non-entity clients\n
//  namespacelonglonglong4/,,191,171,20\n
//  namespacelonglonglong4/,auth/method/uMGBU,35,20,15\n'

export default class DownloadCsvService extends Service {
  download(filename, content) {
    let { document, URL } = window;
    let anchor = document.createElement('a');
    anchor.download = filename;
    anchor.href = URL.createObjectURL(
      new Blob([content], {
        type: 'text/csv',
      })
    );

    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
  }
}
