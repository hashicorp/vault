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

import { maybeQueryRecord } from 'vault/macros/maybe-query-record';

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
  return maybeQueryRecord(
    'capabilities',
    context => {
      // pull all context attrs
      let contextObject = context.getProperties(...keys);
      // remove empty ones
      let nonEmptyContexts = Object.keys(contextObject).reduce((ret, key) => {
        if (contextObject[key] != null) {
          ret[key] = contextObject[key];
        }
        return ret;
      }, {});
      // if all of them aren't present, cancel the fetch
      if (Object.keys(nonEmptyContexts).length !== keys.length) {
        return;
      }
      // otherwise proceed with the capabilities check
      return {
        id: templateFn(nonEmptyContexts),
      };
    },
    ...keys
  );
}
