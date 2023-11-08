# OverviewCard

`<OverviewCard>` is often used in dashboards to display information. The `<div>` surrounding the card manages the layout for multiple cards

```hbs preview-template
<div class='flex row-wrap column-gap-16'>
  <OverviewCard
    @cardTitle='My card title'
    @subText='A description about the information in this card'
    @actionText='A link somewhere'
    @actionTo='docs'
  >
    <Hds::Text::Display class='has-top-padding-m' @tag='h2' @size='500'>
      Something yielded
    </Hds::Text::Display>
  </OverviewCard>
  <OverviewCard
    @cardTitle='Number of planets'
    @subText="In August 2006 the International Astronomical Union (IAU) tragically downgraded the status of Pluto to that of 'dwarf planet.'"
  >
    <Hds::Text::Display class='has-top-padding-m' @tag='h2' @size='500'>
      8
    </Hds::Text::Display>
  </OverviewCard>
</div>
```
