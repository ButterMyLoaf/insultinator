from random import seed, choice, shuffle
import requests

class Insultinator:
    def getInsults(self):
        docs = requests.get('https://docs.googleapis.com/v1/documents/'+self.path)
        self.insults = []
        

    def __init__(self, path, seedNum):
        insults = open(path, "r")
        self.insults = insults.readlines()
        insults.close()
        self.index = [i for i in range(len(self.insults))]
        shuffle(self.insults)
        shuffle(self.index)
        seed(seedNum)
        self.i = 0

    def hitMe(self):
        insult = self.insults[self.index[self.i]]
        self.i += 1
        if self.i == len(self.insults):
            shuffle(self.index)
            self.i = 0
        return insult
