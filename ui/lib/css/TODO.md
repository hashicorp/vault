# Consolidating Bulma

The work required to consolidate bulma so we can eventually remove it

## Remove from main app style

- Disable bulma output from lib/css addon and see what breaks on build
- How to get bulma/sass/utilities/\_all.sass into main app.css build so we are no longer importing all of bulma
  - write custom from/until breakpoint mixins
  - $colors for buttons, form background, notifications, helpers needs to be replaced/revamped
