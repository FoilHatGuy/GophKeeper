FROM golang:1.20.6-bullseye
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go build -o /app/main /app/cmd/server
CMD ["/app/main"]
#CMD ["tail", "-f", "/dev/null"]