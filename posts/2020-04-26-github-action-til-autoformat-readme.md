[TIL Auto-Format REAME Action on Github Marketplace][4]
[cflynn07/til][5]

While doing my daily browsing of Hacker News, I came across [this][1] post
([Hacker News][2]) by Simon Wilson on the merits of writing small, actionable
TILs "Today I Learned's." and how he leverages GitHub Actions to automatically
generate a README file based on the contents of his [TIL
repository (simonw/til)][3].

I like the idea of a TIL repo and using GitHub Actions to automate indexing.
Other people also have had the idea to use GitHub actions to index their TIL
repository READMEs as well, however all the examples I could find used GitHub
Actions to run a script that was included in their repository. This works but
it seemed like a reusable GitHub Action that could be quickly dropped into a
TIL repo would be useful.

I've been using GitHub actions for a few months now and I'm enjoying the
product. Free and easy to use CI/CD platforms that integrate with GitHub event
hooks to run arbitrary code on push events have existed for year. What is nice
about GitHub actions is the tight integration with GitHub and the focus on
encouraging users to create small, reusable discrete "Actions" that can be
dropped into others' workflows.



[1]: https://simonwillison.net/2020/Apr/20/self-rewriting-readme/
[2]: https://news.ycombinator.com/item?id=22920437
[3]: https://github.com/simonw/til
[4]: https://github.com/marketplace/actions/til-auto-format-readme
[5]: https://github.com/cflynn07/til
