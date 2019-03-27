import { configure, addParameters, addDecorator } from '@storybook/ember';
import Centered from '@storybook/addon-centered/ember';

function loadStories() {
  // automatically import all files ending in *.stories.js
  const req = require.context('../stories/', true, /.stories.js$/);
  req.keys().forEach(filename => req(filename));
}

addDecorator(Centered);

configure(loadStories, module);
