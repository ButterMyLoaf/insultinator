# Insultinator

It insults you on demand, what more do you want.

You can see (or even add more) the list of insults here [here](https://docs.google.com/spreadsheets/d/18J1dfIk2ckKd8885XvytVONG1cYu0Bjo_NP69ZmB6co/edit?usp=sharing).

## How to use

dude why

First setup your google cloud stuff, enable billing and the APIs.

1. Cloud functions
2. Cloud build (optional)
3. Text-to-speech
4. Google Sheets

Also get [`gcloud`](https://cloud.google.com/sdk/docs/quickstart) and do the init and stuff.

Create a sheet and get the SheetID (it's in the URL, it's those long string of seemingly random chars). This code here only reads insults from column A, so be sure to add all the insults in column A. It also relies on the cell `F1` to provide the number of insults. Put the query `=COUNTA(A1:A100)` to count all non-empty cells in cell A1 to A100.

In GCP console, create a service account with the role `Browser` and share the email to your designated google sheet. (idk if this step actually helps but I found someone mention it in some stackoverflow post and I'm too afraid to change it)

And then to deploy, you can setup cloud build to use the `cloudbuild.yaml` or do the command

```bash
gcloud functions deploy [FUNCTION NAME HERE] \
    --region=[REGION] \
    --source=. \
    --trigger-http \
    --entry-point=InsultMe \
    --runtime=go113 \
    --set-env-vars=PLEASE_NO_MORE=[INSULT LIMIT HERE],CHAR_LIMIT=[INSULT CHARACTER LIMIT HERE],SHEET_ID=[YOUR SHEET ID HERE] \
    --service-account=[YOUR SERVICE ACCOUNT EMAIL HERE]
```
