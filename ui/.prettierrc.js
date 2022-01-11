'use strict';

module.exports = {
  singleQuote: true,
  trailingComma: 'es5',
  printWidth: 110,
  overrides: [
    {
      files: '*.hbs',
      options: {
        singleQuote: false,
        printWidth: 125,
      },
    },
  ],
};
