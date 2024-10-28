from chatbot import chatbot
from chatterbot.trainers import ListTrainer
import os
import json

E_BOT = 'bot'
E_PARTNER = 'partner'

PATH = '../data/conversations'
trainer = ListTrainer(chatbot)

if __name__ == '__main__':
    processedIds = []
    conversation = []
    # Open all .json file in a directory
    for file in os.listdir(PATH):
        if file.endswith('.json'):
            conversation = []
            with open(f'{PATH}/{file}', 'r') as f:
                data = json.load(f)

                if data['id'] not in processedIds:
                    processedIds.append(data['id'])
                    joinedMsg = ''
                    last_entity = None

                    for message in data['messages']:
                        entity = message['entity']
                        msg = message['message']

                        if entity == last_entity:
                            joinedMsg += f'\n{msg}'
                        else:
                            if joinedMsg:
                                conversation.append(joinedMsg)
                            joinedMsg = msg
                        last_entity = entity

                    if len(joinedMsg) > 0:
                        conversation.append(joinedMsg)



            trainer.train(conversation)

    print('Training complete ' + str(len(processedIds)) + ' files processed')