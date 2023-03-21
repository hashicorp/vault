# using this to test locally, will be deleting this 
echo "Installing staticcheck"
go install honnef.co/go/tools/cmd/staticcheck@latest


echo "Running staticcheck to look for deprecations"

# use entire package name for packages
ignoreList=("github.com/golang/protobuf/proto"
"errwrap.Wrapf"
"io/ioutil"
"ptypes"
".RawRequestWithContext"
".RawRequest"
"cloud.google.com/go/monitoring/apiv3"
"crypto/dsa"
"x509.ParseCRL"
"ACL().Create"
"ACL().Destroy" 
"Statements.CreationStatements"
"Statements.RevocationStatements"
"Statements.RollbackStatements"
"strings.Title"
"connState.NegotiatedProtocolIsMutual")

filesChanged=("builtin/logical/pki/chain_test.go"
"builtin/logical/database/path_roles.go"
"vault/diagnose/tls_verification.go")

# run staticcheck
staticcheck ./... | grep deprecated > staticcheckOutput.txt  

# include details of only changed files in the PR
for val in ${filesChanged[@]}; do
  grep $val staticcheckOutput.txt >> tmpfile
done

mv tmpfile staticcheckOutput.txt
rm -rf tmpfile

# delete ignored values from the output
for val in ${ignoreList[@]}; do
  grep -v $val staticcheckOutput.txt > tmpfile && mv tmpfile staticcheckOutput.txt
  rm -rf tmpfile
done


if [ -s staticcheckOutput.txt ]
then
     echo "Use of deprecated function, variable, constant or field found"
     cat staticcheckOutput.txt 
     # output file clean up 
     rm -rf staticcheckOutput.txt  
     exit 1 
else
     echo "Use of deprecated function, variable, constant or field not found"
     rm -rf staticcheckOutput.txt  
fi





          
