import { helper } from '@ember/component/helper';

export function secretQueryParams([backendType, type = '']) {
  if (backendType === 'transit') {
    return { tab: 'actions' };
  }
  if (backendType === 'database') {
    return { type: type };
  }
  return;
}

export default helper(secretQueryParams);
