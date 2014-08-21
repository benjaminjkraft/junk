class FrozenObjectError(Exception):
    pass

def fail(*args, **kwargs):
    raise FrozenObjectError

def freeze(instance):
    typ = type(instance)
    newdict = typ.__dict__.copy()
    newdict['__setattr__'] = fail
    newdict['__delattr__'] = fail
    cls = type(typ.__name__, (), newdict)
    newinstance = object.__new__(cls)
    newinstance.__dict__ = instance.__dict__.copy()
    return newinstance

    
class FrozenObject(object):
    def __new__(cls, instance):
        inst = object.__new__(cls)
        object.__setattr__(inst, '_instance', instance)
        return inst

    def __getattribute__(self, name):
        if name == '_instance':
            raise AttributeError
        if name in ('__class__', '__dir__'):
            return object.__getattribute__(self, name)
        print "getting", name
        instance = object.__getattribute__(self, '_instance')
        print "from", instance
        return getattr(instance, name)

    def __str__(self):
        base = object.__getattribute__(self, '_instance').__str__()
        return "FrozenObject(" + base + ")"

    def __repr__(self):
        base = object.__getattribute__(self, '_instance').__repr__()
        return "FrozenObject(" + base + ")"

    def __unicode__(self):
        base = object.__getattribute__(self, '_instance').__unicode__()
        return "FrozenObject(" + base + ")"

    __setattr__ = fail
    __delattr__ = fail


class Foo(object):
    def __init__(self):
        self.a = 1

    def b(self):
        return 2

    @staticmethod
    def c():
        return 3

    def e(self):
        self.f = 1
