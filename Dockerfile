FROM amazoncorretto:8

RUN ln -s /usr/lib/jvm/java-1.8.0-amazon-corretto/jre/lib/tzdb.dat /usr/lib/jvm/java-1.8.0-amazon-corretto/lib/tzdb.dat
RUN ln -s /usr/lib/jvm/java-1.8.0-amazon-corretto/jre/lib/currency.data /usr/lib/jvm/java-1.8.0-amazon-corretto/lib/currency.data
