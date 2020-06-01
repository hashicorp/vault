import { helper } from '@ember/component/helper';

export function changelogUrlFor([version]) {
  // returns a url to the changelog for the specified version
  let url = 'http://www.github.com/hashicorp/vault/blob/master/CHANGELOG.md#';

  // strip the '+prem' from enterprise versions and remove periods
  let versionNumber = version
    .split('+')[0]
    .split('.')
    .join('');

  // only the most recent versions have a predictable url
  if (versionNumber >= '140') {
    url = url.concat(versionNumber);
  }

  return url;
}

export default helper(changelogUrlFor);
