FROM golang:1.19 AS build-env

WORKDIR /build
COPY . .

RUN make build

FROM gcr.io/distroless/base
COPY --from=build-env /build/bin ./
COPY --from=build-env /build/docs ./docs

ENTRYPOINT ["./runtime"]
CMD ["server"]