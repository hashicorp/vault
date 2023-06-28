import { helper } from '@ember/component/helper';

export default helper(function filterItemsFn([items, searchTerm], { attr = 'id' }) {
  if (!items) return [];
  if (!searchTerm) return items;
  return items.filter((item) => (item[attr] ? item[attr].includes(searchTerm) : false));
});
