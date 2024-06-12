# Email for my Server

A simple email server that I made for use on my personal projects.

## How to use

This is quite a bare bones email server, all it does is send in emails and 
store any emails that it sends.
The email works as a reverse email proxy, so the end user's information is
routed through the server, and the server then sends the email on the user's
behalf.

The email server also stores the user's information in a postgres database.
