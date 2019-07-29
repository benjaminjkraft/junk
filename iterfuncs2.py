import itertools


def calling(cb):
    """Iterfunc which returns its input, but calling callback after every item.

    In particular, the returned iterator calls cb() after it has been asked for
    a next element, but before returning it.
    """
    def iterfunc(iterator):
        for item in iterator:
            cb()
            yield item
    return iterfunc


def coalescing(iterfunc):
    """Modifies an iterfunc to group elements output "because of" an input.

    In particular, given an iterfunc which takes some input, and given each
    item of the input, returns zero or more items of output,
    coalescing(iterfunc) returns the same items of output, but grouped by input
    item, such that its output is the same length as the input.  For example:
        coalescing(mapping(f))(iterator) == [[f(x)] for x in iterator]
        coalescing(filtering(pred))(iterator) == [
            [x] if pred(x) else [] for x in iterator]
    """
    def coalesced(iterator):
        buf = [[]]
        def cb():
            buf.append([])

        for item in iterfunc(calling(cb)(iterator)):
            buf[-1].append(item)
            while len(buf) > 1:
                yield buf.pop(0)

        yield from buf
    return coalesced


def interleaving(*iterfuncs):
    """Calls several iterfuncs alternately.

    Given several iterfuncs, interleaving returns an iterfunc which, given an
    iterator, returns an iterator which yields:
        - the outputs of the first iterfunc based on the first input
        - the outputs of the second iterfunc based on the first input
        - etc.
        - the same, repeated for the second input
        - etc.
    """
    n = len(iterfuncs)
    def interleaved(iterator):
        wrapped_iters = [
            coalescing(f)(itercopy)
            for f, itercopy in zip(
                iterfuncs, itertools.tee(iterator, n))]
        while wrapped_iters:
            for i, it in enumerate(wrapped_iters):
                try:
                    yield from next(it)
                except StopIteration:
                    wrapped_iters[i] = None
            wrapped_iters = list(filter(None, wrapped_iters))
    return interleaved


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
