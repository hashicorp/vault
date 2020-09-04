import { helper as buildHelper } from '@ember/component/helper';

export function isWildcardString([string]) {
  if (!string) {
    return false;
  }
  // string is actually an object which is what comes in in the searchSelect component
  if (typeof string === 'object') {
    if (!string.id && !string.name) {
      return false;
    }
    if (string.id) {
      return string.id.includes('*');
    }
    if (string.name) {
      return string.name.includes('*');
    }
  }
  // otherwise string is a string
  if (string.includes('*')) {
    return true;
  }
  // default to false
  return false;
}

export default buildHelper(isWildcardString);
