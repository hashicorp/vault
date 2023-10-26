/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { helper } from '@ember/component/helper';

/*
This helper returns a url to the changelog for the specified version.
It assumes that Changelog headers for Vault versions >= 1.4.3 are structured as:

## 1.5.0
### Month, DD, yyyy

## 1.4.5
### Month, DD, YYY

etc.
*/

export function changelogUrlFor([version]) {
  const url = 'https://www.github.com/hashicorp/vault/blob/main/CHANGELOG.md#';
  if (!version) return url;
  try {
    // strip the '+prem' from enterprise versions and remove periods
    const versionNumber = version.split('+')[0].split('.').join('');

    // only recent versions have a predictable url
    if (versionNumber >= 143) {
      return url.concat(versionNumber);
    }
  } catch (e) {
    console.log(e); // eslint-disable-line
    console.log('Cannot generate URL for version: ', version); // eslint-disable-line
  }
  return url;
}

export default helper(changelogUrlFor);
