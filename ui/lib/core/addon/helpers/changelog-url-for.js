import { helper } from '@ember/component/helper';

/*
This helper returns a url to the changelog for the specified version.
It assumes that Changelog headers for Vault versions >= 1.4.3 are structured as:

## v1.5.0
### Month, DD, YYYY

## v1.4.5
### Month, DD, YYY

etc.
*/

export function changelogUrlFor([version]) {
  let url = 'http://www.github.com/hashicorp/vault/blob/master/CHANGELOG.md#';

  // strip the '+prem' from enterprise versions and remove periods
  let versionNumber = version
    .split('+')[0]
    .split('.')
    .join('');

  // only recent versions have a predictable url
  if (versionNumber >= '143') {
    return url.concat('v', versionNumber);
  } else {
    return url;
  }
}

export default helper(changelogUrlFor);
