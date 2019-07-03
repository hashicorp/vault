import { assert } from '@ember/debug';

// This is a tagged template function that will
// replace placeholders in the form of 'id' with the value from the passed context
//
// usage:
// let fn = apiPath`foo/bar/${'id'}`;
// let output = fn({id: 'an-id'});
// output will result in 'foo/bar/an-id';

export default function apiPath(strings, ...keys) {
  return function(data) {
    let dict = data || {};
    let result = [strings[0]];
    assert(
      `Expected ${keys.length} keys in apiPath context, only recieved ${Object.keys(data).join(',')}`,
      Object.keys(data).length >= keys.length
    );
    keys.forEach((key, i) => {
      result.push(dict[key], strings[i + 1]);
    });
    return result.join('');
  };
}
