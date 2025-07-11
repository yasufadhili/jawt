## `.jawt/` Workspace Layout

The `.jawt` directory contains **everything Jawt needs internally**, with no exposure to the user.

```
.jawt/
├── config/                  # Cached config
├── libs/                    # Installed reusable JML libraries
├── node/                    # Embedded npm modules (logic libraries only)
├── output/                  # Build output (production-ready assets)
├── generated/               # Compiler-generated files (routes, manifests)
├── runtime/                 # Jawt runtime API: browser.ts, store.ts, etc.
└── scripts/                 # Shared internal helper code
```
