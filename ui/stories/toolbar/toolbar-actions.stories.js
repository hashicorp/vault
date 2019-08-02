/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select } from '@storybook/addon-knobs';
import notes from './toolbar-actions.md';

storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ToolbarActions`, () => ({
    template: hbs`
        <h5 class="title is-5">ToolbarActions</h5>
        <Toolbar>
          <ToolbarActions>
            <ToolbarLink
              @type="add"
              @params={{array "#"}}
            >
              Add item
            </ToolbarLink>
          </ToolbarActions>
        </Toolbar>
    `,
    context: {},
  }),
  {notes}
);
