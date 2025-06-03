# `CheckIt`

[![Code Quality](https://img.shields.io/badge/Code_Quality-A+-brightgreen?style=for-the-badge&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/lucasepe/checkit)



> Render your checklists into clean printable PDF documents – straight from the terminal.

Ideal for quick task lists, audit checklists, packing lists, etc.


## 🔧 Usage

```sh
checkit [OPTIONS] [INPUT_FILE]
```

If no `INPUT_FILE` is provided, input is read from `stdin`.


### Options

```
  -o, --output       Output directory for generated PDF (default: current directory)
  -s, --font-size    Base font size in points (default: 12)
  -v, --version      Show version and exit
  -h, --help         Show help and exit
```

### Examples

```sh
# Generate a PDF checklist from a file
checkit /path/to/my-checklist.md

# Pipe input from another command
cat /path/to/my-checklist.md | checkit
```


### Input Format

`checkit` expects a markdown `.md`-style structured file. Example:

```md
# My CheckList Title

## Morning Routine

- Coffee
- Shower
- Emails

## Work

- Code Review
- Team Meeting
```

You can find some sample input files along with the generated PDF documents [here](./testdata).

## 👍 Support

All tools are completely free to use, with every feature fully unlocked and accessible.

If you find one or more of these tool helpful, please consider supporting its development with a donation.

Your contribution, no matter the amount, helps cover the time and effort dedicated to creating and maintaining these tools, ensuring they remain free and receive continuous improvements.

Every bit of support makes a meaningful difference and allows me to focus on building more tools that solve real-world challenges.

Thank you for your generosity and for being part of this journey!

[![Donate with PayPal](https://img.shields.io/badge/💸-Tip%20me%20on%20PayPal-0070ba?style=for-the-badge&logo=paypal&logoColor=white)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=FV575PVWGXZBY&source=url)


## 🛠️ How To Install

### Download the latest binaries from the [releases page](https://github.com/lucasepe/checkit/releases/latest):

- [macOS](https://github.com/lucasepe/checkit/releases/latest)
- [Windows](https://github.com/lucasepe/checkit/releases/latest)
- [Linux (arm64)](https://github.com/lucasepe/checkit/releases/latest)
- [Linux (amd64)](https://github.com/lucasepe/checkit/releases/latest)

### Using a Package Manager

» macOS » [Homebrew](https://brew.sh/)

```sh
brew tap lucasepe/cli-tools
brew install checkit
```