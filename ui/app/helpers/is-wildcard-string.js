import { helper as buildHelper } from '@ember/component/helper';

export function isWildcardString([string]) {
  if (!string) {
    return false;
  }
  // string is actually an object which is what comes in in the searchSelect component
  if (typeof string === 'object') {
    string = Object.values(string)[0];
  }

  return string.includes('*');
}

export default buildHelper(isWildcardString);
