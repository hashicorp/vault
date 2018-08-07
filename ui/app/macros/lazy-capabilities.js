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
