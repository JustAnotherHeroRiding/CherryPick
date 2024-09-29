# CherryPick

## What is this?

Have you ever needed to download a single folder from a github repo, and spent way too time after finding out there is no way to do this on github.com?
I went into a rabbit hole afterwards of trying multiple extensions that no longer work, ones that allow you to download files only(i needed a fonts folder) or stack overflow threads with multiple long winded answers none of which I wanted to follow for something as simple as this.

This is why I made CherryPick, so that you can quicky download the folder you need and get on with your work.

## Prerequisites

All you need is to have Go installed on your machine.

## How to use

Clone the repo by running git clone https://github.com/JustAnotherHeroRiding/CherryPick.git

Afterwards run the following command, with the url being a direct link to the folder you want to download, which in my case was the fonts folder contained in the assets directory.

go run main.go -url="https://github.com/KristijanKocev/stipsa/tree/main/assets/fonts"

After the download is complete, you will find the contents of that folder in the 'cherrypicked' directory, created automatically in the project folder.

## Github Token

If you get the message that you have sent too many requests, or you want to download a folder from a private repo, make sure to add in your github username and token. The way you can do this is to create a .env file or use the provided .env.example by removing .example from the file name.
Afterwards set the two required environmental variables.

## TODO

- accept multiple folders
- add a progress indicator
- accept single files
- accept multiple file
- test with private repos

## Optimization

Currently it feels very slow when the target repo is very big or the download speed is not super fast. I want to make this as fast as possible, so any Go tips or pull requests are appreciated to make this happen.
