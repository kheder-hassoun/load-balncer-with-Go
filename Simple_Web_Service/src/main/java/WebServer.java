import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
public class WebServer {
    public static void main(String[] args) throws IOException {
        if (args.length != 1) {
            System.err.println("Usage: java -jar server.jar <port>");
            System.exit(1);
        }
        int port = Integer.parseInt(args[0]);
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/", new MainHandler(port));
        server.start();
        System.out.println("Starting web service on port " + port);
    }
    static class MainHandler implements HttpHandler {
        private int port;
        public MainHandler(int port) {
            this.port = port;
        }
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            String response = "<html><body><h1>Hello, World!</h1>" +
                    "<p>Server Address: " + exchange.getLocalAddress().getHostName() + "</p>" +
                    "<p>Server Port: " + port + "</p></body></html>";

            exchange.sendResponseHeaders(200, response.getBytes().length);
            OutputStream os = exchange.getResponseBody();
            os.write(response.getBytes());
            os.close();
        }
    }
}
