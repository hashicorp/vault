import { helper as buildHelper } from '@ember/component/helper';

export function isWildcardString([string]) {
  if (string && string.id && string.id.includes('*')) {
    // string contains a wildcard
    return true;
  }
}

export default buildHelper(isWildcardString);
