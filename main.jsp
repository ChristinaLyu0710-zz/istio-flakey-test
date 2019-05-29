<html>
   <head>
      <title>Using GET Method to Read Form Data</title>
   </head>
   
   <body>
      <h1>Using GET Method to Read Form Data</h1>
      <ul>
         <li><p><b>First Name:</b>
            <%= request.getParameter("firstname")%>
         </p></li>
         <li><p><b>Last  Name:</b>
            <%= request.getParameter("lastname")%>
         </p></li>
      </ul>
   
   </body>
</html>