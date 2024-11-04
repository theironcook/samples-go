### Steps to run this sample:
1) Configure a [Temporal Server](https://github.com/temporalio/samples-go/tree/main/#how-to-use) (such as Temporal Cloud) with mTLS.

2) Run the following command to start the worker
```
go run ./helloworldmtls/worker -target-host enrichments-testing-b1.11a66.tmprl.cloud:7233 -namespace enrichments-testing-b1.11a66 -client-cert /Users/bartwood/projects/TemporalCerts/b1_cert.pem -client-key /Users/bartwood/projects/TemporalCerts/b1_key.pem
```
3) Run the following command to start the example
```
go run ./helloworldmtls/starter -target-host enrichments-testing-b1.11a66.tmprl.cloud:7233 -namespace enrichments-testing-b1.11a66 -client-cert /Users/bartwood/projects/TemporalCerts/b1_cert.pem -client-key /Users/bartwood/projects/TemporalCerts/b1_key.pem
```

If the server uses self-signed certificates and does not have the SAN set to the actual host, pass one of the following two options when starting the worker or the example above:
1. `-server-name` and provide the common name contained in the self-signed server certificate
2. `-insecure-skip-verify` which disables certificate and host name validation
