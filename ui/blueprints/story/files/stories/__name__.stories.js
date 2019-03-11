/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';

storiesOf('<%= classifiedModuleName %>/', module)
  .addParameters({ options: { showPanel: false } })
  .add(`<%= classifiedModuleName %>`, () => ({
    template: hbs`
        <h5 class="title is-5"></h5>
        <<%= classifiedModuleName %>/>
    `,
    context: {},
  }));
