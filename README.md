# gojiai(WIP)

`gojiai` is a **toy** JVM(compatible with Java SE8) implementation by Go.

## Usage

```shell
# Build gj
git clone git@github.com:murakmii/gojiai.git && cd gojiai
go build -o gojiai cmd/main.go

# Compile sample code.
docker pull amazoncorretto:8

echo 'public class HelloGojiai {
    public static void main(String[] args) {
        System.out.println("Hello, gojiai!");
    }   
}' > HelloGojiai.java

docker run -v $(pwd):/gojiai -w /gojiai amazoncorretto:8 javac HelloGojiai.java

# Run it
docker run -v $(pwd):/gojiai -w /gojiai amazoncorretto:8 ./gojiai --config dist/config.json --main HelloGojiai
-> VM initialized!(59 ms)
-> Loaded classes: 165
-> Execute main method...
--------------------------------------
Hello, gojiai!
```
