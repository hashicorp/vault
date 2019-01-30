import { helper as buildHelper } from '@ember/component/helper';

export const TABS = {
  entity: ['details', 'aliases', 'policies', 'groups', 'metadata'],
  'entity-alias': ['details', 'metadata'],
  //group will be used in the model hook of the route
  group: ['details', 'aliases', 'policies', 'members', 'parent-groups', 'metadata'],
  'group-internal': ['details', 'policies', 'members', 'parent-groups', 'metadata'],
  'group-external': ['details', 'aliases', 'policies', 'members', 'parent-groups', 'metadata'],
  'group-alias': ['details'],
};

export function tabsForIdentityShow([modelType, groupType]) {
  let key = modelType;
  if (groupType) {
    key = `${key}-${groupType}`;
  }
  return TABS[key];
}

export default buildHelper(tabsForIdentityShow);
