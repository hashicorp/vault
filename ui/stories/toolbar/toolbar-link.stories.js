/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select, text } from '@storybook/addon-knobs';
import notes from './toolbar-link.md';


storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ToolbarLink`,() => ({
    template: hbs`
      <h5 class="title is-5">ToolbarLink</h5>
      <div style="width: 400px;">
        <Toolbar>
          <ToolbarActions>
            <ToolbarLink
              @params={{array '#'}}
              @type={{type}}
            >
              {{label}}
            </ToolbarLink>
          </ToolbarActions>
        </Toolbar>
      </div>
    `,
    context: {
      type: select('Type', ['', 'add']),
      label: text('Button text', 'Edit secret'),
    },
  }),
  {notes}
);
