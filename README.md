# shorten_path: A fairly flexible path shortening utility

You can use little tool to create uniquely shortened Unix pathnames for presentation purposes in
a number of user interface contexts. I use it:

- In tmux, as a component of the displayed tab name
- In zsh, in a [powerlevel10k][] prompt element
- In (neo)vim as a buffer and window title component

Being written in go, it is highly portable, statically linked (ideal for a frequently run program)
and solidly performant. The slight downside of using go is that the resident set size is pretty
large. The tradeoff is fine for me.

## Installation

    go install github.com/jdelkins/shorten_path@latest

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

The "lead-in" and "lead-out" options are the main distinct features of this tool, otherwise it might
as well be implemented a shell function. They can be any string or binary data. See some examples
below.

The `--length` option will forego any shortening (simply outputting the command line argument)
unless the unshortened path is longer than this integer.

The `--minimum-element-savings` option won't shorten a path element unless the shortened path
element (excluding any lead-in or lead-out strings) is shorter than the unshortened path element by
at least this number of characters. This option is handy when the lead-in and -out include printable
text. This utility doesn't know or care whether the lead-in and -out will be printable, or rather
some zero-width formatting code, and we don't care. We just assume they are zero width, unless you,
the caller, dictates a different economic with this option.

## Examples

### Simple usage with printable lead-ins

    $ pwd
    /sys/devices/system/cpu/cpu0
    $ shorten_path -i '[' -o ']' -m 2 `pwd`
    /sys/[devi]/[sy]/cpu/cpu0

Here, we are using square brackets as lead-in and -out characters. We tell shorten_path not to
shorten unless the shortened version is at least 2 characters
($2 = \textrm{len}('[') + \textrm{len}(']')$). This kicks in on the `cpu` directory in the
example.

### [powerlevel10k][] path component

Put the following function into `.zshrc`, and change `dir` to `shortdir` in your `.p10k.zsh` file.

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
  33   ┊ typeset -g POWERLEVEL9K_LEFT_PROMPT_ELEMENTS=(
     1 ┊   # =========================[ Line #1 ]=========================
     2 ┊   os_icon                 # os identifier
     3 ┊   shortdir                     # current directory <--**CHANGE HERE**
     4 ┊   vcs                     # git status
     5 ┊   context                 # user@hostname
     6 ┊   # =========================[ Line #2 ]=========================
     7 ┊   newline                 # \n
     8 ┊   prompt_char             # prompt symbol
     9 ┊ )
...
```

[powerlevel10k]: https://github.com/romkatv/powerlevel10k
