# vidispine configuration files

In order for vsstoragecheck to work, it needs to access the actual
Vidispine API and for this it needs credentials.

The XML files in this directory configure a group with minimal required
permissions for it to work, and a user that is associated with this group.

In order to use vsstoragecheck, simply ingest the documents to Vidispine using
an admin account and then configure the `VIDISPINE_API_USER` and 
`VIDISPINE_API_PASSWD` parameters appropriately.

You can ingest the documents from the commandline using curl:

```bash
$ curl -D- https://vidispine-server/API/group/StorageCheck -u admin -X PUT -d@storage_check_role.xml --header "Content-Type: application/xml"
Enter host password for user 'admin':
HTTP/2 200 
server: nginx/1.19.1
date: Fri, 12 Mar 2021 11:45:41 GMT
content-type: text/plain
content-length: 0
set-cookie: rememberMe=deleteMe; Path=/; Max-Age=0; Expires=Thu, 11-Mar-2021 11:45:41 GMT
$ curl -D- https://vidispine-server/API/user -u admin -X POST -d@storage_check_user.xml --header "Content-Type: application/xml"
Enter host password for user 'admin':
HTTP/2 200 
server: nginx/1.19.1
date: Fri, 12 Mar 2021 11:47:55 GMT
content-type: text/plain
content-length: 0
set-cookie: rememberMe=deleteMe; Path=/; Max-Age=0; Expires=Thu, 11-Mar-2021 11:47:55 GMT
$ echo -n somepassword > passwd.txt
$ curl -D- -u admin https://vidispine-server/API/user/vidispine-monitor/password?type=raw -X PUT -d@passwd.txt --header "Content-Type: text/plain"
$ rm -f passwd.txt
```

Obviously, you change 'somepassword' to the password you want to set!

Then you use `vidispine-monitor` for the user id and `somepassword` (well your changed
password) for the passwd