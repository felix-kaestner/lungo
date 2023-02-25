# Lungo

<p align="center">
    <span>A tiny, zero-dependency web framework based on net/http with an intuitive API.</span>
    <br><br>
    <a href="https://github.com/felix-kaestner/lungo/issues">
        <img alt="Issues" src="https://img.shields.io/github/issues/felix-kaestner/lungo?color=29b6f6&style=flat-square">
    </a>
    <a href="https://github.com/felix-kaestner/lungo/stargazers">
        <img alt="Stars" src="https://img.shields.io/github/stars/felix-kaestner/lungo?color=29b6f6&style=flat-square">
    </a>
    <a href="https://github.com/felix-kaestner/lungo/blob/main/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/felix-kaestner/lungo?color=29b6f6&style=flat-square">
    </a>
    <a href="https://pkg.go.dev/github.com/felix-kaestner/lungo">
        <img alt="Stars" src="https://img.shields.io/badge/go-documentation-blue?color=29b6f6&style=flat-square">
    </a>
    <a href="https://goreportcard.com/report/github.com/felix-kaestner/lungo">
        <img alt="Issues" src="https://goreportcard.com/badge/github.com/felix-kaestner/lungo?style=flat-square">
    </a>
    <a href="https://codecov.io/gh/felix-kaestner/lungo">
        <img src="https://img.shields.io/codecov/c/github/felix-kaestner/lungo?style=flat-square&token=KK7ZG7A90X"/>
    </a>
    <a href="https://twitter.com/kaestner_felix">
        <img alt="Twitter" src="https://img.shields.io/badge/twitter-@kaestner_felix-29b6f6?style=flat-square">
    </a>
</p>

## Quickstart

```go
package main

import "github.com/felix-kaestner/lungo"

func main() {
    app := lungo.New()

    app.Get("/", func(c *lungo.Context) error {
        return c.Text("Hello, World!")
    })

    app.Listen(":3000")
}
```

##  Installation

Install Lungo with the `go get` command:

```
$ go get -u github.com/felix-kaestner/lungo
```

## Contribute

All contributions in any form are welcome! 🙌🏻  
Just use the [Issue](.github/ISSUE_TEMPLATE) and [Pull Request](.github/PULL_REQUEST_TEMPLATE) templates and I'll be happy to review your suggestions. 👍

## Support

Any kind of support is well appreciated! 👏🏻  
If you want to tweet about the project, make sure to tag me [@kaestner_felix](https://twitter.com/kaestner_felix). You can also support my open source work on [GitHub Sponsors](https://github.com/sponsors/felix-kaestner).

---

Released under the [MIT License](LICENSE).
