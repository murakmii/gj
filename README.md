# gj(WIP)

`gj` is a **toy** JVM(compatible with Java SE8) implementation by Go.

## Usage

```shell
git clone git@github.com:murakmii/gj.git && cd gj
docker pull amazoncorretto:8

echo 'public class HelloGj {
    public static void main(String[] args) {
        System.out.println("Hello, gj!");
    }   
}' > HelloGj.java

docker run -v $(pwd):/gj -w /gj amazoncorretto:8 javac HelloGj.java
go run cmd/main.go HelloGj.class

go build -o gj cmd/main.go
docker run -v $(pwd):/gj -w /gj amazoncorretto:8 ./gj --config dist/config.json --print --main HelloGj
```
