# Simple Web Server for HAProxy Lab

This project provides a basic web server implementation written in Java. It is designed to be used as a practical example for lab exercises involving HAPROXY, a popular open-source load balancer and proxy server.

## Table of Contents
1. [Introduction](#introduction)
2. [Features](#features)
3. [Getting Started](#getting-started)
   * [Running the Web Server](#running-the-web-server)
4. [HAProxy Configuration](#haproxy-configuration)
5. [Usage Examples with HAProxy](#usage-examples-with-haproxy)


## Introduction

The Simple Web Server is intended to be deployed as part of a HAProxy proxy lab setup. This server can serve a simple HTML page and respond to status check requests, making it an excellent candidate for load balancing experiments using HAProxy.

## Features

- Serves a static HTML page.
- Allows configuration of the server name via command-line arguments.
- Responsive to health checks for load balancing applications (using `/status` endpoint).


## Getting Started

### Running the Web Server

To run a server instance, you need to specify `PORT_NUMBER` and `SERVER_NAME`:

```sh
java -jar SimpleWebServer.jar [PORT_NUMBER] [SERVER_NAME]
```

Or build and run using Maven:

```sh
mvn exec:java -Dexec.mainClass=org.ds.Main -Dexec.args="[PORT_NUMBER] [SERVER_NAME]"
```

Replace `[PORT_NUMBER]` with an actual port number (e.g., `8080`) and `[SERVER_NAME]` with your desired server name (e.g., `Server1`).

## HAProxy Configuration

Here is a sample HAPROXY configuration for load balancing across instances of the Simple Web Server:

```conf
frontend http_front
    bind *:80
    default_backend servers

backend servers
    balance roundrobin
    server srv1 127.0.0.1:8081 check
    server srv2 127.0.0.1:8082 check
```

In this example:
- Instances of the Simple Web Server are run on `127.0.0.1` with ports `8081` and `8082`.
- The HAPROXY frontend listens on port 80 for incoming HTTP requests.
- Requests are distributed to the instances using a round-robin algorithm.

## Usage Examples with HAProxy

For your lab setup:
1. Start two Simple Web Server instances:

```sh
java -jar SimpleWebServer.jar 8081 Server1 &
java -jar SimpleWebServer.jar 8082 Server2 &
```

2. Configure and start HAPROXY using the sample configuration provided above.

3. Access any of the server instances through your browser at `http://localhost`. You should see requests distributed between `Server1` and `Server2`.
