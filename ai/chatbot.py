from chatterbot import ChatBot


chatbot = ChatBot('AnoAI_AIO')

if __name__ == '__main__':
    print(chatbot.get_response({
        'text': 'Jól köszönöm, te hogy vagy?',
    }))
