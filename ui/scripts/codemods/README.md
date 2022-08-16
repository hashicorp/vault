# Running Codemods

The handlebars codemods use [ember-template-recast](https://github.com/ember-template-lint/ember-template-recast) and can be run with the following:

- navigate to the UI directory of the Vault project
- execute `npx ember-template-recast "**/*.hbs" -t ./path/to/transform-file.js`

This will run the transform on all .hbs files within the ui directory which covers the app and all addons.
The terminal will output the number of files processed as well as the number of changed, unchanged, skipped and errored files.
It's a good idea to validate the output to ensure that the intended transforms have taken place.
If there are issues with some of the files, simply revert the changes via git, tweak the codemod and run again.

## Example
`npx ember-template-recast "**/*.hbs" -t ./scripts/codemods/no-quoteless-attributes.js`