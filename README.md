# Pager
> A clean pager

## Why?
GNU less has been around for a long time and does the job

I wanted to look at less features and make an opinionated version of it that could also handle stdin and file arguments at the same time

> Check this out
As of version 668 of less(1) if you pass in both an stdin and a list of file arguments the stdin won't be accessible.
```sh
# doesn't work as expected
echo "some redirected input" | less file.txt another_file.txt
```
You will be able to switch between those two files with :n and :p as in next and previous, but won't be able to reach the buffer with the echo output.
```sh
# works as expected
echo "some redirected input" | pager file.txt another_file.txt
```
This is a minor and not important feature and you could get less to do the same in a number of ways, again, this is just me making a pager to learn gdamore/tcell and try out spf13/pflag.

## Installation
* Download a release
* Install with Go
```sh
go install github.com/0xalby/pager@latest
```
* Compile it from source
```sh
git clone https://github.com/0xalby/pager
cd pager
go mod tidy
make build
```

## Usage
### Help
```sh
Usage of ./bin/pager:
  -n, --numbers      Show line numbers
  -o, --offset int   Initial view offset
  -q, --quit         Quits when EOF is reached
  -r, --relative     Show relative line numbers
  -v, --version      Show version number
```
### Controls
* quit with q, Q or Z
* move around with standard vi keys (h,j,k,l) or arrow keys
* skip a paragraph with ^D and ^U or PGDN and PGUP
* teleport to the top of the file with g and to the bottom with G
* go to the next buffer with n or the previous using either p or b
* use r to redraw the buffer

## Todo
* Fix -o
* Partial reading for faster startup like less does
* Implement regex search with / and the reverse with ?
* Optional syntax highlighting
* Disable Go's garbage collector to optimize
* Try to lower executable size as this is about 3MB and less is about 224KB

## Size
I do understand this really doesn't matter unless you work in an embedded environment.
```
┌───────────────────────────────────────────────────────────────────────┐
│ pager                                                                 │
├─────────┬────────────────────────────────────────┬────────┬───────────┤
│ PERCENT │ NAME                                   │ SIZE   │ TYPE      │
├─────────┼────────────────────────────────────────┼────────┼───────────┤
│ 30.98%  │ runtime                                │ 834 kB │ std       │
│ 18.65%  │ .rodata                                │ 502 kB │ section   │
│ 10.37%  │ github.com/gdamore/tcell/v2            │ 279 kB │ vendor    │
│ 6.14%   │ regexp                                 │ 165 kB │ std       │
│ 3.16%   │ os                                     │ 85 kB  │ std       │
│ 2.97%   │ reflect                                │ 80 kB  │ std       │
│ 2.74%   │ .noptrdata                             │ 74 kB  │ section   │
│ 2.54%   │ time                                   │ 68 kB  │ std       │
│ 2.20%   │ fmt                                    │ 59 kB  │ std       │
│ 1.81%   │                                        │ 49 kB  │ generated │
│ 1.57%   │ syscall                                │ 42 kB  │ std       │
│ 1.47%   │ .data                                  │ 40 kB  │ section   │
│ 1.44%   │ strconv                                │ 39 kB  │ std       │
│ 1.27%   │ internal/poll                          │ 34 kB  │ std       │
│ 1.18%   │ slices                                 │ 32 kB  │ std       │
│ 1.02%   │ github.com/spf13/pflag                 │ 28 kB  │ vendor    │
│ 0.97%   │ sync                                   │ 26 kB  │ std       │
│ 0.93%   │ unicode                                │ 25 kB  │ std       │
│ 0.66%   │ strings                                │ 18 kB  │ std       │
│ .....   │ .......                                │ .. kB  │ ...       │
│ 0.00%   │ pager                                  │ 0 B    │ vendor    │
├─────────┼────────────────────────────────────────┼────────┼───────────┤
│ 100.00% │ Known                                  │ 2.7 MB │           │
│ 100%    │ Total                                  │ 2.7 MB │           │
└─────────┴────────────────────────────────────────┴────────┴───────────┘
```
