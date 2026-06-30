def sign(x):
    if x < 0:
        return -1
    elif x == 0:
        return 0
    else:
        return 1
    
x = sign(-5)
y = sign(5)
z = sign(0)