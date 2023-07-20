# GophKeeper

Application made while learning Go. The main aim is to provide secure storage option for clients.

## Services

Primarily there are two services required to run this app:

- Server - can be either remote or local
- Client - CLI app for interacting with server storage


## Features
### Server can store varying data

- Login+password pair
- Text data;
- Files (binary data);
- Credit card info. 
Each data supports storing text metadata 

### Global features
- [ ] CRUD of private data
- [ ] Account management
  - [ ] Registration
  - [ ] Logging in
- [ ] Data sync between devices
- [ ] Server data storage
- [ ] Redis for session storage
- [ ] Postgres for private data
- [ ] File storage for large files
- [ ] Data encoding on client
- [ ] gRPC server-client connection


### Also:
- [ ] Client is a CLI-application for Windows, Linux and Mac OS
  - [ ] TUI is planned for future
- [ ] Client provides build version and date



