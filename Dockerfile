FROM golang:1.14 as build

RUN apt-get update && apt-get install -y ninja-build

# TODO: Змініть на власну реалізацію системи збірки
RUN go get -u github.com/SunRiseGG/ArchitectureLab2/build/cmd/bood

WORKDIR /go/src/ArchitectureLab3/
COPY . .

RUN CGO_ENABLED=0 bood

# ==== Final image ====
FROM alpine:3.11
WORKDIR /opt/ArchitectureLab3/
COPY entry.sh ./
COPY --from=build /go/src/ArchitectureLab3/out/bin/* ./
ENTRYPOINT ["/opt/ArchitectureLab3/entry.sh"]
CMD ["server"]
