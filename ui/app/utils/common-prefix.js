import { get } from '@ember/object';

export default function(arr, attribute = 'id') {
  let content = arr || [];
  // this assumes an already sorted array
  // if the array is sorted, we want to compare the first and last
  // item in the array - if they share a prefix, all of the items do
  let firstString = get(content[0], attribute);
  let lastString = get(content[arr.length - 1], attribute);

  // the longest the shared prefix could be is the length of the match
  let targetLength = firstString.length;
  let prefixLength = 0;
  // walk the two strings, and if they match at the current length,
  // increment the prefixLength and try again
  while (
    prefixLength < targetLength &&
    firstString.charAt(prefixLength) === lastString.charAt(prefixLength)
  ) {
    prefixLength++;
  }
  // slice the prefix from the first item
  return firstString.substring(0, prefixLength);
}
