# fediclient

Do you ever wonder what's behind that `@username@social.example.net` username that you keep seeing mentioned?

Chances are it's probably a Fediverse username, and the way you turn that into a profile and list of posts is complicated

1. strip the initial `@` and then `curl https://social.example.net/.well-known/webfinger?resource=acct:username@social.example.net`
2. look for the rel=self link
3. curl that url with `curl -H "Accept: application/activity+json" <url>`
4. look at the summary
5. look for the outbox URL
6. curl the outbox url with `curl -H "Accept: application/activity+json" <url>`
7. look for the `last` URL
8. fix the character encoding, and then curl the `last` URL with `curl -H "Accept: application/activity+json" <url>`

Or you could just run this CLI app that will do all of that for you

```shell
fediclient @username@social.example.net
```

Different implementations return different results, and I may have misread the docs, so this may not work for all instances. It has been tested against both Mastodon and Hometown

PRs welcome to add support for activitypub servers not currently working
