[![Build](https://github.com/TheEaterr/qnapsmsc/actions/workflows/main.yml/badge.svg)](https://github.com/TheEaterr/qnapsmsc/actions/workflows/main.yml)

# qnapsmsc

`qnapsmsc` is a simple Go program meant to be used as a mock SMSC server with a QNAP NAS to allow for custom handling of these notifications. This project is intended to run on an external server. I mostly created this because I couldn't get my ISP's `smtp` server to work with my NAS, but it could be useful for other monitoring solutions without having to run anything special on the NAS.

## Installation

1. Compile the go project using `make build` or build the docker image using `docker build .`.
2. Run the project using the binary or the docker image with (for example) docker compose. 
3. Log in to the NAS web interface
4. Open Notification Center, then Service Account and Device Pairing.
5. In the SMS tab, click the Add SMSC Service button and:
    1. In SMS service provider, select custom.
    2. Set Alias to qnapsmsc.
    3. Set URL template text to http://{host}:{port}/notification?phone_number=@@PhoneNumber@@&text=@@Text@@&username=@@UserName@@&password=@@Password@@. Customize the port and host.
    4. Confirm the settings.



## Customization

`qnapsmsc` supports the following command line flags (or equivalent environment variable using the docker image):

| Flag              | Env variable           | Default value | Description                                     |
| ----------------- | ---------------------- | ------------- | ----------------------------------------------- |
| `--port`          | QNAPSMSC_PORT          | `:9094`       | Address/port where to serve the metrics.        |
| `--username`      | QNAPSMSC_HANDLER       | `admin`       | Username to connect to the notification server. |
| `--password`      | QNAPSMSC_USERNAME      | `notsecure`   | Password to connect to the notification server. |
| `--handler`       | QNAPSMSC_PASSWORD      | `log`         | Handler to use for notifications (log or mail). |
| `--mail-sender`   | QNAPSMSC_MAIL_SENDER   |               | Email address to use as sender.                 |
| `--mail-receiver` | QNAPSMSC_MAIL_RECEIVER |               | Email address to use as receiver.               |
| `--smtp-host`     | QNAPSMSC_SMTP_USERNAME | `localhost`   | SMTP host to use for sending emails.            |
| `--smtp-port`     | QNAPSMSC_SMTP_PASSWORD | `587`         | SMTP port to use for sending emails.            |
| `--smtp-username` | QNAPSMSC_SMTP_HOST     |               | SMTP username to use for sending emails.        |
| `--smtp-password` | QNAPSMSC_SMTP_PORT     |               | SMTP password to use for sending emails.        |
| `--log-file`      | QNAPSMSC_LOG_FILE      |               | Log file path (defaults to empty, i.e. STDOUT). |