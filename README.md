# keys

## init

```sql
CREATE TABLE keys (
  key VARCHAR NOT NULL PRIMARY KEY,
  value TEXT NOT NULL 
);
```

## read

```
curl https://keys.rileysnyder.dev/foo
```

## write

```
curl -d 'bar' https://keys.rileysnyder.dev/foo
```
