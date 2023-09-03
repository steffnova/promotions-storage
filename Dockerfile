FROM golang:1.20 as builder

WORKDIR /app
COPY . .
RUN make build


FROM golang:1.20-alpine

WORKDIR /app

COPY --from=builder /app/bin/server/server .
COPY --from=builder /app/promotions.csv .

EXPOSE 8080

CMD ["./server"]