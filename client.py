from gpiozero import Button
from time import sleep, time
import os
import requests
from io import BytesIO
from pydub import AudioSegment
from pydub.playback import play

def playBytes(data):
    play(AudioSegment.from_file(BytesIO(data), format="mp3"))

def fetchInsult():
    headers = {'Authorization': 'bearer {}'.format(os.environ['GCLOUD_BEARER_TOKEN'])}
    url = 'https://australia-southeast1-upheld-garage-290112.cloudfunctions.net/imtoohappy'
    return requests.get(url, headers=headers).content

button = Button(2)

while True:
    print('waiting')
    button.wait_for_press()
    print('button pressed, creating insult...')
    start_time = time()
    insult = fetchInsult()
    print('insulting you took ', (time() - start_time))
    playBytes(insult)
    print('done')
