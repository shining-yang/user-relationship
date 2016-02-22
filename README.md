README
====

## Environments
  * Ubuntu 14.04 LTS i386 VM
  * golang 1.5
  * PostgreSQL 9.3  
__NOTE:__ Do not use the golang-v1.2.1 which come with ubuntu 14.04 distribution,
	because there is a compile issue with gopkg.in/pg.v3. I used v1.5 here.
  
## PostgreSQL
#### Install PostgreSQL
`sudo apt-get install postgresql`
    
#### Set password for user 'postgres' so that we can access the PostgreSQL server later  
`sudo -u postgres psql postgres`  
`ALTER USER postgres WITH PASSWORD '123456'`
    
#### (Optional) Set password for the linux user account 'postgres'
`sudo password postgres`
  
## Building
Untar the source tarball to a directory. I choosed: "~/user-relationship".
  
#### Switch to the top working directory
`cd ~/user-relationship`
  
#### Set GOPATH
`export GOPATH=~/user-relationship`
  
#### Install required packages
`go get github.com/gorilla/mux`  
`go get gopkg.in/pg.v3`
  
#### Import SQL procedure
1. Switch to linux user 'postgres' with the following command, and enter the password if you've set one:  
	`su postgres`
2. Now, your command prompt should be shown as 'postgres'. Enter postgresql command prompt:  
	`psql`
3. Now, your command prompt should be shown as 'postgres=#'. Import the SQL script:  
	`\i init.sql`
4. When succeeds, you may return to your own linux account:  
	`\q`        # quit psql  
	`exit`      # quit linux account 'postgres'
  
#### Compile & Run
Issue the following commands to compile the project and run the server:  
	`cd ~/user-relationship/bin`  
	`go build app`  
	`./app`
    
## Testing
Frisby is quite suitable for RESTful API test. See [frisby](http://frisbyjs.com/) for details.
	In this project, for simplicity, very few test cases are provided.
  
1. Install node.js & npm  
	Visit [nodejs](https://nodejs.org) and download a suitable package to install.
2. Install frisby  
	`cd test/`  
	`npm install -g frisby`
3. Install jasmine-node  
	`sudo npm install -g jasmine-node`
4. Run test cases scripts under directory 'test/'  
	`jasmine-node ./`

## Limitations
* For simplicity, I used 'Bigint' rather than 'string' for user-id field

