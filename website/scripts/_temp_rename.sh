# Renames the slew of markdown extensions in middleman all to .mdx
# Call with the path to the root folder, will convert recursively
# For example, bash _temp_rename.bash pages/nomad
# This file can be removed once we have finished porting from the old version!

find $1 -name "*.html.md" -exec rename 's/\.html.md$/.mdx/' '{}' \;
find $1 -name "*.html.markdown" -exec rename 's/\.html.markdown$/.mdx/' '{}' \;
find $1 -name "*.html.md.erb" -exec rename 's/\.html.md.erb$/.mdx/' '{}' \;
find $1 -name "*.md" -exec rename 's/\.md$/.mdx/' '{}' \;
