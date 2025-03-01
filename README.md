
---

# **Load Balancer Project**

This project implements a **simple yet powerful load balancer** in Go. It distributes incoming HTTP requests across multiple backend servers using a **least active connections** algorithm. The load balancer also performs periodic health checks to ensure only healthy servers are used.

---

## **Features**
- **Load Balancing**: Distributes requests evenly across backend servers.
- **Health Checks**: Periodically checks the health of backend servers.
- **Logging**: Logs all requests and server activity for monitoring.
- **Configurable**: Backend servers and health check intervals are configurable via `config.json`.

---

## **Scenario**

Imagine you have three backend servers running on ports `9090`, `9091`, and `9092`. These servers handle incoming HTTP requests. However, you want to:
1. Distribute the load evenly across the servers.
2. Ensure that unhealthy servers are not used.
3. Monitor request handling and server activity.

This load balancer solves these problems by:
- Forwarding requests to the server with the least active connections.
- Periodically checking server health and excluding unhealthy servers.
- Logging all requests and server activity for visibility.

---

## **How It Works**

1. **Backend Servers**:
    - Backend servers are defined in `config.json`.
    - They handle the actual requests forwarded by the load balancer.

2. **Load Balancer**:
    - Listens on a configurable port (e.g., `:9095`).
    - Forwards requests to the backend server with the least active connections.
    - Logs all requests and server activity.

3. **Health Checks**:
    - Periodically sends HTTP GET requests to backend servers.
    - Marks servers as unhealthy if they fail to respond or return a status code >= 500.

---

## **Getting Started**

### **Prerequisites**
- Go installed on your machine.
- Python (optional, for simulating backend servers).

---

### **Step 1: Clone the Repository**
```bash
git clone https://github.com/kheder-hassoun/load-balncer-with-Go.git
cd load-balancer
```

---

### **Step 2: Set Up Backend Servers**
You can use simple HTTP servers to simulate backend servers. Run the following commands in separate terminal windows:

#### Terminal 1 (Server 1):
```bash
 java -jar .\out\artifacts\simpleWebApp_jar\simpleWebApp.jar 9090 "server 1"

```

#### Terminal 2 (Server 2):
```bash
 java -jar .\out\artifacts\simpleWebApp_jar\simpleWebApp.jar 9091 "server 2"

```

#### Terminal 3 (Server 3):
```bash
 java -jar .\out\artifacts\simpleWebApp_jar\simpleWebApp.jar 9092 "server 3"

```

These servers will act as your backend servers.

---

### **Step 3: Configure the Load Balancer**
Edit the `config.json` file to define your backend servers and load balancer settings:
```json
{
  "healthCheckInterval": "3s",
  "servers": [
    {
      "Name": "Server1",
      "URL": "http://localhost:9090",
      "ActiveConnections": 0,
      "Mutex": {},
      "Healthy": true
    },
    {
      "Name": "Server2",
      "URL": "http://localhost:9091",
      "ActiveConnections": 0,
      "Mutex": {},
      "Healthy": true
    },
    {
      "Name": "Server3",
      "URL": "http://localhost:9092",
      "ActiveConnections": 0,
      "Mutex": {},
      "Healthy": true
    }
  ],
  "listenPort": ":9095"
}
```

---

### **Step 4: Run the Load Balancer**
In the project directory, run:
```bash
go run .
```

The load balancer will start listening on port `9095`.

---

### **Step 5: Test the Load Balancer**
Send requests to the load balancer using `curl` or a browser:
```bash
curl http://localhost:9095
```

You can simulate multiple requests using a loop:
```bash
for i in {1..10}; do
  curl http://localhost:9095
done
```

---

### **Step 6: Monitor Logs**
The load balancer logs all requests and server activity to `logfile.log`. To monitor the logs in real-time, use:
```bash
tail -f logfile.log
```

Example log output:
```
Received request: GET / Server: localhost:9095
Request completed: GET / Server: localhost:9095 Server URL: http://localhost:9090
Received request: GET / Server: localhost:9095
Request completed: GET / Server: localhost:9095 Server URL: http://localhost:9091
```

---

## **Project Structure**
```
.
├── main.go            # Entry point of the application
├── config.go          # Configuration loading logic
├── server.go          # Server-related logic and health checks
├── loadbalancer.go    # Load balancing logic and request handling
├── config.json        # Configuration file for servers and settings
└── logfile.log        # Log file for request and server activity
```

---

## **Customization**
- **Add More Servers**: Edit `config.json` to add more backend servers.
- **Change Health Check Interval**: Update the `healthCheckInterval` field in `config.json`.
- **Modify Load Balancing Algorithm**: Edit the `nextServerLeastActive` function in `loadbalancer.go`.

---

## **Example Scenario**

1. Start three backend servers on ports `9090`, `9091`, and `9092`.
2. Start the load balancer on port `9095`.
3. Send 10 requests to the load balancer:
   ```bash
   for i in {1..10}; do
     curl http://localhost:9095
   done
   ```
4. Observe how the load balancer distributes requests across the backend servers in `logfile.log`.

---

## **Contributing**
Feel free to contribute to this project! Open an issue or submit a pull request.

---
Dev kheder hassoun 