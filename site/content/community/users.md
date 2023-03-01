+++
title = "Who is using wazero"
+++

## Who is using wazero?

Below are an incomplete list of projects that use wazero strategically as a
part of their open source or commercial work. Please support our community by
considering their efforts before starting your own!

### General purpose plugins

| Name           | Description                                                                   |
|:---------------|-------------------------------------------------------------------------------|
| [go-plugin][2] | implements [Protocol Buffers][8] services with WebAssembly vi code generation |
| [waPC][5]      | implements [Apex][6] interfaces with WebAssembly via code generation          |
| [scale][13]    | implements [Polyglot][14] interfaces with WebAssembly via code generation     | 

### Network Proxies

| Name      | Description                                         |
|:----------|-----------------------------------------------------|
| [mosn][9] | implements 3rd party extension via [proxy-wasm][10] |

### Middleware

| Name                   | Description                                        |
|:-----------------------|----------------------------------------------------|
| [http-wasm-host-go][3] | serves HTTP handlers implemented in [http-wasm][4] |

### Libraries

| Name             | Description                                  |
|:-----------------|----------------------------------------------|
| [go-re2][7]      | high performance regular expressions         |
| [go-sqlite3][11] | [SQLite][12] bindings, `database/sql` driver |

## Updating this list

This is a community maintained list. It may have an inaccurate or outdated
entries, or missing something entirely. Changes to the [source][1] are
welcome, but please be conscious that not all projects desire to be on lists.
To ensure we promote community members, please do not add works that don't use
wazero to this list. Please keep descriptions short for a better table
experience.

[1]: https://github.com/tetratelabs/wazero/tree/main/site/content/community/users.md

[2]: https://github.com/knqyf263/go-plugin

[3]: https://github.com/http-wasm/http-wasm-host-go

[4]: https://http-wasm.io

[5]: https://wapc.io

[6]: https://apexlang.io

[7]: https://github.com/wasilibs/go-re2

[8]: https://protobuf.dev/overview/

[9]: https://mosn.io/

[10]: https://github.com/proxy-wasm/spec

[11]: https://github.com/ncruces/go-sqlite3

[12]: https://sqlite.org

[13]: https://scale.sh

[14]: https://github.com/loopholelabs/polyglot-go