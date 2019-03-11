import { configure, addParameters } from '@storybook/ember';

function loadStories() {
  // automatically import all files ending in *.stories.js
  const req = require.context('../stories', true, /.stories.js$/);
  req.keys().forEach(filename => req(filename));
}

configure(loadStories, module);
