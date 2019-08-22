<!--THIS FILE IS AUTO GENERATED. This file is generated from JSDoc comments in lib/core/addon/components/confirm.js. To make changes, first edit that file and run "yarn gen-story-md confirm" to re-generate the content.-->

## Confirm
`Confirm` components prevent users from performing actions they do not intend to by showing a confirmation message as an overlay. This is a contextual component that should always be rendered with a `Message` which triggers the message.

See the `Message` component for a description of properties.

**Example**
  
```js
<div class="box">
  <Confirm as |c|>
    <c.Message
      @id={{item.id}}
      @triggerText="Delete"
      @message="This will permanently delete this secret and all its versions."
      @onConfirm={{action "delete" item "secret"}}
      />
  </Confirm>
</div>
```

**See**

- [Uses of Confirm](https://github.com/hashicorp/vault/search?l=Handlebars&q=Confirm+OR+confirm)
- [Confirm Source Code](https://github.com/hashicorp/vault/blob/master/ui/lib/core/addon/components/confirm.js)

---
