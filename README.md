# hookup

hookup get webhooks from some service and proxy requests by config file.

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
./hookup -c config.yml
```

you set `http://some.hookup.domain/<from name>`.

ex)

- github webhook url `http://some.hookup.domain/github`
- bitbucket webhook url `http://some.hookup.domain/bitbucket`

hookup proxy `/github` webhook to localhost:1234 and `/bitbucket` webhook to localhost:1235.
