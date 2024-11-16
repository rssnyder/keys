# keys

## read

```
curl https://keys.rileysnyder.dev/foo
```

## write

```
curl -d 'bar' https://keys.rileysnyder.dev/foo
```

you can also create random keys

```
curl -d 'bar' https://keys.rileysnyder.dev
```

## development notes

1. test with `make test`
2. need to split out routes eventually to leverage httptest
