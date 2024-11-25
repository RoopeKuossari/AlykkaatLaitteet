import RPi.GPIO as GPIO
import time
import os
import datetime
import requests

# GPIO setup
GPIO.setwarnings(False)
GPIO.setmode(GPIO.BOARD)
GPIO.setup(7, GPIO.IN)

# Ensure the "intrudes" directory exists
os.makedirs("intruders", exist_ok=True)

# Main loop
while True:
	i = GPIO.input(7)
	if i == 0:
		print("No Intruders")
		time.sleep(0.5) # Debounce delay
	elif i == 1:
		print("Intruder Alert")
		# Get current timestamp and format it for the filename
		timestamp = datetime.datetime.now().strftime("%Y-%m-%d-%H-%M-%S")
		# Capture image using fswebcam
		image_path = f"intruders/{timestamp}.jpg"
		os.system(f"fswebcam --no-banner -r 640x480 {image_path}")
		print(f"Image captured and saved as {image_path}")
		time.sleep(2.5) # Delay to avoid repeated triggers.

		url = "http://192.168.1.102:8080/upload"
		with open(image_path, "rb") as image_file:
			files = {"image": image_file}
			try:
				# Send the image using a POST
				response = requests.post(url, files=files)
				print(f"HTTP Status Code: {response.status_code}")
				print(f"Server Response: {response.text}")
			except requests.exceptions.RequestException as e:
				print(f"Error sending data: {e}")


		# Wait before capturing the next image
		print("Waiting for the next intruder...")
		time.sleep(2)
