
# goshop

Prepare go

```shell
GOPATH=~/projects/go
mkdir -p $GOPATH
echo "export GOPATH=$GOPATH" >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
```

goshop download

```shell
go get github.com/rofrol/goshop
cd $GOPATH/src/github.com/rofrol/goshop
```


Download css and js

```shell
wget http://foundation.zurb.com/cdn/releases/foundation-5.0.3.zip
unzip foundation-5.0.3.zip -d static/foundation
rm foundation-5.0.3.zip
```

perm

download openssl http://www.openssl.org/related/binaries.html

generate cert http://stackoverflow.com/questions/10175812/how-to-build-a-self-signed-certificate-with-openssl

```shell
mkdir tls && cd tls
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

Run it

```shell
sh db/init.sh
go get .
goshop
```

Go to https://localhost:9000 and https://localhost:9000/admin (login and password from `db/schema.sql`)

## Go - Managing versions

Q: How about automatic version managment in go?

A: Versioning cannot be correctly done automatically for non trivial cases. The trivial cases are few minutes of manual work.

Q: What is the proper way to manage two different versions of the same package

A: The name of a package is equivalent to the first number of a semantic version -- e.g., a package foo exposes a backward compatible interface always and forever. If a breaking change becomes necessary, the name of the package changes to, e.g., foo2 or something.

Q: where does the convention of package names being the major version come from? That is to use foo, foo2, foo3, etc.

A: It comes from Rob Pike and the core Go development team. See the FAQ http://golang.org/doc/faq#get_version . He doesn't state it like I do though; but, the end result is the same. According to semver, if you make a breaking, incompatible change, you bump the 1st semver number. According to Go community convention, if you make a breaking change, you are advised to call your package a new name. The syntax of the "name" may be different, but the core ideas are the same.

https://plus.google.com/113468371879331813621/posts/P4rcZAsHPTB

## HTTP Redirect

Remember to add return after Redirect, if you want to exit function.

```shell
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
	return
```

http://en.wikipedia.org/wiki/Post/Redirect/Get

http://en.wikipedia.org/wiki/HTTP_303

303 for HTTP 1.1, maybe problem with old corporate proxies, so 302 could be better
http://stackoverflow.com/questions/46582/response-redirect-with-post-instead-of-get

The common practice is to redirect only after successful forms.
So forms with errors are treated by the same POST request, and so have
access to the data.
https://groups.google.com/forum/?fromgroups#!msg/golang-nuts/HeAoybScSTU/qxp1H7mWZVYJ
