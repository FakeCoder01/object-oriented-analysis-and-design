FROM alpine:3.18 AS builder
RUN apk add --no-cache build-base cmake wget linux-headers

WORKDIR /app
COPY . .

RUN wget https://raw.githubusercontent.com/yhirose/cpp-httplib/master/httplib.h -O httplib.h
RUN wget https://raw.githubusercontent.com/nlohmann/json/develop/single_include/nlohmann/json.hpp -O json.hpp

RUN mkdir build && cd build && cmake .. && make

FROM alpine:3.18
RUN apk add --no-cache libstdc++
WORKDIR /app
COPY --from=builder /app/build/backend .
EXPOSE 8080
CMD ["./backend"]
