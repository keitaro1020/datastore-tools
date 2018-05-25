# datastore-tool
CLI for Google Cloud Datastore

## Before you begin

- Create a [service account](https://cloud.google.com/docs/authentication/getting-started#creating_a_service_account).

## Installing the tools

`go get github.com/keitaro1020/datastore-tools`

## Using the tools

- select 
```
Usage:
  datastore-tools select [flags]

Flags:
  -p, --project string     datastore project id [required]
  -k, --kind string        datastore kind [required]
  -n, --namespace string   datastore namespace
  -f, --key-file string    gcp service account JSON key file
  -t, --table              output table view
  -c, --count              count only
  -w, --filter string      filter query (Property=Value)
```

- truncate
```
Usage:
  datastore-tools truncate [flags]

Flags:
  -p, --project string     datastore project id [required]
  -k, --kind string        datastore kind [required]
  -n, --namespace string   datastore namespace
  -f, --key-file string    gcp service account JSON key file
```