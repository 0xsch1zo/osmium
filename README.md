# Osmium
It's just a command and controll server. It's awfull. Don't use it.

## Build agent
On windows there will probably be a project file created for visual studio
```sh
cmake --build build
```
Or if you are crosscompiling then:
```sh
cmake --build build -DCMAKE_TOOLCHAIN_FILE=../crosscompile.cmake
```
And then on linux:
```
cd build
make -j
```

## Build teamserver
```sh
go build -o teamserver ./cmd/main.go
```
