FROM    golang:1.12

# RUN     apt-get update && \
#     apt-get install -y \
#     postgresql-client \
#     redis-tools

ENV     GO111MODULE=on
ENV     GOFLAGS=-mod=vendor

WORKDIR /go/src/github.com/RoryShively/api

COPY    . .
RUN     go install .

EXPOSE  3100
CMD     ["api"]
