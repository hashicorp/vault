<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/confirm.js. To make changes, first edit that file and run "yarn gen-story-md confirm" to re-generate the content.-->

## Confirm
`Confirm` components prevent users from performing actions they do not intend to. This component should always be rendered with a Trigger (usually a link or button) and Message.

**Example**
  
```js
<div class="box">
  <Confirm as |c|>
    <nav class="menu">
      <ul class="menu-list">
        <li class="action">
          <c.Trigger>
            <button
              type="button"
              class="link is-destroy"
              onclick={{action c.onTrigger id}}>
              Delete
            </button>
          </c.Trigger>
        </li>
      </ul>
    </nav>
    <c.Message
      @id={{item.id}}
      @onCancel={{action c.onCancel}}
      @onConfirm={{onConfirm}}
      @title={{title}}
      @message={{message}}
      @confirmButtonText={{confirmButtonText}}
      @cancelButtonText={{cancelButtonText}}>
    </c.Message>
  </Confirm>
</div>
```

**See**

- [Uses of Confirm](https://github.com/hashicorp/vault/search?l=Handlebars&q=Confirm+OR+confirm)
- [Confirm Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/confirm.js)

---
