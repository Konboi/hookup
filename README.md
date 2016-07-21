# hookup

hookup get request and proxy.

# How to use

## Config file

```yaml
port: 9876
hooks:
  - from: github
    to: localhost:1234
  - from: bitbucket
    to: localhost:1235
```

## Useage

```
go get github.com/Konboi/hookup
hookup -c config.yml
```

you set `http://some.hookup.domain/<from name>`.

ex)

- github webhook url `http://some.hookup.domain/github`
- bitbucket webhook url `http://some.hookup.domain/bitbucket`

hookup proxy `/github` request to localhost:1234 and `/bitbucket` request to localhost:1235.
