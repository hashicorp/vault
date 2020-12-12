import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs } from '@storybook/addon-knobs';

storiesOf('<%= classifiedModuleName %>', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`<%= classifiedModuleName %>`, () => ({
    template: hbs`
      <h5 class="title is-5"><%= header %></h5>
      <<%= classifiedModuleName %>/>
    `,
    context: {},
  }));
