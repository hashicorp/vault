(function() {
  function vendorModule() {
    'use strict';

    return { 'default': self['autosize'] };
  }

  define('autosize', [], vendorModule);
})();
