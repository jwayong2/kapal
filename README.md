# kapal

Kapal is an orchestration tool utilizing several Linux file systems to provide smart volume capabilities to docker containers that includes cloning/snapshotting, backup, replication, in-memory volume.

# Contributing

1. Clone the project into `$KAPAL_HOME`

2. Install all the project dependencies into your `$GOPATH`
```
make deps
```

3. Run all the tests from `$KAPAL_HOME`
```
go test ./...
```

