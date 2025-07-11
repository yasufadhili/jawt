## Runtime Architecture

* **Pages**: Compiled into custom elements like `<page-home>`, injected into `<router-view>`
* **Routing**: SPA model, powered by Vaadin Router, dynamic lazy-loading via import()
* **State**: Managed via NanoStores, scoped per feature
* **Runtime APIs**:

  * `browser`: DOM & BOM helpers
  * `store`: Unified client-side persistence
  * `clipboard`, `env`, `events`, `network`, `date`, etc.
* **Logic & Scripts**: Written in `.ts`, imported into components

---

## Built-in Component Library (Lit-based)

Included out of the box:

* `Button`, `Input`, `Text`, `Card`, `Container`, `Grid`, `List`, `Modal`, `Form`, `Header`, `Footer`, `RouterView`

These can be extended, overridden, or ignored in favour of user libraries.

---

## Built-in Runtime APIs (TypeScript)

These modules are embedded into Jawt and available from any `.jml` file:

| Module | Functions |
|---|---|
| `browser` | `query()`, `setTitle()`, `scrollToTop()` |
| `store` | `get()`, `set()`, `clear()`, `observe()` |
| `events` | `emit()`, `on()`, `off()` |
| `clipboard` | `copy()`, `paste()` |
| `network` | `fetch()`, `getJSON()`, `post()` |
| `date` | `now()`, `format()`, `parse()` |
| `env` | `isProd()`, `getEnv()` |
