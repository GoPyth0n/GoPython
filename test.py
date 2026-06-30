def mask(x, y):
    return x & y

def mix(a, b):
    return (a + b) & (a - b)

def combine(a, b, c):
    return mask(mix(a, b), c) + mix(b, c)

def chain(x):
    return combine(x, x + 3, 7)

def deep(x):
    return chain(x) + combine(x, 3, x & 9)

def root(x):
    return deep(x) + chain(x + 1) + mask(x * 2, 15)

r1 = root(4)
r2 = root(10)
r3 = root(7)
r4 = root(12)
r5 = root(20)