def a():
    def b():
        return 1
    return b()

x = a()