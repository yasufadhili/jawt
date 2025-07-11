## Project Types

Jawt supports both **Applications** and **Libraries** via `jawt.config.json`.

### Applications

* Contain pages (`_doctype page`)
* Define routing
* Compiled to self-contained SPAs
* Use component libraries or logic libraries

### Libraries

* Export **JML components** or **TypeScript scripts**
* Can be reused in other Jawt projects via `jawt add`
* Defined via `jawt.config.json`:

  ```json
  {
    "type": "library",
    "exports": ["components/button.jml", "scripts/utils.ts"]
  }
  ```
