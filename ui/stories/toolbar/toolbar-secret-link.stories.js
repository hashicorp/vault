/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select, text } from '@storybook/addon-knobs';
import notes from './toolbar-secret-link.md';


storiesOf('Toolbar/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ToolbarSecretLink`,() => ({
    template: hbs`
      <h5 class="title is-5">ToolbarLink</h5>
      <div style="width: 400px;">
        <Toolbar>
          <ToolbarActions>
            <ToolbarSecretLink
              @secret={{model.id}}
              @mode="edit"
              @data-test-edit-link=true
              @replace=true
              @type={{type}}
            >
              {{label}}
            </ToolbarSecretLink>
          </ToolbarActions>
        </Toolbar>
      </div>
    `,
    context: {
      type: select('Type', ['', 'add']),
      label: text('Button text', 'Edit role'),
    },
  }),
  {notes}
);
