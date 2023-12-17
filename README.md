# C103-book

このブランチはコミックマーケットC103で頒布の同人誌『Make JVM』用に、`main`ブランチの内容を添削しコミットしています。

```
$ docker pull amazoncorretto:8
$ docker run -v $(pwd):/MakeJVM -w /MakeJVM amazoncorretto:8 javac MakeJVM.java
go build -o make_jvm main.go
./make_jvm
return: 55
```
