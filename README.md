# Google Cloud Datastore Tools
CLI for Google Cloud Datastore

## Before you begin

1. Create a [service account](https://cloud.google.com/docs/authentication/getting-started#creating_a_service_account).(roles: datastore.user)
2. Save JSON key file

## Installing the tools

```
go get github.com/keitaro1020/datastore-tools
cd $GOPATH/src/github.com/keitaro1020/datastore-tools
make
make install
```
Or download binaries from the [releases](https://github.com/keitaro1020/datastore-tools/releases) page(Linux/Windows/macOS).

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

- update
```
Usage:
  datastore-tools update [flags]

Flags:
  -p, --project string     datastore project id [required]
  -k, --kind string        datastore kind [required]
  -n, --namespace string   datastore namespace
  -f, --key-file string    gcp service account JSON key file
  -w, --filter string      filter query (Property=Value)
  -s, --set string         set update value (Property=Value)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Inspiration

https://github.com/boiyaa/datastore-tools