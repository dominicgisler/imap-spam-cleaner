# Whitelists

Whitelists are individual lists identified by a name (`whitelist1` and `whitelist2` in the following example) which can be reused in inboxes.
The lists contain regexes for the sender of a mail, if the regex matches, the message will be skipped.

Example:

```yaml
whitelists:                       
  whitelist1:                     
    - ^.* <info@example.com>$
    - ^.* <contact@domain.com>$
  whitelist2:
    - ^.* <.*@example.com>$
    - ^.* <.*@domain.com>$
```
