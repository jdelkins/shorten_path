# shorten_path: A fairly flexible path shortening utility

You can use little tool to create uniquely shortened Unix pathnames for
presentation purposes in a number of user interface contexts. I use it:

- In tmux, as a component of the displayed tab name
- In zsh, in a [powerlevel10k][] prompt element
- In (neo)vim as a buffer and window title component

Being written in go, it is highly portable, statically linked (ideal for
a frequently run program) and solidly performant. The slight downside of using
go is that the resident set size is pretty large. The tradeoff is fine for me.

## Installation

    go install github.com/jdelkins/shorten_path@latest

The repo also contains a Makefile, but there's no need to use it for normal
installation. It is just for creating binary artifacts for various platforms in
the releases. This tiny project isn't worth setting up CI.

## Usage

```
Usage: shorten_path [-?] [-h value] [-H value] [-i value] [-l value] [-m value] [-o value] [-t value] [-T value] [parameters ...]
 -?, --help  display help
 -h, --head-lead-in=value
             character sequence to begin the first element
 -H, --head-lead-out=value
             character sequence to end the first element
 -i, --short-lead-in=value
             character sequence to begin abbreviated elements
 -l, --length=value
             length of path above which shortening will be attempted [1]
 -m, --minimum-element-savings=value
             don't abbreviate a path element unless doing so will result
             in at least this many charcaters saved; use to compensate
             for printable lead-in or lead-out strings, if any [1]
 -o, --short-lead-out=value
             character sequence to end abbreviated elements
 -t, --tail-lead-in=value
             character sequence to begin the last element
 -T, --tail-lead-out=value
             character sequence to end the last element
```

The "lead-in" and "lead-out" options are the main distinct features of this
tool, otherwise it might as well be implemented a shell function. They can be
any string or binary data. See some examples in the section below.

The first and last elements are never shortened by this tool under any
circumstances, although their output may be modified by the optional `-h`,
`-H`, `-t`, and `-T` options. These work just like `-i` and `-o` but are
applied respectively to the first and last path elements, in case you want
their formatting to be distinct.

The `--length` option will forego any shortening (simply outputting the command
line argument) unless the unshortened path is longer than this integer.

The `--minimum-element-savings` option won't shorten a path element unless the
shortened version (excluding any lead-in or lead-out strings) is shorter than
the unshortened path element by at least this number of characters. This option
is handy when the lead-in and -out include printable text. This utility doesn't
know or care whether the lead-in and -out will be printable, or rather some
zero-width formatting code. By default, it is assumed that these are zero
width, unless you, the caller, dictates a different economic with this option.

Subject to the above, each path element is shortened, at the suffix end, to the
fewest number of characters so as to remain uniquely identifiable by inspection
in the filesystem. It is hoped that, for example with the aid of
tab-completion, that one could navigate to the same place using `cd` on a Unix
shell, provided only with the visual output of this tool.

The tool should work fine with "fake" or non-existent paths, although the
uniqueness test is meaningless (every non-existent path is unique relative to
all of the existing paths) and all of the non-existent elements
will be shortened to the minimum length.

## Examples

### Simple usage with printable lead-ins

    $ pwd
    /sys/devices/system/cpu/cpu0
    $ shorten_path -i '[' -o ']' `pwd`
    /sys/[devi]/[sy]/[cp]/cpu0
    $ shorten_path -i '[' -o ']' -m 3 `pwd`
    /sys/[devi]/[sy]/cpu/cpu0

Here, we are using square brackets as lead-in and -out characters. In the
second invocation, which is 1 character shorter overall, we tell `shorten_path`
not to shorten any path element unless it would save at least 3 characters.
This option kicks in on the `cpu` directory in the example (`"[cp]"`, the shortest
unique prefix decorated with brackets, is longer than the original `"cpu"`, so
don't bother shortening. See?). 

### Using elipsis

    $ pwd
    /sys/devices/system/cpu/cpu0
    $ shorten_path -o '…' `pwd`
    /sys/devi…/sy…/cp…/cpu0
    $ shorten_path -o '…' -m 2 `pwd`
    /sys/devi…/sy…/cpu/cpu0

Same thing as above, but using elipsis characters. Note again that, with `-m
2`, we don't bother shortening the `cpu` component into `cp…`. Although the
resulting shortened paths are the same visual length when rendered in
a monospace font, the first is less legible.

### [powerlevel10k][] path component

Put the following function into `.zshrc`, and change `dir` to `shortdir` in
your `.p10k.zsh` file.

`~/.zshrc`
```
function prompt_shortdir() {
    local hi="%F{$POWERLEVEL9K_DIR_ANCHOR_FOREGROUND}"
    local ho="%F{$POWERLEVEL9K_DIR_FOREGROUND}"
    local li="%F{$POWERLEVEL9K_DIR_SHORTENED_FOREGROUND}"
    local lo="%F{$POWERLEVEL9K_DIR_FOREGROUND}"
    local ti="%F{$POWERLEVEL9K_DIR_WORK_FOREGROUND}"
    local to="%F{$POWERLEVEL9K_DIR_FOREGROUND}"
    local shorty=$(shorten_path -i $li -o $lo -h $hi -H $ho -t $ti -T $to $PWD)
    p10k segment -t "%F{$POWERLEVEL9K_DIR_FOREGROUND}$shorty"
}
```

`~/.p10k.zsh`
```
...
typeset -g POWERLEVEL9K_LEFT_PROMPT_ELEMENTS=(
  # =========================[ Line #1 ]=========================
  os_icon                 # os identifier
  shortdir                     # current directory <--**CHANGE HERE**
  vcs                     # git status
  context                 # user@hostname
  # =========================[ Line #2 ]=========================
  newline                 # \n
  prompt_char             # prompt symbol
)
...
```

[powerlevel10k]: https://github.com/romkatv/powerlevel10k
