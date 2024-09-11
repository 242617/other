# other

Messing around with Golang.

## tasks

* Linked list implementation
* Leetcode challenges

## generator

Generates words combinations. Available charsets:
* `generator.C_09`: all digits ([0-9])
* `generator.C_AZ`: capital letters ([A-Z])
* `generator.C_az`: letters ([a-z])

### Import

```
import "github.com/242617/other/generator"
```

### Usage

```
ch := generator.New(4, generator.C_az|generator.C_09)
for word := range ch {
    fmt.Println(word)
}
```

Will generate
```
0000
0001
0002
...
zzzy
zzzz
```