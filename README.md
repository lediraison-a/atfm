This project is not in active development.
I started it to learn goland and because I couldn't find any terminal file manager that suited my needs.
I don't plan on using a terminal file manager anymore, the tmux + command line combo is enough for me.

# ATFM
ATFM (Another Terminal File Manager) is a fully featured and vim inspired terminal file manager, written in go, and packed in a single
executable with no dependencies needed.

## Features
- file exploration
- file modification (rename, copy, new...)
- zip & tar.gz support
- embedded shell
- dual pane, tabs
- sftp support
- extensive customization
- file preview
- xdg trash support

## Installation

### Build it from source
first install the latest version of [golang](https://go.dev/doc/install)

clone this repo
```bash
git clone https://git.alediraison.com/alban/atfm.git
```
build it
```bash
cd atfm
go build .
```

## Usage
```bash
atfm
```

## Contributing

Pull requests are welcome. Please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)

## Acknowledgements
I learned go while working on this project. The code is not commented.

Big thanks to all of these awesome projects I used to build ATFM.
- golang
- afero
- tview
- cobra
- tcell
