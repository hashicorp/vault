module.exports = {
  ...require('@hashicorp/nextjs-scripts/.stylelintrc.js'),
  ignoreFiles: ['out/**'],
  rules: {
    'selector-pseudo-class-no-unknown': [
      true,
      {
        ignorePseudoClasses: ['first', 'last'],
      },
    ],
  },
  /* Specify overrides here */
}
