<%@ page import="TotalFlakey"%>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <title>JSP Page</title>
</head>
<body>
<%
  TotalFlakey tc = new TotalFlakey();
  print(tc.testHTML());
%>
</body>
</html>