FROM golang:1.15 AS builder

WORKDIR /temp
COPY . .
RUN CGO_ENABLED=0 go build -o ./app .

FROM alpine:3.13.4 as base
RUN apk add --no-cache bash && \
    echo 'export PS1="[docker \u@\h:\[\e[34;1m\]\w\[\033[m\] \[\033[1m\]\[\033[m\]] # "' >> /root/.bashrc && \
    echo 'alias ls="ls --color=auto"' >> /root/.bashrc && \
    echo 'alias ll="ls -l"' >> /root/.bashrc && \
	rm -rf /var/lib/apt/lists/* /var/cache/apk/*
COPY --from=builder /temp/app /usr/local/bin/app

ENTRYPOINT ["app"]