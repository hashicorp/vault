# find plugins that will NOT be in minimal 
while IFS= read -r line; do
	# find import entry of "github.com/hashicorp/<plugin_name>.*"
  found=$(grep -r \"github.com/hashicorp/${line}.*\" ./)

  num_lines=$(echo -e "$found" | wc -l)

  # if it is referenced once (by helper/builtinplugins/registry.go),
  # it can be safely removed from minimal
  if [ "$num_lines" -eq 1 ]; then
    echo "$line NOT in minimal"
    echo "$found\n"
  fi
done < "$1"
