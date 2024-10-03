# tender

tender is template engine that has compatibility with Terraform string template.

## Installation

Download binary from [releases page](https://github.com/ysugimoto/tender/releases) according to your platform and place it under the `$PATH`, or you can install via `go get` command:

```shell
go get github.com/ysugimoto/tender
```

## Usage (CLI)

You can execute templating from cli.

```shell
tender /path/to/template/file.tpl
```

## Usage (Programmable)

Mostly you will use as templating library.

```go
package main

import (
    "fmt"
    "strings"

    "github.com/ysugimoto/tender"
)

func main() {
    tmpl := "Hello ${name}!"

    rendered := tender.Must(tender.Render(tmpl, map[string]any{
        "name": "tender",
    }))

    fmt.Println(rendered) //=> Hello tender!"
}
```

## Control Syntax

`tender` has some control syntax that has Terraform string template.
The control syntax must be enclosed `%{ ... }`.

### for
`for` control can loop rendering inside a block.

```
%{ for v in list }
This is loop block for list variable.
%{ endfor }
```

### If-elseif-else

`if` control can switch rendering block from provided condition.

```
%{ if i == 0 }
This is block if i is 0.
%{ else if i == 1 }
This is block if i is 1.
%{ elseif i == 2 }
This is block if i is 2.
${ else }
This is block if i is other value.
%{ endif }
```

#### Operator

`tender` recognize following operators:

| operator   | comparison                            |
|:----------:|:--------------------------------------|
| a `==` b   | compare a is equal to b.              |
| a `!=` b   | compare a is not equal to b.          |
| a `>` b    | compare a is greater than b.          |
| a `>=` b   | compare a is greater than equal to b. |
| a `<` b    | compare a is less than b.             |
| a `<=` b   | compare a is less than equal to b.    |
| a `&&` b   | a and b are truthy.                   |
| a `\|\|` b | a or b is truthy.                     |

> [!NOTE]
> Template assigned variables are readonly. Therefore you can't do arithmetic operations variable in template like `x + 1`.

## Interporation

Template variable will be interporated in `${...}` synatax.

```
The template vairable is ${value}.
```

If you provide "value" variable with "tender", the result will be `The template variable is tender`.

### Environment variables

`tender` can also reference environment variable if interporation name is `[A-Z_]+` format.

```
The environment vairable is ${SERVICE_NAME}.
```

If you specify "SERVICE_NAME" environment variable with "tender", the result will be `The environment variable is tender`.

## Performance Benchmark

```
goos: darwin
goarch: arm64
pkg: benchmark
BenchmarkRaymondRender-10          21631             54697 ns/op            9719 B/op        278 allocs/op
BenchmarkNativeRender-10          120549              9344 ns/op            7556 B/op        141 allocs/op
BenchmarkTenderRender-10           77595             16370 ns/op           18521 B/op        442 allocs/op
```

We need to improve to reduce allocation X(

## Contribution

- Fork this repository
- Customize / Fix problem
- Send PR :-)
- Or feel free to create issues for us. We'll look into it

## License

MIT License

## Contributors

- [@ysugimoto](https://github.com/ysugimoto)
