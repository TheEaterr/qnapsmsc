# syntax=docker/dockerfile:1
# check=skip=SecretsUsedInArgOrEnv
FROM alpine:latest

# Install required packages
RUN apk add --no-cache go gcc musl-dev make

# Set the working directory
WORKDIR /app

EXPOSE $PORT

# Copy the source code
COPY . .

# Build the application
RUN make build

ENV QNAPSMSC_PORT=":9094"
ENV QNAPSMSC_HANDLER="log"
ENV QNAPSMSC_USERNAME="admin"
ENV QNAPSMSC_PASSWORD="placeholder"
ENV QNAPSMSC_MAIL_SENDER="placeholder"
ENV QNAPSMSC_MAIL_RECEIVER="placeholder"
ENV QNAPSMSC_SMTP_USERNAME="placeholder"
ENV QNAPSMSC_SMTP_PASSWORD="placeholder"
ENV QNAPSMSC_SMTP_HOST="localhost"
ENV QNAPSMSC_SMTP_PORT="587" 

# Run the application
CMD ["sh", "-c", "/app/bin/qnapsmsc --port $QNAPSMSC_PORT --handler \"$QNAPSMSC_HANDLER\" --username \"$QNAPSMSC_USERNAME\" --password \"$QNAPSMSC_PASSWORD\" --mail-sender \"$QNAPSMSC_MAIL_SENDER\" --mail-receiver \"$QNAPSMSC_MAIL_RECEIVER\" --smtp-username \"$QNAPSMSC_SMTP_USERNAME\" --smtp-password \"$QNAPSMSC_SMTP_PASSWORD\" --smtp-host \"$QNAPSMSC_SMTP_HOST\" --smtp-port $QNAPSMSC_SMTP_PORT"]
