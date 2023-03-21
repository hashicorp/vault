# using this to test locally, will be deleting this 
echo "Installing staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@2023.1.2 #v0.4.2

echo "Get files changed in the PR"
changed_files=$(git diff --name-only main)
echo $changed_files
for val in ${changed_files[@]}; do
  echo "$val changed"
done


echo "Running staticcheck to look for deprecations"

# use entire package name for packages
ignoreList=()

filesChanged=("builtin/logical/pki/chain_test.go"
"builtin/logical/database/path_roles.go"
"vault/diagnose/tls_verification.go")

# run staticcheck
staticcheck ./... | grep deprecated > staticcheckOutput.txt  

# include details of only changed files in the PR
count=0
for val in ${changed_files[@]}; do
    if grep -q  $val staticcheckOutput.txt; then
        grep $val staticcheckOutput.txt
        count=$((count+1))
    fi
done

# mv tmpfile staticcheckOutput.txt
# rm -rf tmpfile

# # delete ignored values from the output
# for val in ${ignoreList[@]}; do
#   grep -v $val staticcheckOutput.txt > tmpfile && mv tmpfile staticcheckOutput.txt
#   rm -rf tmpfile
# done
rm -rf staticcheckOutput.txt  


if [ "$count" -ne "0" ]
then
    echo $count
    echo "Use of deprecated function, variable, constant or field found"
    #  cat staticcheckOutput.txt 
    #  # output file clean up 
    #  rm -rf staticcheckOutput.txt  
     exit 1 
else
     echo "Use of deprecated function, variable, constant or field not found"
    #  rm -rf staticcheckOutput.txt  
fi





          
