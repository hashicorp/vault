/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import timestamp from 'core/utils/timestamp';

interface Extensions {
  csv: string;
  hcl: string;
  sentinel: string;
  json: string;
  jsonl: string;
  pem: string;
  txt: string;
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
const EXTENSION_TO_MIME: Extensions = {
  csv: 'txt/csv',
  hcl: 'text/plain',
  sentinel: 'text/plain',
  json: 'application/json',
  jsonl: 'application/json',
  pem: 'application/x-pem-file',
  txt: 'text/plain',
};

export default class DownloadService extends Service {
  formatFilename(filename: string, extension: keyof Extensions = 'txt') {
    const downloadTimestamp = timestamp.now().toISOString();
    // replace spaces with underscores
    const name = filename ? filename?.replace(/\s+/g, '_') : 'vault_data';
    // appends extension to filename
    return `${name}_${downloadTimestamp}.${extension}`;
  }

  download(filename: string, content: any, extension: keyof Extensions) {
    const formattedFilename = this.formatFilename(filename, extension);

    // map extension to MIME type or use default
    const mimetype = EXTENSION_TO_MIME[extension] || 'text/plain';

    // commence download
    const downloadElement = document.createElement('a');
    const data = new File([content], formattedFilename, { type: mimetype });
    downloadElement.download = formattedFilename;
    downloadElement.href = URL.createObjectURL(data);
    downloadElement.click();
    URL.revokeObjectURL(downloadElement.href);
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

  miscExtension(filename: string, content: string, extension: keyof Extensions) {
    this.download(filename, content, extension);
  }
}
