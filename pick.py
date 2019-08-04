from typing import TypeVar

T = TypeVar('T', bound='Mapping[str, object]')
T = TypeVar('K', bound='str')

def pick(o: T, keys: K) -> :
    return {k: v for k, v in keys.iteritems() if k in keys}
