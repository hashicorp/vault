import { helper } from '@ember/component/helper';

export default helper(function formatTtl([timestring], { removeZero = false }) {
  // Expects a number followed by one of s, m, or h
  // eg. 40m or 1h20m0s
  let matches = timestring?.match(/([0-9]+[h|m|s])/g);
  if (!matches) {
    return timestring;
  }
  return matches
    .map((set) => {
      // eslint-disable-next-line no-unused-vars
      let [_, number, unit] = set.match(/([0-9]+)(h|m|s)/);
      if (removeZero && number === '0') {
        return null;
      }
      const word = { h: 'hour', m: 'minute', s: 'second' }[unit];
      return `${number} ${number === '1' ? word : word + 's'}`;
    })
    .filter((s) => null !== s)
    .join(' ');
});
