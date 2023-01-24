import { helper } from '@ember/component/helper';

export function secretQueryParams([backendType, type = ''], { asQueryParams }) {
  const values = {
    transit: { tab: 'actions' },
    database: { type },
    keymgmt: { itemType: type === 'provider' ? 'provider' : 'key' },
  }[backendType];
  // format required when using LinkTo with positional params
  if (values && asQueryParams) {
    return {
      isQueryParams: true,
      values,
    };
  }
  return values;
}

export default helper(secretQueryParams);
