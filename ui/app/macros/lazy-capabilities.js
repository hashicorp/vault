// usage:
//
// import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
//
// export default DS.Model.extend({
//   //pass the template string as the first arg, and be sure to use '' around the
//   //paramerters that get interpolated in the string - that's how the template function
//   //knows where to put each value
//   zeroAddressPath: lazyCapabilities(apiPath`${'id'}/config/zeroaddress`, 'id'),
//
// });
//

import { queryRecord } from 'ember-computed-query';

export function apiPath(strings, ...keys) {
  return function(data) {
    let dict = data || {};
    let result = [strings[0]];
    keys.forEach((key, i) => {
      result.push(dict[key], strings[i + 1]);
    });
    return result.join('');
  };
}

export default function() {
  let [templateFn, ...keys] = arguments;
  return queryRecord(
    'capabilities',
    context => {
      return {
        id: templateFn(context.getProperties(...keys)),
      };
    },
    ...keys
  );
}
