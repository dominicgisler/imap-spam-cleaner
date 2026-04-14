# IMAP spam cleaner

![logo](docs/assets/icon_128.png)

[![Latest release](https://img.shields.io/github/v/release/dominicgisler/imap-spam-cleaner?style=for-the-badge)](https://github.com/dominicgisler/imap-spam-cleaner/releases)
[![License](https://img.shields.io/github/license/dominicgisler/imap-spam-cleaner?style=for-the-badge)](https://github.com/dominicgisler/imap-spam-cleaner/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/dominicgisler/imap-spam-cleaner?style=for-the-badge)](https://github.com/dominicgisler/imap-spam-cleaner/issues)
[![Contributors](https://img.shields.io/github/contributors/dominicgisler/imap-spam-cleaner?style=for-the-badge)](https://github.com/dominicgisler/imap-spam-cleaner/graphs/contributors)

[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-grey?style=for-the-badge&logo=docker)](https://hub.docker.com/r/dominicgisler/imap-spam-cleaner)
[![Buy me a coffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-grey?style=for-the-badge&logo=ko-fi)](https://ko-fi.com/dominicgisler/tip)

A tool to clean up spam in your imap inbox.

Check the [Documentation](https://dominicgisler.github.io/imap-spam-cleaner) for detailed information.

## Summary

This application loads mails from configured imap inboxes and checks their contents using the defined provider.
Depending on a spam score, the message can be moved to the spam folder, keeping your inbox clean.

## Example

```console
$ docker run -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner:latest
INFO   [2026-02-28T16:53:41Z] Starting imap-spam-cleaner v0.5.3
DEBUG  [2026-02-28T16:53:41Z] Loaded config
INFO   [2026-02-28T16:53:41Z] Scheduling inbox info@example.com (*/5 * * * *)
INFO   [2026-02-28T16:55:00Z] Handling info@example.com
DEBUG  [2026-02-28T16:55:00Z] Available mailboxes:
DEBUG  [2026-02-28T16:55:00Z]   - INBOX
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Drafts
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Sent
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Trash
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Spam
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Spam.Cleaner
DEBUG  [2026-02-28T16:55:00Z] Found 34 messages in inbox
DEBUG  [2026-02-28T16:55:00Z] Found 5 UIDs in timerange
INFO   [2026-02-28T16:55:00Z] Loaded 5 messages
DEBUG  [2026-02-28T16:55:06Z] Spam score of message #478 (Herzlichen Glückwunsch! Ihr Decathlon-Geschenk wartet auf Sie. 🎁): 90/100
DEBUG  [2026-02-28T16:55:12Z] Spam score of message #479 (Leider ist bei der Verarbeitung Ihrer Zahlung ein Problem aufgetreten.): 90/100
DEBUG  [2026-02-28T16:55:18Z] Spam score of message #480 (Das neue Geheimnis gegen Bauchfett!): 92/100
DEBUG  [2026-02-28T16:55:26Z] Spam score of message #481 (Schnell: 1 Million / Lady Million): 80/100
DEBUG  [2026-02-28T16:55:32Z] Spam score of message #483 (Vermögen x4 zu Fest): 85/100
INFO   [2026-02-28T16:55:32Z] Moved 4 messages
```

## Contributors

Feel free to contribute to this project, check the [Contribution Guidelines](CONTRIBUTING.md) for more information.

<!-- readme: contributors -start -->
<table>
	<tbody>
		<tr>
            <td align="center">
                <a href="https://github.com/dominicgisler">
                    <img src="https://avatars.githubusercontent.com/u/13015514?v=4" width="100;" alt="dominicgisler"/>
                    <br />
                    <sub><b>Dominic Gisler</b></sub>
                </a>
            </td>
            <td align="center">
                <a href="https://github.com/nistei">
                    <img src="https://avatars.githubusercontent.com/u/1652722?v=4" width="100;" alt="nistei"/>
                    <br />
                    <sub><b>Niklas Steiner</b></sub>
                </a>
            </td>
            <td align="center">
                <a href="https://github.com/alanpbaldwin">
                    <img src="https://avatars.githubusercontent.com/u/1648277?v=4" width="100;" alt="alanpbaldwin"/>
                    <br />
                    <sub><b>alanpbaldwin</b></sub>
                </a>
            </td>
		</tr>
	<tbody>
</table>
<!-- readme: contributors -end -->
