import { helper as buildHelper } from '@ember/component/helper';

export function includes([haystack, needle]) {
  return haystack.includes(needle);
}

export default buildHelper(includes);
