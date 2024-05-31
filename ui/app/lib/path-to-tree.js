/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { unflatten } from 'flat';
import deepmerge from 'deepmerge';

const DOT_REPLACEMENT = '☃';

//function that takes a list of path and returns a deeply nested object
//representing a tree of all of those paths
//
//
// given ["foo", "bar", "foo1", "foo/bar", "foo/baz", "foo/bar/baz"]
//
// returns {
//    bar: null,
//    foo: {
//      bar: {
//        baz: null
//      },
//      baz: null,
//    },
//    foo1: null,
// }
export default function (paths) {
  // first sort the list by length, then alphanumeric
  const list = paths.slice(0).sort((a, b) => b.length - a.length || b.localeCompare(a));
  // then reduce to an array
  // and we remove all of the items that have a string
  // that starts with the same prefix from the list
  // so if we have "foo/bar/baz", both "foo" and "foo/bar"
  // won't be included in the list
  let tree = list.reduce((accumulator, ns) => {
    const nsWithPrefix = accumulator.find((path) => path.startsWith(ns));
    // we need to make sure it's a match for the full path part
    const isFullMatch = nsWithPrefix && nsWithPrefix.charAt(ns.length) === '/';
    if (!isFullMatch) {
      accumulator.push(ns);
    }
    return accumulator;
  }, []);

  tree = tree.sort((a, b) => a.localeCompare(b));
  // after the reduction we're left with an array that contains
  // strings that represent the longest branches
  // we'll replace the dots in the paths, then expand the path
  // to a nested object that we can then query with Ember.get
  return deepmerge.all(
    tree.map((p) => {
      p = p.replace(/\.+/g, DOT_REPLACEMENT);
      return unflatten({ [p]: null }, { delimiter: '/', object: true });
    })
  );
}
