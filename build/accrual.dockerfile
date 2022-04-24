FROM gcr.io/distroless/base-debian11
WORKDIR ./
COPY cmd/accrual/accrual_linux_amd64 /accrual
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/accrual"]