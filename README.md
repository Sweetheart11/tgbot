# Telegram news summarizer  ![Build Status](https://img.shields.io/badge/build-passing-green)

## How it works

First, we have to add a source to the list of sources. Then after some fetcher will collect all rss articles from the sources and then notifier will summarize its contents and the message to the telegram channel.

## How to use it 

- `/addsource {name: "source_name", url: "source_url"}` - Add a source 
- `/sources` - List sources
- `/deletesource [name]` - Delete a source
