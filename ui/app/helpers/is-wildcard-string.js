import { helper as buildHelper } from '@ember/component/helper';

export function isWildcardString([string]) {
  if (!string) {
    return false;
  }
  // string is actually an object which is what comes in in the searchSelect component
  if (typeof string === 'object') {
    // if the dropdown is used the object is a class which cannot be converted into a string
    if (Object.prototype.hasOwnProperty.call(string, 'store')) {
      return false;
    }
    string = Object.values(string).toString();
  }

  return string.includes('*');
}

export default buildHelper(isWildcardString);
