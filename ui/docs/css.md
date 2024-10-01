# CSS/SCSS

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Guidelines](#guidelines)
  - [Helper classes](#helper-classes)
  - [Core class styles](#core-class-styles)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Guidelines

- [**Helper classes**](#helper-classes) should be used if a styling block does not already exist and/or a reasonable number of helper classes can cover the desired style.
- [**Core classes**](#core-class-styles) provide styling for any classes not associated with a component. The scope of each file is defined by the file name.
- **Component specific styling** should only be added to, or created when you cannot achieve the styling with a helper class or core class.
- **Utils'** files define mixins, and variables. 

> ### Known issues
> The following are known issues that we are working to address in conjunction with the adoption of HDS.
> 1. **Size variables** The UI does not follow a consistent size variable pattern. We use both px and rem to define font-size and we use px, rem, and ems to define margins, heights and widths. For accessibility reasons we _should_ define all font-sizing at the very least by rem, though this is not consistently done in the app.
> 2. **Variable naming** The UI does not have a consistent pattern to variable naming. We use a mix of numbers and words (ex: `ui-gray-050` is the same as `ui-gray-lightest`).
> 3. **Random variables** We have dieing but not dead variables. For example, we have some variables that define box-shadows and we have some variables to define animations, but we are missing many box-shadow definitions and we do not consistently use the animation variables.
> 4. **Missing variables** The UI does not have a variable for all commonly occurring sizes and colors. For example, we do not have a variable that covers the `14px` though it is a commonly used size.
> 5. **!Important** `!important` is sprinkled throughout helper, core and component files. Ideally, the cascading and order of styles would eliminate the need of this keyword. However, because `!important` exist randomly in many of our files, we now have cascading issues inside helper files and core files. In all known cases where these issues exist a comment has been left as to why the order of classes in that particular area matters.

### Helper classes

A good portion of our class definitions have come from Bulma. Bulma has since been removed, but we still rely on many of its class definitions. Bulma class definitions, specifically their helper classes, always end in the keyword  `!important`. This use of `important!` and Bulma specific naming patterns has led to a mix of inconsistent helper class names. Moving forward, we have agreed as a team to pursue the following patterns. When it makes sense, please default to these instead of relying on existing helper names for guidance.

- Drop the starting verb. Many of our helper classes start with a verb `has` or `is`. Going forward we prefer to drop the verb. `margin-bottom` instead of `is-margin-bottom`.
- Start your helper class name with what the class controls. `margin-bottom` instead of `bottom-margin`.
- Match your helper class size to a pre-existing size variable. `margin-bottom-large` instead of `margin-bottom-big`.

### Core class styles

All files under `app/styles/core` directory define style for the class of the file name. Think of these as files for the heavily used classes that are not defined as a component. Things like `.box` or `.title`. These classes are used in our app over many files that span multiple components, but they are themselves are _not_ components.

If the core file ends in a `s` (e.g. `lists` or `containers`) the plural indicates that the file defines more than just the style for the class `container`. The `core/containers.scss` file defines classes for all things relating to container-type classes: `page-container`, `section` as well as the `container` class. There are only a few plural core files as they are the exception and not the norm.
