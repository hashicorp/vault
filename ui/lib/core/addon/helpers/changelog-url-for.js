import { helper } from '@ember/component/helper';

export function changelogUrlFor([version]) {
  // returns a url to the changelog for the specified version
  let url = 'http://www.github.com/hashicorp/vault/blob/master/CHANGELOG.md#';

  // strip the '+prem' from enterprise versions and remove periods
  let versionNumber = version
    .split('+')[0]
    .split('.')
    .join('');

  if (versionNumber === '142') {
    url = url.concat('142-may-21st-2020');
  } else if (versionNumber === '141') {
    url = url.concat('141-april-30th-2020');
  } else if (versionNumber === '140') {
    url = url.concat('140-april-7th-2020');
  }

  return url;
}

export default helper(changelogUrlFor);
