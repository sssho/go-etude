- `app0`: Standalone app without any other local package
- `app1`: Use package `mylib0` in `lib0`
- `app2`: Use pacakge `mylib1` and it's sub pacakge `mylib1/sublib` in `lib1`
- `lib0`: Simple package `mylib0`
- `lib1`: Package `mylib1` with subpackage `mylib1sub` in `lib1/sublib`

```
./
├── app0/
│  ├── go.mod
│  ├── main.go
│  └── sub/
│     ├── lib.go
│     └── subsub/
│        └── sublib.go
├── app1/
│  ├── go.mod
│  └── main.go
├── app2/
│  ├── go.mod
│  └── main.go
├── lib0/
│  ├── go.mod
│  ├── lib.go
│  └── util.go
├── lib1/
│  ├── go.mod
│  ├── lib.go
│  └── sublib/
│     └── util.go
└── README.md
```
