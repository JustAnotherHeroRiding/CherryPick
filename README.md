 # CherryPick
<div justify="center" align="center">
<img alt="Cherries" height="280" src="logo.png" />
</div>

## What is this?

Have you ever needed to download a single folder from a GitHub repository, only to find out there's no straightforward way to do it on GitHub.com? I found myself going down a rabbit hole trying multiple extensions that either didn't work, only allowed file downloads (when I needed a fonts folder), or led to convoluted Stack Overflow threads with long-winded answers that I didn't want to follow for something as simple as this.

This is why I created **CherryPick**: to allow you to quickly download the folder you need and get on with your work without the hassle.

## Prerequisites

To use CherryPick, you need to have **Go** installed on your machine. You can download and install Go from the official Go website ([https://golang.org/dl/](https://golang.org/dl/)).

## How to Use

IMPORTANT: The example repo used below is private and will not work for you. Make sure to replace the url with a repo you have access to or a public one.

1. Clone the repository by running:
    
    `git clone https://github.com/JustAnotherHeroRiding/CherryPick.git`
    
2. Navigate to the project directory:
    
   `cd CherryPick`
    
3. Run the following command, replacing the URL with a direct link to the folder you want to download. For example, to download the fonts folder contained in the assets directory, you would run:
    
    `go run main.go https://github.com/KristijanKocev/stipsa/tree/main/assets/fonts`
    
4. Once the download is complete, you will find the contents of that folder in a directory named **`cherrypicked`**, which is created automatically in the project folder.

#### Install using Go

1. Install the repository
    `go install github.com/JustAnotherHeroRiding/cherrypick`

2.Set your environmental variables
    `export GITHUB_USERNAME=YourUsername`
    `export GITHUB_TOKEN=YourGithubToken` 
    `export CHERRYPICK_DOWNLOAD_DIR=YOUR_DOWNLOAD_DIR`
    
3. Download your chosen directory
    `cherrypick https://github.com/KristijanKocev/stipsa/tree/main/assets/fonts`
   
This will create a new directory called cherrypicked with your target directories and/or files.
It will be located in the same directory where you ran the command.

NOTE: You must add the Go binary to your PATH in your shell config for cherrypick to work.

### Downloading multiple directories or files

Provide a comma separated list of directory or file URLs and all of them will be downloaded into the target directory.
IMPORTANT: Do not add spaces after the comma

    go run main.go https://github.com/KristijanKocev/stipsa/tree/main/assets/fonts,https://github.com/KristijanKocev/stipsa/tree/main/assets/images
    

## GitHub Token

If you receive an error message indicating that you have sent too many requests, or if you want to download a folder from a private repository, you will need to provide your GitHub username and a personal access token.

To do this:

1. Create a .env file in the project directory or use the provided .env.example file by renaming it (remove .example from the file name).
    
2. Set the following two required environment variables in your .env file:
    
    `GITHUB_USERNAME=your_username GITHUB_TOKEN=your_token`

### Changing download target directory

To change the default location, change set the environmental variable to your desired location.
For example:
    `export CHERRYPICK_DOWNLOAD_DIR=/desired/download/path`


## TODO

- [x] Accept multiple folders
- [x] Accept single files
- [x] Accept multiple files
- [x] Private repositories(with a token that can access them)

## Contribution

If you'd like to contribute to CherryPick, feel free to open an issue or submit a pull request. Any suggestions for improvements, bug fixes, or new features are welcome!

## License

This project is licensed under the MIT License - see the LICENSE file for details.
