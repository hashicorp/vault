/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './search-select.md';
import { withKnobs, text } from '@storybook/addon-knobs';

const onChange = (value) => alert(`New value is "${value}"`);
const models = ["policies/acl"];
storiesOf('SearchSelect/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs({ escapeHTML: false }))
  .add(`SearchSelect`, () => ({
    template: hbs`
        <h5 class="title is-5">Search Select</h5>
        <SearchSelect 
          @id="policies" 
          @models={{models}} 
          @onChange={{onChange}} 
          @helpText={{helpText}} 
          @label={{label}}
          @fallbackComponent="string-list" 
          @storyOptions={{storyOptions}}
        />
    `,
    context: {
      label: text("Label", "Policies"),
      helpText: text("Help Tooltip Text", "Policies associated with this group"),
      models: models,
      onChange: onChange,
      storyOptions: [{ id: "my-policy" }, { id: "my-other-policy" }, { id: "123" }, { id: "example-1" }]
    },
  }),
    { notes }
  );
