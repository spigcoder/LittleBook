FROM ubuntu:latest
COPY webook /app/webook
RUN chmod +x /app/webook  # 添加执行权限
WORKDIR /app
ENTRYPOINT ["/app/webook"]