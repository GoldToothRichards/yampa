FROM golang:1.19 as build

# Create yampa user
ENV USER=yampa
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -o yampa-cli

# FROM gcr.io/distroless/static-debian11 as runtime
FROM golang:1.19 as runtime
COPY --from=build /app/yampa-cli yampa-cli
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER yampa:yampa
ENTRYPOINT [ "./yampa-cli" ]
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "./yampa-cli", "healthcheck" ]