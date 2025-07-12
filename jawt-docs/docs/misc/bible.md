# 📖 The Book of JML — The Unofficial Jawt Bible

> *"In the beginning, there was HTML. Then chaos. Then JSX. And then, from the void, rose JML."*

---

## ✨ Chapter 1: The Origin of Structure

> *"And the developer said, let there be less noise. And JML appeared, pure and declarative."*

* You shall write no `<div>`.
* You shall not repeat thyself with boilerplate.
* The `_doctype` shall declare the purpose of every JML file.
* All things must begin with structure, for structure is clarity.

```jml
_doctype page home

Page {
  title: "Welcome"

  Container {
    Text {
      content: "And it was good."
    }
  }
}
```

---

## 🔥 Chapter 2: The Rules of Declaration

> *"Thou shalt describe thy UI, not assemble it."*

* Components are sacred.
* Props are passed, not injected.
* Logic resides in script. Structure lives in clarity.

```jml
_doctype component HelloButton

import script actions from "scripts/actions"

Button {
  text: "Say Hello"
  onClick: () => actions.sayHello()
}
```

---

## 🧱 Chapter 3: The Temple of Components

> *"Let thy UI be built from blocks, like holy Lego bricks of type safety and composability."*

* Reuse is divine.
* Every component must start with `_doctype component`.
* Props must be respected as the interface of the sacred.

```jml
_doctype component UserCard

Card {
  Text {
    content: props.username
  }
  Text {
    content: props.email
  }
}
```

---

## 🔀 Chapter 4: The Path of Routing

> *"Ye shall not manually configure thy routes. The compiler knows thy ways."*

* Routing shall be derived.
* Paths become components. Components become paths.
* Let there be no `router.js`, only pages.

```jml
_doctype page about

Page {
  title: "About Us"
  Container {
    Text { content: "We exist because JSX is pain." }
  }
}
```

---

## 🧠 Chapter 5: The Wisdom of Store

> *"You needed localStorage. Now it is `store.set()`."*

* The `store` is simple.
* The `store` remembers.
* The `store` observes.

```jml
import store

Input {
  onEnter: (value) => store.set("username", value)
}

Text {
  content: store.get("username")
}
```

---

## ⚙️ Chapter 6: The Rites of Runtime

> *"The browser API shall be unified under one name: `browser`"*

* You shall scroll with `browser.scrollToTop()`
* You shall update the title with `browser.setTitle()`
* You shall query not with `document.querySelector`, but with `browser.query()`

---

## 📦 Chapter 7: The Scroll of Reuse

> *"Thou shalt not copy-paste thy components across projects. Thou shalt `jawt add`."*

* Share JML components through libraries.
* Use `jawt build --as-lib` to sanctify them.
* Let others use them — if you dare.

---

## 🙅 Chapter 8: The Warnings of Temptation

> *"Thou may be tempted to eject the toolchain, to override it, to reach for React."*

* You shall not question the compiler.
* You shall not eject.
* You shall not say “But in Vue, I could just—”.

JML is law. Jawt is order. JSX is noise.

---

## 🧙‍♂️ Final Words: The Path of One

> *"Jawt was never meant to be for everyone. It was forged in solitude, to suit the needs of the one."*

But if you walk the path, if you accept the rules, and if you embrace the declarative truth — you may never write `<div class="container">` again.

*And that, dear pilgrim, is freedom.*
