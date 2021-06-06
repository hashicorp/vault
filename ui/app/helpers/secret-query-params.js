import { helper } from '@ember/component/helper';

export function secretQueryParams([backendType]) {
  if (backendType === 'transit') {
    return { tab: 'actions' };
  }
  return;
}

export default helper(secretQueryParams);
