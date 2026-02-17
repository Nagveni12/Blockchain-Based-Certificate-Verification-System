from PIL import Image, ImageDraw, ImageFont
import os

# Create a simple image
img = Image.new('RGB', (400, 300), color='lightblue')
d = ImageDraw.Draw(img)
d.text((100, 140), "Test Certificate\nJohn Doe", fill='blue')
img.save('/home/shreegowri/certificate-blockchain/test-images/test-cert.png')
print("Image created successfully!")
