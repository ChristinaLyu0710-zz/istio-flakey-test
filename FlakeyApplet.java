import java.applet.Applet; // Provides the Applet class.
import java.awt.*;         // Provides Button class, etc.
import java.awt.event.*;   // Provides ActionEvent, ActionListener 

public class FlakeyApplet extends Applet {
	public static void run() {
		String test = TotalFlakey.testHTML();
		System.out.println(test);
	}
}