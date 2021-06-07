## golanglin-ci

```yaml
linters-settings:
  errcheck:
    check-type-assertions: true
  goconst:
    min-len: 2
    min-occurrences: 3
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
  govet:
    check-shadowing: true
  nolintlint:
    require-explanation: true
    require-specific: true

linters:
  disable-all: true
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - exportloopref
    - exhaustive
    - goconst
    - gocritic
    - gofmt
    - goimports
    - gomnd
    - gocyclo
    - gosec
    - gosimple
    - govet
    - ineffassign
    - misspell
    - nolintlint
    - nakedret
    - prealloc
    - predeclared
    - revive
    - staticcheck
    - structcheck
    - stylecheck
    - thelper
    - tparallel
    - typecheck
    - unconvert
    - unparam
    - varcheck
    - whitespace
    - wsl

run:
  issues-exit-code: 1

```

### ignore error
```go
func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Int()) //nolint # 这个忽略会忽略所有的lint
	fmt.Println(rand.Int()) //nolint:gosec # 忽略指定的lint
}

```

### 排除规则

```yaml
issues:
  exclude-rules:
    - path: _test\.go # 对 test 文件禁用指定的linters
	  linters:
	    - gocyclo
		- gosec
		- dupl
	- linters： # 排除一些linters消息
	  - gosec
	  text： "weak cryptographic primitive"
```

### 从指定的提交记录开始进行检测
```yaml
issues:
  new-from-rev: 02270a6 # 仅在 git revision 02270a6 后开始检测
```