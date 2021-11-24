# zip

[![Build Status](https://travis-ci.org/josa42/go-zip.svg?branch=master)](https://travis-ci.org/josa42/go-zip)

**🚧  Work in progress**

## Examples

See also: [`examples/main.go`](examples/main.go).

### Compress

```go
a, _ := zip.CreateArchive("test.zip")
defer a.Close()

a.Add(".", "test")
```

Optionally function function can be provided, to ignore specific files or
directories.

```go
a, _ := zip.CreateArchive("test.zip")
defer a.Close()

a.Add(".", "test", func(path string, sourcePath string) bool {

  // The file / directory will be ignored if `false` is returned
  if path == ".git" {
    fmt.Println("> Ignore:", path)
    return false
  }

  fmt.Println("> Add:   ", path)
  return true
})
```

### List archive content

```go
	a, _ := zip.OpenArchive("test.zip")
	defer a.Close()

	filePaths, _ := a.List()
	for _, f := range filePaths {
		fmt.Println(f)
	}
```

## Todo

- [ ] Extract archive

## License

[MIT © Josa Gesell](LICENSE)