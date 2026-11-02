# FormCodeChallenge


## Description
A form builder that lets users create forms
Users can fill out forms, edit them and see the results

Forms are stored in local memory  

## Usage Instructions
in the project folder, 
`run go run .`
The application will run on localhost:4000


## API Call
Get a random fact from https://uselessfacts.jsph.pl/api/v2/facts/random?language=en
It returns in a JSON file like this
{"id":"b8e2f218a511aba6c467a63df7b88842","text":"Each year in America there are about 300,000 deaths that can be attributed to obesity.","source":"djtech.net","source_url":"https://www.djtech.net/humor/shorty_useless_facts.htm","language":"en","permalink":"https://uselessfacts.jsph.pl/api/v2/facts/b8e2f218a511aba6c467a63df7b88842"}

the TemplateQuestion Struct maps the text sections to a string.
Then append that string to the form questions. 
Repeats 3 times to get three randomly generated facts. (Ya I know facts aren't questions so it doesn't make any sense.)

## If I had more time...
- Multiple choice questions
- Better state management
- Refresh and rebuild handling. 
- Question ordering is kind of buggy. 


## AI Usage
Things AI implemented 
- All styling
- All of the html structures. 
- API call error handling
- Most of the question ordering. 