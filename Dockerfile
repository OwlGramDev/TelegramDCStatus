ARG dname="amd64"
FROM laky64/tdlib:linux-${dname} AS golang
COPY TelegramDCStatus/* .
RUN cd TelegramDCStatus && go build -o /usr/src/outputs/TgStatus .
WORKDIR /usr/src/file_manager
COPY linux_mount.sh /usr/src/file_manager
RUN chmod +x /usr/src/file_manager/linux_mount.sh
VOLUME ['/usr/src/file_manager', '/usr/src/outputs']