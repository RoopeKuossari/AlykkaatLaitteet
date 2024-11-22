import RPi.GPIO as GPIO
import time
import os
GPIO.setwarnings(False)
GPIO.setmode(GPIO.BOARD)
GPIO.setup(7,GPIO.IN)
while True:
            i=GPIO.input(7)
            if i==0:
                print("No Intruders")
                time.sleep(0.5)

            elif i==1:
                print("Intruder Alert")
                os.system("fswebcam --no-banner -r 640x480 intruders/%Y-%m-%d-%H-%M-%S.jpg")
                time.sleep(2.5)
		
