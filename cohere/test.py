import cohere
co = cohere.Client("vlXwe8qP8xFvknbHOKGHdVT755uXaePhu9QEUWfr")

message = "Please generate a brief news-style article (500-1000 words) regarding one of the following topics: breaking news, business, international affairs, arts and culture, sports,  politics, human interest, science and technology. You can use real-world topics as a premise, but please make up all of the details in the article. Do not generate an article that is very similar to one that has already been mentioned in this conversation."
file = open("output.txt", "a")
history = []
for i in range(500):
    response = co.chat(message=message, chat_history=history)
    file.write(response.text)
    file.write("\\\n")
    history.append({"role": "USER", "message": message})
    history.append({"role": "CHATBOT", "message": response.text})


file.close()