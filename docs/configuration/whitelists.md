# Whitelists

Whitelists are named lists (`whitelist1` and `whitelist2` in the example below) which can be reused in inboxes.
Each entry is a [regular expression](https://pkg.go.dev/regexp/syntax). If one of them matches the sender, the message is skipped and not analyzed.

The expression is matched against the raw `From` header, e.g. `Example Info <info@example.com>`, not just the address:

- It matches anywhere in the header unless it is anchored with `^` and `$`. `info@example\.com` alone also matches `info@example.com.evil.org`.
- Escape dots as `\.`, otherwise `.` matches any character.
- Some senders use a bare address without a name (`info@example.com`). `^.* <info@example\.com>$` doesn't match those, `^(.* )?<?info@example\.com>?$` matches both forms.

**Warning:** anyone can put any address into the `From` header. Every message claiming to come from a whitelisted sender skips the spam analysis, so keep the lists short.

Example:

```yaml
whitelists:
  whitelist1:                          # exact addresses
    - ^.* <info@example\.com>$
    - ^(.* )?<?contact@domain\.com>?$  # also matches a bare address
  whitelist2:                          # whole domains
    - ^.* <.*@example\.com>$
    - ^(.* )?<?[^<>@ ]+@domain\.com>?$ # also matches a bare address
```
