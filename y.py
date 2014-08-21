def g(x):
    return x(x)

class Foo:
    __metaclass__ = Bar(Bar)
