def caller(cb):
    while True:
        try:
            yield cb()
        except StopIteration:
            break


sentinel = object()


def interleaving(*iterfuncs):
    def iterfunc(iterator):
        outbuffers = [[[]] for _ in iterfuncs]
        inbuffers = [[] for _ in iterfuncs]
        done = 0

        def cb(outb, inb):
            def callback():
                nonlocal done
                outb.append([])
                if not inb:
                    if done:
                        raise StopIteration
                    try:
                        item = next(iterator)
                    except StopIteration:
                        done = 1
                        raise
                    for b in inbuffers:
                        b.append(item)
                return inb.pop(0)
            return callback

        iters = [f(caller(cb(outb, inb)))
                 for outb, inb, f in zip(outbuffers, inbuffers, iterfuncs)]
        while done < 2:
            if done:
                done += 1
            for b, it in zip(outbuffers, iters):
                if b and b[0] is None:  # this iter is exhausted
                    continue
                while len(b) < 2:
                    try:
                        elem = next(it)
                    except StopIteration:
                        b.append(None)
                    else:
                        b[-1].append(elem)
                yield from b.pop(0)

    return iterfunc


def mapping(f):
    def iterfunc(iterator):
        for i in iterator:
            yield f(i)
    return iterfunc


def filtering(pred):
    def iterfunc(iterator):
        for i in iterator:
            if pred(i):
                yield i
    return iterfunc


def main():
    expected = [0, 0, -1, -2, 2, -3, -4, 4, -5, -6, 6, -7, -8, 8, -9]
    actual = list(interleaving(
        mapping(lambda x: -x),
        filtering(lambda x: x % 2 == 0)
    )(iter(range(10))))
    if actual == expected:
        print("ok")
    else:
        print("not ok: %s" % actual)


if __name__ == '__main__':
    main()
