# Tabs Component

> An MDX-compatible Tabs component

This React component renders tabbed content.

## Usage

- Use the `<Tabs>` tag in your markdown file to begin a tabbed content section.
- Use the `<Tab>` tag with a `heading` prop to separate your markdown

### Important

A line must be skipped between the `<Tab>` and your markdown (for both above and below said markdown). [This is a limitation of MDX also pointed out by the Docusaurus folks 🔗 ](https://v2.docusaurus.io/docs/markdown-features/#multi-language-support-code-blocks)

### Example

```mdx
<Tabs>
<Tab heading="CLI command">
             <!-- Intentionally skipped line.. -->
### Content
            <!-- Intentionally skipped line.. -->
</Tab>
<Tab heading="API call using cURL">

### Content

</Tab>
</Tabs>
```

### Component Props

`<Tabs>` can be provided any arbitrary `children` so long as the `heading` prop is present the React or HTML tag used to wrap markdown, that said, we provide the `<Tab>` component to separate your tab content without rendering extra, unnecessary markup.

This works:

```mdx
<Tabs>
<Tab heading="CLI command">

### Content

</Tab>
....
</Tabs>
```

This _does not_ work:

```mdx
<Tabs>
<Tab> <!-- missing the `heading` prop to provide a tab heading -->

### Content

</Tab>
....
</Tabs>
```
