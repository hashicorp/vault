import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { text, withKnobs, boolean } from '@storybook/addon-knobs';
import notes from './linkable-item.md';

storiesOf('LinkableItem', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `LinkableItem with attributes`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Linkable Item</h5>
        <LinkableItem @disabled={{disabled}} @link={{hash route="vault" model="myModel"}} as |Li|>
          <Li.content @accessor={{accessor}} @link={{hash route="vault" model="myModel"}} @title={{title}} @description={{description}} />
          <Li.menu>{{menu}}</Li.menu>
        </LinkableItem>
      `,
      context: {
        disabled: boolean('disabled', false),
        title: text('title', 'My Title'),
        accessor: text('accessor', 'v2 secret'),
        description: text(
          'description',
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam volutpat vulputate lacus sit amet lobortis. Nulla fermentum porta consequat. Mauris porttitor libero nibh, ac facilisis ex molestie non. Nulla dolor est, pharetra et maximus vel, varius eu augue. Maecenas eget nisl convallis, vehicula massa quis, pharetra justo. Praesent porttitor arcu at gravida dignissim. Vestibulum condimentum, risus a fermentum pulvinar, enim massa venenatis velit, a venenatis leo sem eget dolor. Morbi convallis dui sit amet egestas commodo. Nulla et ultricies leo.'
        ),
        menu: text('menu component', 'menu button goes here'),
      },
    }),
    { notes }
  )
  .add(
    `LinkableItem with blocks`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Linkable Item with block content</h5>
        <LinkableItem @disabled={{disabled}} @title={{title}} @link={{hash route="vault" model="myModel"}} as |Li|>
          <Li.content @accessor={{accessor}} @description={{description}} />
          <Li.menu>{{menu}}</Li.menu>
        </LinkableItem>
      `,
      context: {
        disabled: boolean('disabled', false),
      },
    }),
    { notes }
  );
