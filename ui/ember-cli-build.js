/*jshint node:true*/
/* global require, module */
var EmberApp = require('ember-cli/lib/broccoli/ember-app');

module.exports = function(defaults) {
  var config = defaults.project.config(EmberApp.env());
  var app = new EmberApp(defaults, {
    favicons: {
      faviconsConfig: {
        appName: 'Vault Enterprise',
        path: config.rootURL,
        url: null,
        icons: {
          android: false,
          appleIcon: false,
          appleStartup: false,
          coast: false,
          favicons: true,
          firefox: false,
          opengraph: false,
          twitter: false,
          windows: false,
          yandex: false
        }
      }
    },
    codemirror: {
      modes: ['javascript','ruby'],
      keyMaps: ['sublime']
    },
    babel: {
      plugins: [
        'transform-object-rest-spread'
      ]
    }
  });

  app.import('vendor/string-includes.js');
  app.import(app.bowerDirectory + '/string.prototype.startswith/startswith.js');
  app.import(app.bowerDirectory + '/autosize/dist/autosize.js');
  app.import('vendor/shims/autosize.js');

  app.import(app.bowerDirectory + '/jsonlint/lib/jsonlint.js');
  app.import(app.bowerDirectory + '/codemirror/addon/lint/lint.css');
  app.import(app.bowerDirectory + '/codemirror/addon/lint/lint.js');
  app.import(app.bowerDirectory + '/codemirror/addon/lint/json-lint.js');
  app.import(app.bowerDirectory + '/base64-js/base64js.min.js');
  app.import(app.bowerDirectory + '/text-encoder-lite/index.js');
  app.import(app.bowerDirectory + '/Duration.js/duration.js');

  // Use `app.import` to add additional libraries to the generated
  // output files.
  //
  // If you need to use different assets in different
  // environments, specify an object as the first parameter. That
  // object's keys should be the environment name and the values
  // should be the asset to use in that environment.
  //
  // If the library that you are including contains AMD or ES6
  // modules that you would like to import into your application
  // please specify an object with the list of modules as keys
  // along with the exports of each module as its value.

  return app.toTree();
};
