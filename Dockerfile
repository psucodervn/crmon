FROM scratch
COPY crmon /
ENTRYPOINT ["/crmon"]
