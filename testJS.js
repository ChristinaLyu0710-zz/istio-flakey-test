
browser.addScriptContextListener(new ScriptContextAdapter() {
    @Override
    public void onScriptContextCreated(ScriptContextEvent event) {
        Browser browser = event.getBrowser();
        JSValue window = browser.executeJavaScriptAndReturnValue("window");
        window.asObject().setProperty("java", new JavaObject());
    }
});
window.java.print('Hello Java!');

var input = function() {

	var myClass = Java.type("TotalFlakey");
	print(myClass.testHTML());
}