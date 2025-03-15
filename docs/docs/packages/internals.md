---
sidebar_position: 20
---

# Internals

Those astute among you may have noticed that `zog` has an `internals` package. You may be right in thinking that this is a bit weird as standard practice in golang is to name the package for internal code `internal` which allows golang to not export that code to user space.

And you would be right. However, `zog` takes a different approach. Our `internals` package holds code that is not meant for user space just like a typical `internal` package. Its not recommended that you use it as code inside it may have breaking changes at any time. However you can if you need or want to. This is a nice way for us to experiment with API's and allow you to build things on top of experimental code. Often time a feature will be hidden inside the internals package for a long time before it is promoted to the main package.
