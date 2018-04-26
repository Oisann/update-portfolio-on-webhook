# Update Portfolio On Webhook
Downloads and stores the PORTFOLIO file of the repo it recieves a webhook from.


## Environment Variables
- SECRET *(Same as in the webhook settings on Github)*
- DB_FOLDER *(Folder to store the PORTFOLIO files)*

## Docker

`docker run --rm -it -e SECRET=somereallyobscuresecretthatnobodywillguess -e DB_FOLDER=/home/portfolio -v ${PWD}:/home/portfolio -p 80:80 oisann/update-portfolio-on-webhook:latest`

This container updates its internal database when it recieves a webhook. It actually just downloads the raw PORTFOLIO file in the root of the repo on the master branch, and stores it in the /home/portfolio folder. 

![alt text](https://i.imgur.com/QL1fR7K.png "Example Github webhook settings")
