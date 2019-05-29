// Import required java libraries
import java.io.*;
import javax.servlet.*;
import javax.servlet.http.*;
 import javax.servlet.annotation.WebServlet;

// Extend HttpServlet class
@WebServlet("/helloworld")
public class HelloWorld extends HttpServlet {
 
   private String message;
   public HelloWorld(String m) {
      System.out.println(m);
   }

   public void init() throws ServletException {
      // Do required initialization
      message = "Hello World";
   }

   public void doGet(HttpServletRequest request, HttpServletResponse response)
      throws ServletException, IOException {
      
      // Set response content type
      response.setContentType("text/html");
      String firstName = request.getParameter("firstname");

      // Actual logic goes here.
      PrintWriter out = response.getWriter();
      out.println("<h1>" + firstName + "</h1>");
      request.setAttribute("message", firstName); // This will be available as ${message}
      request.getRequestDispatcher("response.jsp").forward(request, response);
   }

   @Override
    protected void doPost(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        
      String firstName = request.getParameter("firstname");
      System.out.println("First name = " + firstName);
      request.setAttribute("message", firstName); 
      request.getRequestDispatcher("response.jsp").forward(request, response);
    }



   public void destroy() {
      // do nothing.
   }
}