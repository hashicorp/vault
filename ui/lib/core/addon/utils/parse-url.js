// adapted from https://gist.github.com/jed/964849
let fn = (function(anchor) {
  return function(url) {
    anchor.href = url;
    let parts = {};
    for (let prop in anchor) {
      if ('' + anchor[prop] === anchor[prop]) {
        parts[prop] = anchor[prop];
      }
    }

    return parts;
  };
})(document.createElement('a'));

export default fn;
