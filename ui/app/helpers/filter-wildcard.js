import { helper as buildHelper } from '@ember/component/helper';

export function filterWildcard([string, array]) {
  let stringId = string.id;
  const filterBy = stringId =>
    array.filter(item => new RegExp('^' + stringId.replace(/\*/g, '.*') + '$').test(item));
  return filterBy(stringId).length;
}

export default buildHelper(filterWildcard);
