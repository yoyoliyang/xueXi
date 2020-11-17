from PIL import Image

i = Image.OPEN("../sc.png")

colors = sorted(i.getcolors())

print(colors)

