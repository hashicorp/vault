# Enterprise Alert Component

This component is an easy way to mark some content as only applicable to the enterprise version of vault. It can be used in any documentation pages in a variety of ways. The basic implementation is written as such, on its own line within a markdown file:

```jsx
<EnterpriseAlert />
```

And renders [like this](https://p176.p0.n0.cdn.getcloudapp.com/items/geuWOzkz/Screen%20Shot%202020-05-08%20at%204.17.34%20PM.png?v=2ace1c70f48cf1bbdd17f9ce96684453)

The default text can also be replaced with custom text as such:

```jsx
<EnterpriseAlert>
  Custom text <a href="">with a link</a>
</EnterpriseAlert>
```

Which renders [as such](https://p176.p0.n0.cdn.getcloudapp.com/items/v1uDE2vQ/Screen%20Shot%202020-05-08%20at%204.18.22%20PM.png?v=3a45268830fac868be50047060bb4303)

Finally, it can be rendered inline as a "tag" to mark a section or option as enterprise only by adding the `inline` attribute:

```jsx
<EnterpriseAlert inline>
```

This is typically used after a list item, or after a headline. It renders [as such](https://p176.p0.n0.cdn.getcloudapp.com/items/KouqnrOm/Screen%20Shot%202020-05-08%20at%204.16.34%20PM.png?v=ac21328916aa98a1a853cde5989058bd)
