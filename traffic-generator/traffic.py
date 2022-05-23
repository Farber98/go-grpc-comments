import json

from random import randrange
from sys import getsizeof
from locust import HttpUser, between, task

debug = False

# Start: locust -f traffic.py


def printDebug(msg):
    if debug:
        print(msg)


class Reader():

    def __init__(self):
        self.array = []

    def pickRandom(self):
        lenght = len(self.array)

        if (lenght > 0):
            random_index = randrange(0, lenght - 1) if lenght > 1 else 0
            return self.array.pop(random_index)
        else:
            print(">> Reader: No values detected in file.")
            return None

    def load(self):
        print(">> Reader: Opening file...")
        try:
            with open("traffic.json", 'r') as data_file:
                self.array = json.loads(data_file.read())
        except Exception as error:
            print(f'>> Reader: No data loaded, error: {error}')


class MessageTraffic(HttpUser):
    wait_time = between(0.1, 0.9)
    reader = Reader()
    reader.load()

    def on_start(self):
        print(">> MessageTraffic: Starting traffic")

    @task
    def PostMessage(self):
        random_data = self.reader.pickRandom()

        if (random_data is not None):
            data_to_send = json.dumps(random_data)
            printDebug(data_to_send)

            self.client.post("/", json=random_data)
        else:
            print(">> MessageTraffic: Ending traffic.")
            self.stop(True)

    @task
    def GetMessages(self):
        self.client.get("/")
