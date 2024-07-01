# oldstable

Check if version of go directive in go.mod is oldstable.

## As a CLI

### Usage

``` console
$ oldstable
Error: version of go directive in go.mod is not latest oldstable (oldstable: 1.21.11, current: 1.22.4)
```

## As a GitHub Action

### Usage

``` yaml
# .github/workflows/oldstable.yml
[...]
    steps:
      -
        name: Check if version of go directive in go.mod is oldstable
        uses: k1LoW/oldstable@v1
```
