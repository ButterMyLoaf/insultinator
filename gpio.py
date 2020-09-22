from gpiozero import Button
from time import sleep, time
from speech import audio
from player import playBytes
from insult import Insultinator

button = Button(2)
mistake = Insultinator("./insults.txt", 1)

while True:
    print('waiting')
    button.wait_for_press()
    print('button pressed, creating insult...')
    start_time = time()
    s = audio(mistake.hitMe())
    print('insulting you took ', (time() - start_time))
    playBytes(s)
    print('done')
    