# Ephem & WebUI

### _One quick thing..._

 Ephem was written as a weekend project to learn Go for the first time. As with all new learning, it is almost certainly not perfect, and almost certainly does not follow common "go" conventions (or in some cases at all). It may not always be consistent in it's implementation (this was a make-it-work exercise, not an architectural one). I sincerely didn't quite know what I was doing yet. I'm also not a front-end developer, so WebUI may also be a bit clunky.

 _When reviewing this project, please be kind._

## Ephem
**Ephem**, short for "Ephemeral", is an encrypted store for short-term secrets persistence, retrieval, and expiration.

Ephem owns the responsibility of:
* Accepting a secret on the provided URL path
* Encrypting that secret
* Generating a hash of that path
* Persisting the encrypted using the hashed path as the key to locate the secret in the future
* Enforcing maximum read counts and secret lifetimes by purging secrets that exceed their configured limits

Simply stated, Ephem owns _encryption and persistence_. It does not have a user interface and it does not own URL generation.

## WebUI
**WebUI** is a front end user interface for Ephem. WebUI owns the responsibility of:
* Accepting the user's provided secret via HTTP POST
* Generating a random URL that will be provided to Ephem and used as the encryption key for the secret
* Calling Ephem with both the URL and secret
* On success response from Ephem, returning the generated URL to the user
* For retrieval, accepting the previously provided URL from the user
* Requesting the secret from Ephem using the key embedded within that URL
* Returning the decrypted secret returned by Ephem to the user.
  
Simply stated, WebUI owns _user interactivity and URL generation_. It does not own encryption or persistence.

In the current implementation, WebUI and Ephem are integrated in code (vs. simply separate components that interact over a standard protocol). However, with the correct configuration it should be possible to deploy just the GRPC API if desired.

## Security Implementation Note
The actual implementation of the core security aspects such as proper implementation and use of best practice APIs at this point cannot be guaranteed. Although Ephem does encryption and hashing as part of the implementation, perfect implementation of these things was _not_ the primary purpose of this project. The purpose was to learn the fundamentals of writing a functional Go application and a crypto tool was an interesting problem space to play with. Perfectly implementing cryptography within the application was not critical to the primary goal of simply learning Go.

## Objectives
* Keep it simple. Do one thing, do it well, and do no more. Simplicity increases security. Whiz-bang should be avoided unless it provides direct functional value to the platform.
* Maintain separation of layers. Allow each layer to completely own it's functionality and make the implementation of that functionality invisible to the caller.
* Never persist or log identifiable information that may identify or expose a secret. Use hashes and encryption for anything provided by the caller. 
* Security first. Functionality should never compromise security. Encryption should be the first step when accepting a secret and the last step when returning it. Use the best ciphers available with their maximum key length(s). Assume that a given algorythem will eventually be compromisable. Do not trust the layers.
* When possible, attempt to optimize resources (memory, compute, network) through the use of pointers rather than values.
* Do not manage state in the application. The only thing that should be stateful is the persistence store (database). Restarting the application should lose nothing.
* Avoid External Dependencies - The ability to develop, start, or run should not require external resources to be consistently available except when absolutely necessary (such as a database or SSO provider). If a resource is required, embed it when possible/practical, do not reference it.
* Platform Independence - Do not rely on the functionality of a specific platform or operating system. Linux and Windows should both be equal in their functionality and capabilities. 

## Next Steps
Bug fixes, optimizations, and refactoring to align better with common Go implementations and models is always first priority. Beyond that, future enhancements include:
* Unit Tests - Create the proper go-style unit testing files that can be consumed and executed by Go automatically on build.
* CI/CD - Utilize a defined integration and deployment pipeline for accepted pull reqests, and deploy PRs merged into master automatically.
* Containerization - Enable the CI/CD pipeline to automatically generate the appropriate Docker images that can be directly deployed within an orchestration platform.
* File Uploads - Currently WebUI only accepts text-based posts. However, the design intention is to also accept browser-based file uploads and retrievals. Ephem already allows this through the use of a BLOB as the secret store.
* Authenticated Secrets - Allow a submitter to specify the an individual that is allowed to retrieve a given secret, and only that individual. This further ensures that even if the secret URL is compromised the secret is still unretrievable by anyone except the specified recipient.
* SSO Integration - Authenticated Secrets implies the need for authentication. Authentication itself _is outside the scope of WebUI_, but integration with an authentication provider via SSO should be supported. 
* Notifications - When a target recipient is specified, automatically notify that recipient that a secret is available to them. Optionally, include the retrieval URL and expiration date/time. Notify the submitter when that secret is retrieved and by whom.
* Purge Startup Secrets - Both Ephem and WebUI require SSL certificates and Ephem requires database credentials and depending on configuration AES keys. These are expected to be present on the file system and are loaded on startup. After successfully starting and consuming these resources, purge them from the file system so as to prevent them from being retrieved afterward.
* Comments and Documentation - Implement the correct comments so as to enable Go's automatic documentation generation.

## Development Environment
* Obviously [Go](https://golang.org/dl/) is required and must be properly configured.
* Static assets are compiled directly into go binaries using [go-bindata](https://github.com/go-bindata/go-bindata) and any static asset updates must be compiled into the application using this tool.
* [Visual Studio Code](https://code.visualstudio.com/download) is the preferred development environment for this project. 
* The [Go Visual Studio Code Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) is required and should be properly installed and configured.
* [MySQL](https://dev.mysql.com/downloads/) is the current default backing database for development, and for development a locally available running instance of MySQL is expected.

## Contribution
* Use Github issues before making any changes.
* Use Branches for development with the branch naming format being '```[Issue ID]-[brief-change-description]```'. Avoid making multiple changes within a single branch unless implementation of the primary change requires the implementation of a secondary change.
* Commit frequently when something appears to work as intended before moving on to the next step in the implementation of your feature. 
* Use pull requests referencing the appropriate issue. PRs will be reviewed by the primary development team and accepted, commented, or rejected. Note that PRs that overly extend the scope of this project beyond its primary purpose may not be accepted at all. The use of issues before beginning development helps prevent this scenario.
* **DO NOT EMBED REAL/VALID SECRETS**. This project does contain non-real-world secrets (certificates, DB credentials) for ease and consistency of development. For development, these should always be used. Secrets that are "real world" valid should NEVER be used for development and should NEVER be committed to this project. 
* In case of fire:
  * git commit
  * git push
  * Locate exit

## Compiling Static Assets
Static assets such as content for the Web UI and .SQL files for the mysql data provider have been converted to binary objects using the go-bindata tool and compiled directly
into the ephem executable. Although not strictly necessary, this reduces deployment 
complexity by enabling a single-file deployment process and increases security by reducing the
likelihood that content files might be modified in the filesystem and served to clients or executed.

Whenever static assets are added or modified they must be re-converted to their binary representations using the following process:
1. Modify the appropriate static files in the filesystem
2. Run the go-bindata tool from the root of the project as follows:
   1. MySQL: go-bindata.exe -o src/data/mysql/procedures.go src/data/mysql/sql/...
   2. Web UI: go-bindata.exe -fs -prefix "src/webui/content/" -o src/webui/content.go src/webui/content/...
3. Modify the ```package main``` header such that the generated file properly references the correct project:
   1. MySQL: ```package mysql```
   2. Web UI: ```package webui```

In the future this may be integrated as part of a MakeFile or other automated build process
but is a manual step in the current project.

## Thanks for reading 😊
