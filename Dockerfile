FROM gcr.io/distroless/static:nonroot

COPY meta-ads-mcp-server /meta-ads-mcp-server

ENTRYPOINT ["/meta-ads-mcp-server"]
