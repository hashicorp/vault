# Testing

To run these tests, first start MSSQL in Docker. Please do make sure to view the EULA before 
accepting it as it includes limits on the number of users per company who can be using the 
image, and how it can be used in testing.

```
sudo docker run -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=<YourStrong!Passw0rd>' \
   -p 1433:1433 --name sql1 \
   -d mcr.microsoft.com/mssql/server:2017-latest
```

Then use the following env variables for testing:

```
export VAULT_ACC=1
export MSSQL_URL="sqlserver://SA:%3CYourStrong%21Passw0rd%3E@localhost:1433"
```

Note that the SA password passed into the Docker container differs from the one passed into the tests.
It's the same password, but Go's libraries require it to be percent encoded.

Running all the tests at once against one Docker container will likely fail because they interact with
each other. Consider running one test at a time.
