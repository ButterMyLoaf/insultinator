from pydub import AudioSegment
from pydub.playback import play
import io

def playBytes(data):
    play(AudioSegment.from_file(io.BytesIO(data), format="mp3"))
