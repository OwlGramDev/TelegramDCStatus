ARG dname="amd64"
FROM laky64/tdlib:linux-${dname} AS golang
RUN apt install libc++-dev libc++abi-dev -y
RUN git clone https://github.com/OwlGramDev/TelegramDCStatus
RUN cd TelegramDCStatus && export CGO_LDFLAGS="-lstdc++" && go build -o /usr/src/outputs/TgStatus .