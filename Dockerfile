FROM golang:1.17 AS builder

ENV PKG=/go/src/github.com/olavtar/chromedp
WORKDIR ${PKG}

# compile test binary
COPY . .
RUN make

FROM chromedp/headless-shell:latest

COPY --from=builder /go/src/github.com/olavtar/chromedp/chromedp chromedp

ENTRYPOINT [ "/chromedp" ]

