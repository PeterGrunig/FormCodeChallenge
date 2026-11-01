# FormCodeChallenge

## API Call
Get a random fact from https://uselessfacts.jsph.pl/api/v2/facts/random?language=en
It returns in a JSON file like this
{"id":"b8e2f218a511aba6c467a63df7b88842","text":"Each year in America there are about 300,000 deaths that can be attributed to obesity.","source":"djtech.net","source_url":"https://www.djtech.net/humor/shorty_useless_facts.htm","language":"en","permalink":"https://uselessfacts.jsph.pl/api/v2/facts/b8e2f218a511aba6c467a63df7b88842"}

the TemplateQuestionStruct maps the text sections to a string.
Then append that string to the form questions. 