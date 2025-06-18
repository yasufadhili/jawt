# Just A Weird Web Tool


## Syntax

### Pages

```jml

_doctype page index

import Layout from "somponents/layout"

Page {
  title: "Homepage"
  description: "Built with JAWT"
  
  Layout {}
  
}

```

### Components

```jml

_doctype component Layout

Row {
  style: "flex-1 p-4"
}

```

## CLI Commands

```shell
jawt init my_app
```

```shell
jawt run
```

```shell
jawt build
```
