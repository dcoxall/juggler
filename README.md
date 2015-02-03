[![Build Status](https://travis-ci.org/dcoxall/juggler.svg?branch=master)](https://travis-ci.org/dcoxall/juggler)

Juggler
=======

A Go library that acts as a proxy to a pool of multiple application servers.

Development
-----------

The recommended environment requires the following:

- [Docker][docker] / [boot2docker][boot2docker]
- [Fig][fig]

Once installed just use the following to run the test suite...

    $ fig run juggler make

You can also preview the documentation by running...

    $ fig up

Once the server is running you can visit it with:

    $ open http://`boot2docker ip`:6060/pkg/github.com/dcoxall/juggler/

Contributing
------------

1. Fork the project
2. Create a feature branch
3. Write tests for the new feature
4. Implement new features to pass the tests
5. Document all functions and attributes
6. Ensure you have run go fmt
7. Create a pull request

[docker]: http://docker.com/
[boot2docker]: https://github.com/boot2docker/boot2docker
[fig]: http://fig.sh/
